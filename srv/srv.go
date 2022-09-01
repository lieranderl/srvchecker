package srv

import (
	"context"
	
	"crypto/tls"
	"crypto/x509"
	"github.com/grantae/certinfo"
	"fmt"

	"net"

	"sort"
	"strings"
	"sync"
	"time"
)

var (
	mra_srv = []string{"_collab-edge:_tls", "_cuplogin:_tcp", "_cisco-uds:_tcp"}
 	b2b_srv = []string{"_h323cs:_tcp", "_sip:_tcp", "_sips:_tcp", "_sip:_udp", "_h323ls:_udp"}
	xmpp_fed_srv = []string{"_xmpp-server:_tcp"}
	cma_srv = []string{"_xmpp-client:_tcp"}
	spark_srv = []string{"_sips:_tcp.sipmtls"}
	mssip_srv = []string{"_sipfederationtls:_tcp"}
	srvtextlist = map[string][]string{
		"mra":mra_srv, 
		"b2b":b2b_srv,
		"xmpp_fed":xmpp_fed_srv,
		"cma":cma_srv,
		"spark":spark_srv,
		"mssip":mssip_srv,
	}
)

type DiscoveredSrvRow struct {
	Srv 		string
	Fqdn 		string
	Ip 			string
	Priority 	string
	Weight 		string 
	Port 		uint16
	Proto 		string
	IsOpened 	bool
	Cert 		string
	Certs 		[]*Cert
	ServiceName string
}

type Cert struct {
	Txt string
	Cn string
	Subject string
	San string
	KeyUsage []string
	ExtKeyUsage []string
	Issuer string
	NotBefore string
	NotAfter string
	Child []*Cert
}

type DiscoveredSrvTable []*DiscoveredSrvRow

type inputSRV struct {
	service 	string
	proto   	string
	domain  	string
	servName 	string
}

type inputSRVlist []inputSRV


func (s *inputSRVlist) Init(domain string) {
	var isrv inputSRV
	isrv.domain = domain
	for k, vl := range srvtextlist {
		isrv.servName = k
		for _, v := range vl {
			ii := strings.Split(v, ":")
			isrv.service = strings.TrimPrefix(ii[0], "_")
			isrv.proto = strings.TrimPrefix(ii[1], "_")
			*s = append(*s, isrv)
		} 
	}
}


func (ps *DiscoveredSrvRow)Connect_cert(ip string, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpened = false
	}
	if err == nil {
		defer conn.Close()
		ps.IsOpened = true
		certs := GetCert(ip , fmt.Sprint(port))
		if certs != nil {
			ps.Cert = certs[0].Subject.CommonName
		}
		for i, cert :=range certs {
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
			
			if i == 0 {
				ps.Certs = append(ps.Certs, c)
			}
			if i == 1 {
	
				ps.Certs[0].Child = append(ps.Certs[0].Child, c)
			}
			if i == 2 {

				ps.Certs[0].Child[0].Child = append(ps.Certs[0].Child[0].Child, c)
			}
			if i == 3 {
		
				ps.Certs[0].Child[0].Child[0].Child = append(ps.Certs[0].Child[0].Child[0].Child, c)
			}
			if i == 4 {
				ps.Certs[0].Child[0].Child[0].Child[0].Child = append(ps.Certs[0].Child[0].Child[0].Child[0].Child, c)
			}
			// if i == 5 {
			// 	ps.Certs.Child.Child.Child.Child.Child = c
			// }
			
		}

	}
}

func GetCert(ip string, port string) []*x509.Certificate {

	conf := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout:  2 * time.Second}, "tcp", ip+":"+port, &conf)
	if err == nil {
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		reversed := make([]*x509.Certificate,0)
		for i := range certs {
				n := certs[len(certs)-1-i]
				reversed = append(reversed, n)
		}
		return reversed
	}
    return nil
}


func (s *DiscoveredSrvTable) ForDomain(domain string) {
	mysrvs := new(inputSRVlist)
	mysrvs.Init(domain)	
	var wg sync.WaitGroup

	for _, srv := range *mysrvs {
		proto := "udp"
		if strings.HasPrefix(srv.proto, "t") {
			proto = "tcp"
		}
		cname := "_"+srv.service+"._"+srv.proto+"."+srv.domain
		_, fqdns, err := net.LookupSRV(srv.service, srv.proto, srv.domain)
		
		if err != nil {
			discoveredSrvRow := new(DiscoveredSrvRow)
			discoveredSrvRow.ServiceName = srv.servName
			discoveredSrvRow.Srv = cname
			discoveredSrvRow.Fqdn = "SRV record not configured"		
			*s = append(*s, discoveredSrvRow)	
		} else {
			for _, fqdn := range fqdns {
				wg.Add(1)
				go s.fetchIps(srv.servName, cname, fqdn, proto, &wg)
			}
		}
	}
	wg.Wait()
	sort.Slice((*s)[:], func(i, j int) bool {
		return (*s)[i].Srv < (*s)[j].Srv
	})
}

func (d *DiscoveredSrvTable) fetchIps(servName, cname string, fqdn *net.SRV, proto string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", fqdn.Target)
	if err != nil {
		discoveredSrvRow := new(DiscoveredSrvRow)
		discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fqdn.Port, "A record not configured" ,proto)
		*d = append(*d, discoveredSrvRow)				
	} 
	if len(ips)>0 {
		for _, ip := range ips {
			discoveredSrvRow := new(DiscoveredSrvRow)
			discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fqdn.Port, ip.To4().String() ,proto)
			if proto == "tcp" {
				discoveredSrvRow.Connect_cert(ip.To4().String(), fmt.Sprint(fqdn.Port))
			}
			*d = append(*d, discoveredSrvRow)				
		}
	}
}

func (d *DiscoveredSrvRow) Init(cname, servName, priority, weight, fqdn string, port uint16, ip, proto string) {
	d.Srv = cname
	d.ServiceName = servName
	d.Priority = priority
	d.Weight = weight
	d.Fqdn = fqdn
	d.Port = port
	d.Ip = ip
	d.Proto = proto
}
