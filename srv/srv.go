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


type SrvResult struct {
	Cname 	 string
	Fqdn     string
	Ips      []string
	Port     string
	Priority string
	Weight   string
	ServName string
}

func (s *SrvResult) fetch(servname string, cname string, fqdn string, ips []string, port uint16, priority uint16, weight uint16) {
	s.Cname = cname
	s.Fqdn = fqdn
	s.Ips = ips
	s.ServName = servname
	if port == 0 {
		s.Port = ""
	} else {
		s.Port = fmt.Sprint(port)
	}
	s.Priority = fmt.Sprint(priority)
	s.Weight = fmt.Sprint(weight)
}


func (s *SrvResult) fetchAddr(addr *net.SRV, cname string, servname string, result chan SrvResult) {
	ips, err := net.LookupHost(addr.Target)
	if err != nil {
		s.fetch(servname, cname, addr.Target, []string{"A record not configured"}, 0, 0, 0)
	} else {
		s.fetch(servname, cname, addr.Target, ips, addr.Port, addr.Priority, addr.Weight)
	}
	result <- *s
}


type SRVResults map[string][]SrvResult

func (s *SRVResults) Init() {
	*s= make(map[string][]SrvResult)
}


func (s *SRVResults) GetForDomain(domain string) {
	mysrvs := new(inputSRVlist)
	s.Init()
	mysrvs.Init(domain)
	input := make(chan SrvResult)
	var wg sync.WaitGroup

	go s.handleResults(input, &wg)

	for _, srv := range *mysrvs {
		mySrvResult := new(SrvResult)
		cname := "_"+srv.service+"._"+srv.proto+"."+srv.domain
		_, addrs, err := net.LookupSRV(srv.service, srv.proto, srv.domain)
		if err != nil {
			wg.Add(1)
			mySrvResult.fetch(srv.servName, cname, "SRV record not configured", []string{""}, 0, 0, 0)
			input <- *mySrvResult
		} else {
			for _, addr := range addrs {
				wg.Add(1)
				go mySrvResult.fetchAddr(addr, cname, srv.servName, input)
			}
		}
	}
	wg.Wait()
	close(input)
	
}

func (s *SRVResults)handleResults(input chan SrvResult, wg *sync.WaitGroup) {
	for result := range input {
		(*s)[result.Cname] = append((*s)[result.Cname], result)
		wg.Done()
	}
}