package portconnectivity

import (
	"net"
	"srvchecker/srv"
	"strings"
	"sync"
	"time"
)



type PortsResult struct {
	Ip    string
	Fqdn  string
	Ports map[string]bool
	ServName string
}

func (p *PortsResult) Init(ip string, fqdn string, servname string) {
	p.Ip = ip
	p.Fqdn = fqdn
	p.ServName = servname
}

func (p *PortsResult) Run(ip string, port string, result chan PortsResult) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	
	if err != nil {
		p.Ports = map[string]bool{port:false}
		
	}
	if conn != nil {
		defer conn.Close()
		p.Ports = map[string]bool{port:true}
	}
	result <- *p
}


type PortsResults []PortsResult

func (p *PortsResults)handleResults(input chan PortsResult, wg *sync.WaitGroup) {
	for result := range input {
		*p = append(*p, result)
		wg.Done()
	}
}


func (p *PortsResults) Connectivity(srvresults srv.SRVResults){

	input := make(chan PortsResult)
	var wg sync.WaitGroup
	go p.handleResults(input, &wg)

	for _,v := range srvresults {
		for _, entry := range v {
			if entry.ServName == "mra" {
				for _, ip := range entry.Ips {
					if strings.Contains(ip, ".") {
						pconn := new(PortsResult)
						pconn.Init(ip, entry.Fqdn, entry.ServName)
						for _, port := range []string{entry.Port, "5060", "5061", "5222"} {
							wg.Add(1)
							go pconn.Run(ip, port, input)
						}
					}
				}
			} else {
				for _, ip := range entry.Ips {
					if strings.Contains(ip, ".") {
						pconn := new(PortsResult)
						pconn.Init(ip, entry.Fqdn, entry.ServName)
						wg.Add(1)
						go pconn.Run(ip, entry.Port, input)
					}
				}
			}
		}
	}
	wg.Wait()
	close(input)
}

