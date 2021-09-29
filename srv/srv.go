package srv

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
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

// type Ip string
// type Ips struct {
// 	Ips 	 map[Ip]*Port
// 	Priority string
// 	Weight   string
// }
// type portnum string
// type portproto string

// type Portidenty string
// func GetPortidenty(portnum portnum,portproto portproto) Portidenty {
// 	return Portidenty(string(portnum) + ":" + string(portproto))
// }

// type Port map[Portidenty]*PortStatus
// type PortStatus struct {
// 	IsOpen			bool
// 	// Cert 			[]*x509.Certificate	
// 	Cert string
// }

// type Fqdn string
// type Fqdns map[Fqdn]*Ips

// type SrvResult struct {
// 	Sname    		string
// 	Fqdn     		Fqdns
// }

// type Cname string
// type SrvResults map[Cname]*SrvResult


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

// func (ps *PortStatus)connect_cert(ip string, port string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	timeout := time.Second
// 	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
// 	if err != nil {
// 		ps.IsOpen = false
// 	}
// 	if conn != nil {
// 		defer conn.Close()
// 		ps.IsOpen = true
// 		if (port == "8443" || port== "5061") {
// 			ps.Cert = GetCert(ip , fmt.Sprint(port))[0].Issuer.CommonName
// 		}
// 	}
// }


// func (s *SrvResult) fetch(fqdn string, ips []string, port uint16, proto string, priority uint16, weight uint16) {
	
// 	if strings.Contains(fqdn, ".") {
// 		myips := new(Ips)
// 		myips.Priority = fmt.Sprint(priority)
// 		myips.Weight = fmt.Sprint(weight)
// 		var wg sync.WaitGroup
// 		myips.Ips = make(map[Ip]*Port)
		
// 		for _, ip := range ips {
// 			if len(ip)>0 {
// 				myips.Ips[Ip(ip)] = &Port{}
// 				if strings.Contains(ip, ".") {
// 					pi := GetPortidenty(portnum(fmt.Sprint(port)), portproto(proto))
// 					myips.Ips[Ip(ip)] = &Port{pi:new(PortStatus)}
// 					currentport := (*myips.Ips[Ip(ip)])[pi]
// 					if port != 0 {
// 						if proto == "tcp" {
// 							wg.Add(1)
// 							go currentport.connect_cert(ip, fmt.Sprint(port), &wg)
// 						} else {
// 							currentport.IsOpen = true
// 						}
// 					}
// 				}
// 			}
// 		} 
// 		wg.Wait()
// 		mu.Lock()
// 		s.Fqdn[Fqdn(fqdn)] = myips
// 		mu.Unlock()
// 	} else {
// 		mu.Lock()
// 		s.Fqdn[Fqdn(fqdn)] = &Ips{}
// 		mu.Unlock()
// 	}
	
// }

func (ps *DiscoveredSrvRow)Connect_cert(ip string, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpened = false
	}
	if conn != nil {
		defer conn.Close()
		ps.IsOpened = true
		if (port == "8443" || port== "5061") {
			ps.Cert = GetCert(ip , fmt.Sprint(port))[0].Issuer.CommonName
		}
	}
}

func GetCert(ip string, port string) []*x509.Certificate {
	conf := &tls.Config{
        InsecureSkipVerify: true,
    }
    conn, err := tls.Dial("tcp", ip+":"+port, conf)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
    if err != nil {
        log.Println("Error in Dial", err)
    }
	if conn != nil {
		defer conn.Close()
		return conn.ConnectionState().PeerCertificates
	}
    return nil
}


// func (s *SRVResults) fetchAddr(cname string, fqdn *net.SRV, proto string, newRes *SrvResult, wg *sync.WaitGroup) {
// 	ips, err := net.LookupHost(fqdn.Target)
// 	if err != nil {
// 		newRes.fetch(fqdn.Target, []string{"A record not configured"}, 0, proto, fqdn.Priority, fqdn.Weight)
// 	} 
// 	if len(ips)>0 {
// 		newRes.fetch(fqdn.Target, ips, fqdn.Port, proto, fqdn.Priority, fqdn.Weight)
// 	}
// 	(*s)[cname] = *newRes
// 	wg.Done()
// }

// type SRVResults map[string]SrvResult

// func (s *SRVResults) Init() {
// 	*s= make(map[string]SrvResult)
// }

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
	ips, err := net.LookupHost(fqdn.Target)
	if err != nil {
		discoveredSrvRow := new(DiscoveredSrvRow)
		discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fmt.Sprint(fqdn.Port), "A record not configured" ,proto)
		*d = append(*d, discoveredSrvRow)				
	} 
	if len(ips)>0 {
		for _, ip := range ips {
			discoveredSrvRow := new(DiscoveredSrvRow)
			discoveredSrvRow.Init(cname, servName, fmt.Sprint(fqdn.Priority), fmt.Sprint(fqdn.Weight), fqdn.Target, fmt.Sprint(fqdn.Port), ip ,proto)
			if proto == "tcp" {
				discoveredSrvRow.Connect_cert(ip, fmt.Sprint(fqdn.Port))
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
