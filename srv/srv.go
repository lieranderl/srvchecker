package srv

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var srvtextlist = []string{
	"_collab-edge._tls", 
	"_sips._tcp", 
	"_sip._tcp", 
	"_sip._udp", 
	"_h323cs._tcp", 
	"_h323ls._udp", 
	"_xmpp-server._tcp", 
	"_xmpp-client._tcp", 
	"_sipfederationtls._tcp", 
	"_sips._tcp.sipmtls"}

type inputSRV struct {
	service string
	proto   string
	domain  string
}

type inputSRVlist []inputSRV

func (s *inputSRVlist) Init(domain string) {
	var isrv inputSRV
	isrv.domain = domain
	for _, i := range srvtextlist {
		ii := strings.Split(i, ".")
		isrv.service = strings.TrimPrefix(ii[0], "_")
		isrv.proto = strings.TrimPrefix(ii[1], "_")
		*s = append(*s, isrv)
	}
}


type SrvResult struct {
	Cname 	 string
	Fqdn     string
	Ips      []string
	Port     string
	Priority string
	Weight   string
}

func (s *SrvResult) fetch(cname string, fqdn string, ips []string, port uint16, priority uint16, weight uint16) {
	s.Cname = cname
	s.Fqdn = fqdn
	s.Ips = ips
	if port == 0 {
		s.Port = ""
	} else {
		s.Port = fmt.Sprint(port)
	}
	if priority == 0 {
		s.Priority = ""
	} else {
		s.Priority = fmt.Sprint(priority)
	}
	if weight == 0 {
		s.Weight = ""
	} else {
		s.Weight = fmt.Sprint(weight)
	}
}


func (s *SrvResult) fetchAddr(addr *net.SRV, cname string, result chan SrvResult) {
	ips, err := net.LookupHost(addr.Target)
	if err != nil {
		s.fetch(cname, addr.Target, []string{"IP address is not resolved"}, 0, 0, 0)
	} else {
		s.fetch(cname, addr.Target, ips, addr.Port, addr.Priority, addr.Weight)
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
	output := make(chan SRVResults)
	var wg sync.WaitGroup

	go s.handleResults(input, output, &wg)
	defer close(output)

	for _, srv := range *mysrvs {
		mySrvResult := new(SrvResult)
		cname := "_"+srv.service+"._"+srv.proto+"."+srv.domain
		_, addrs, err := net.LookupSRV(srv.service, srv.proto, srv.domain)
		if err != nil {
			wg.Add(1)
			mySrvResult.fetch(cname, "SRV record not configured.", []string{""}, 0, 0, 0)
			input <- *mySrvResult
		} else {
			for _, addr := range addrs {
				wg.Add(1)
				go mySrvResult.fetchAddr(addr, cname, input)
			}
		}
	}
	wg.Wait()
	close(input)
	
}

func (s *SRVResults)handleResults(input chan SrvResult, output chan SRVResults, wg *sync.WaitGroup) {
	for result := range input {
		(*s)[result.Cname] = append((*s)[result.Cname], result)
		wg.Done()
	}
}