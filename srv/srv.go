package srv

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.RWMutex{}
var mra_srv = []string{"_collab-edge:_tls", "_cuplogin:_tcp", "_cisco-uds:_tcp"}
var b2b_srv = []string{"_h323cs:_tcp", "_sip:_tcp", "_sips:_tcp", "_sip:_udp", "_h323ls:_udp"}
var xmpp_fed_srv = []string{"_xmpp-server:_tcp"}
var cma_srv = []string{"_xmpp-client:_tcp"}
var spark_srv = []string{"_sips:_tcp.sipmtls"}
var mssip_srv = []string{"_sipfederationtls:_tcp"}
var srvtextlist = map[string][]string{
	"mra":mra_srv, 
	"b2b":b2b_srv,
	"xmpp_fed":xmpp_fed_srv,
	"cma":cma_srv,
	"spark":spark_srv,
	"mssip":mssip_srv,
}

type inputSRV struct {
	service 	string
	proto   	string
	domain  	string
	servName 	string
}

type inputSRVlist []inputSRV

type Ip string
type Ips struct {
	Ips 	 map[Ip]*Port
	Priority string
	Weight   string
}
type portnum string
type portproto string

type Portidenty string
func GetPortidenty(portnum portnum,portproto portproto) Portidenty {
	return Portidenty(string(portnum) + ":" + string(portproto))
}


type Port map[Portidenty]*PortStatus
type PortStatus struct {
	IsOpen			bool
	// Cert 			[]*x509.Certificate	
	Cert string
}


type Fqdn string
type Fqdns map[Fqdn]*Ips

type SrvResult struct {
	Sname    		string
	Fqdn     		Fqdns
}

type Cname string
type SrvResults map[Cname]*SrvResult


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

func (ps *PortStatus)connect_cert(ip string, port string, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpen = false
	}
	if conn != nil {
		defer conn.Close()
		ps.IsOpen = true
		if (port == "8443" || port== "5061") {
			ps.Cert = GetCert(ip , fmt.Sprint(port))[0].Issuer.CommonName
		}
	}
}


func (s *SrvResult) fetch(servname string, fqdn string, ips []string, port uint16, proto string, priority uint16, weight uint16) {
	myips := new(Ips)
	s.Sname = servname
	myips.Priority = fmt.Sprint(priority)
	myips.Weight = fmt.Sprint(weight)
	var wg sync.WaitGroup

	
	for _, ip := range ips {
		if _, ok := myips.Ips[Ip(ip)]; !ok {
			myips.Ips = make(map[Ip]*Port)
		}
		pi := GetPortidenty(portnum(fmt.Sprint(port)), portproto(proto))
		myips.Ips[Ip(ip)] = &Port{pi:new(PortStatus)}

		if port != 0 {
			if proto == "tcp" {
				wg.Add(1)
				go (*myips.Ips[Ip(ip)])[pi].connect_cert(ip, fmt.Sprint(port), &wg)
			} else {
				(*myips.Ips[Ip(ip)])[pi].IsOpen = true
			}
		}

	} 
	wg.Wait()

	s.Fqdn[Fqdn(fqdn)] = myips

	
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

func (s *SRVResults) fetchAddr(cname string, fqdn *net.SRV, servname string, proto string, newRes *SrvResult, wg *sync.WaitGroup) {
	ips, err := net.LookupHost(fqdn.Target)
	if err != nil {
		// mutex.Lock()
		newRes.fetch(servname, fqdn.Target, []string{"A record not configured"}, 0, proto, fqdn.Priority, fqdn.Weight)
		// mutex.Unlock()
	} 
	if len(ips)>0 {
		// mutex.Lock()
		newRes.fetch(servname, fqdn.Target, ips, fqdn.Port, proto, fqdn.Priority, fqdn.Weight)
		// mutex.Unlock()
	}
	(*s)[cname] = *newRes
	wg.Done()
}


type SRVResults map[string]SrvResult

func (s *SRVResults) Init() {
	*s= make(map[string]SrvResult)
}


func (s *SRVResults) ForDomain(domain string) {
	mysrvs := new(inputSRVlist)
	s.Init()
	mysrvs.Init(domain)
	input := make(chan SrvResult)

	var wg sync.WaitGroup

	for _, srv := range *mysrvs {
		proto := "udp"
		if strings.HasPrefix(srv.proto, "t") {
			proto = "tcp"
		}
		cname := "_"+srv.service+"._"+srv.proto+"."+srv.domain
		mySrvResult := new(SrvResult)
		mySrvResult.Fqdn = make(Fqdns)

		_, fqdns, err := net.LookupSRV(srv.service, srv.proto, srv.domain)
		if err != nil {
			mutex.Lock()
			mySrvResult.fetch(srv.servName, "SRV record not configured", []string{""}, 0, proto, 0, 0)
			mutex.Unlock()
			(*s)[cname] = *mySrvResult
		} else {
			for _, fqdn := range fqdns {
				wg.Add(1)
				go s.fetchAddr(cname, fqdn, srv.servName, proto, mySrvResult, &wg)
			}
		}
	}
	wg.Wait()
	close(input)
}
