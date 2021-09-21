package srv

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

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

type Ip struct {
	Ips 	[]string
	Priority string
	Weight   string
}
type Fqdns map[string]*Ip

type SrvResult struct {
	Sname    		string
	Fqdn     		Fqdns
	Port     		string
	Proto	 		string	
}


type SrvResults map[string]*SrvResult


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


func (s *SrvResult) fetch(servname string, fqdn string, ips []string, port uint16, proto string, priority uint16, weight uint16) {
	ip := new(Ip)
	ip.Ips = ips
	s.Proto = proto
	s.Sname = servname
	ip.Priority = fmt.Sprint(priority)
	ip.Weight = fmt.Sprint(weight)
	s.Fqdn[fqdn] = ip
	if port == 0 {
		s.Port = ""
	} else {
		s.Port = fmt.Sprint(port)
	}

}


func (s *SRVResults) fetchAddr(cname string, fqdn *net.SRV, servname string, proto string, newRes *SrvResult, wg *sync.WaitGroup) {
	var mutex = &sync.RWMutex{}
	ips, err := net.LookupHost(fqdn.Target)
	if err != nil {
		mutex.Lock()
		newRes.fetch(servname, fqdn.Target, []string{"A record not configured"}, 0, proto, fqdn.Priority, fqdn.Weight)
		mutex.Unlock()
	} 
	if len(ips)>0 {
		mutex.Lock()
		newRes.fetch(servname, fqdn.Target, ips, fqdn.Port, proto, fqdn.Priority, fqdn.Weight)
		mutex.Unlock()
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
	var mutex = &sync.RWMutex{}

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
