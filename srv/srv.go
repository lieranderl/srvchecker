package srv

import (
	"context"
	"log"

	"crypto/tls"
	"crypto/x509"
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
	Port 		string
	Proto 		string
	IsOpened 	bool
	Cert 		string
	ServiceName string
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
		cert := GetCert(ip , fmt.Sprint(port))
		if cert != nil {
			ps.Cert = cert[len(cert)-1].Issuer.CommonName
		}
	}
}

func GetCert(ip string, port string) []*x509.Certificate {

	conf := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout:  2 * time.Second}, "tcp", ip+":"+port, &conf)
	if err != nil { 
		log.Println("Host:", ip,":",port, "Dial:", err)
	}
	if err == nil {
		defer conn.Close()
		return conn.ConnectionState().PeerCertificates
	}
    return nil
}


func (s *DiscoveredSrvTable) ForDomain(domain string) {
	mysrvs := new(inputSRVlist)
	mysrvs.Init(domain)	
	var wg sync.WaitGroup

	for _, srv := range *mysrvs {
		fmt.Println(srv)
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
		discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fmt.Sprint(fqdn.Port), "A record not configured" ,proto)
		*d = append(*d, discoveredSrvRow)				
	} 
	if len(ips)>0 {
		for _, ip := range ips {
			discoveredSrvRow := new(DiscoveredSrvRow)
			discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fmt.Sprint(fqdn.Port), ip.To4().String() ,proto)
			if proto == "tcp" {
				discoveredSrvRow.Connect_cert(ip.To4().String(), fmt.Sprint(fqdn.Port))
			}
			*d = append(*d, discoveredSrvRow)				
		}
	}
}

func (d *DiscoveredSrvRow) Init(cname, servName, priority, weight, fqdn, port, ip, proto string) {
	d.Srv = cname
	d.ServiceName = servName
	d.Priority = priority
	d.Weight = weight
	d.Fqdn = fqdn
	d.Port = port
	d.Ip = ip
	d.Proto = proto
}
