package srv

import (
	"context"

	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/grantae/certinfo"

	"net"

	"sort"
	"strings"
	"sync"
	"time"
)

type CertsChain []*Cert

type DiscoveredSrvRow struct {
	Srv      string
	Fqdn     string
	Ip       string
	Priority string
	Weight   string
	Port     uint16
	Proto    string
	IsOpened bool
	CertsChain
	ServiceName string
}

type Cert struct {
	Txt         string
	Cn          string
	Subject     string
	San         string
	KeyUsage    []string
	ExtKeyUsage []string
	Issuer      string
	NotBefore   string
	NotAfter    string
}

type DiscoveredSrvTable []*DiscoveredSrvRow

type inputSRV struct {
	service  string
	proto    string
	domain   string
	servName string
}

var SRVTextList = map[string][]string{
	"mra":      {"_collab-edge:_tls", "_cuplogin:_tcp", "_cisco-uds:_tcp"},
	"b2b":      {"_h323cs:_tcp", "_sip:_tcp", "_sips:_tcp", "_sip:_udp", "_h323ls:_udp"},
	"xmpp_fed": {"_xmpp-server:_tcp"},
	"cma":      {"_xmpp-client:_tcp"},
	"spark":    {"_sips:_tcp.sipmtls"},
	"mssip":    {"_sipfederationtls:_tcp"},
}

type inputSRVlist []inputSRV

func (s *inputSRVlist) init(domain string) {
	// Calculate the total number of entries to preallocate the slice.
	totalEntries := 0
	for _, entries := range SRVTextList {
		totalEntries += len(entries)
	}
	*s = make(inputSRVlist, 0, totalEntries) // Preallocate the slice.

	// Create a single inputSRV struct and reuse it.
	var isrv inputSRV
	isrv.domain = domain

	for serviceName, srvEntries := range SRVTextList {
		isrv.servName = serviceName
		for _, srvEntry := range srvEntries {
			parts := strings.Split(srvEntry, ":")
			if len(parts) != 2 {
				// Handle the error according to your application's needs.
				// For example, log an error message or continue to the next entry.
				continue
			}
			isrv.service = strings.TrimPrefix(parts[0], "_")
			isrv.proto = strings.TrimPrefix(parts[1], "_")
			*s = append(*s, isrv)
		}
	}
}

func (ps *DiscoveredSrvRow) Connect_cert(ip string, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpened = false
	}
	if err == nil {
		defer conn.Close()
		ps.IsOpened = true
		certs, err := getCert(ip, fmt.Sprint(port))
		// check error
		if err != nil {
			ps.CertsChain = nil
		}

		for _, cert := range certs {
			c := new(Cert)
			c.Cn = cert.Subject.CommonName
			c.Issuer = cert.Issuer.CommonName
			c.Subject = cert.Subject.String()
			c.San = strings.Join(cert.DNSNames, ", ")
			c.NotBefore = cert.NotBefore.String()
			c.NotAfter = cert.NotAfter.String()

			for _, eku := range cert.ExtKeyUsage {
				if eku == x509.ExtKeyUsageServerAuth {
					c.ExtKeyUsage = append(c.ExtKeyUsage, "TLS Web Server Authentication")
				}
				if eku == x509.ExtKeyUsageClientAuth {
					c.ExtKeyUsage = append(c.ExtKeyUsage, "TLS Web Client Authentication")
				}
			}

			c.Txt, _ = certinfo.CertificateText(cert)
			ps.CertsChain = append(ps.CertsChain, c)

		}

	}
}

func getCert(ip, port string) ([]*x509.Certificate, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "tcp", net.JoinHostPort(ip, port), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func (d *DiscoveredSrvRow) init(cname, servName, priority, weight, fqdn string, port uint16, ip, proto string) {
	d.Srv = cname
	d.ServiceName = servName
	d.Priority = priority
	d.Weight = weight
	d.Fqdn = fqdn
	d.Port = port
	d.Ip = ip
	d.Proto = proto
}

func (d *DiscoveredSrvTable) fetchIps(servName, cname string, fqdn *net.SRV, proto string, wg *sync.WaitGroup) {
	defer wg.Done()
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", fqdn.Target)
	discoveredSrvRow := new(DiscoveredSrvRow)
	discoveredSrvRow.init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fqdn.Port, "A record not configured", proto)
	if err != nil {
		*d = append(*d, discoveredSrvRow)
		return
	}
	if len(ips) > 0 {
		for _, ip := range ips {
			discoveredSrvRow.Ip = ip.To4().String()
			if proto == "tcp" {
				discoveredSrvRow.Connect_cert(ip.To4().String(), fmt.Sprint(fqdn.Port))
			}
			*d = append(*d, discoveredSrvRow)
		}
	}
}

func (s *DiscoveredSrvTable) ForDomain(domain string) {
	mysrvs := inputSRVlist{}
	mysrvs.init(domain)

	var wg sync.WaitGroup

	for _, srv := range mysrvs {
		proto := "udp"
		if strings.HasPrefix(srv.proto, "t") {
			proto = "tcp"
		}
		cname := fmt.Sprintf("_%s._%s.%s", srv.service, srv.proto, srv.domain)
		_, fqdns, err := net.LookupSRV(srv.service, srv.proto, srv.domain)

		if err != nil {
			// Consider logging the error before continuing.
			*s = append(*s, &DiscoveredSrvRow{
				ServiceName: srv.servName,
				Srv:         cname,
				Fqdn:        "SRV record not configured",
			})
			continue
		}

		for _, fqdn := range fqdns {
			wg.Add(1)
			go s.fetchIps(srv.servName, cname, fqdn, proto, &wg)
		}
	}

	wg.Wait()

	// Perform sorting.
	sort.Slice(*s, func(i, j int) bool {
		return (*s)[i].Srv < (*s)[j].Srv
	})
}
