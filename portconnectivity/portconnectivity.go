package portconnectivity

import (
	"log"
	"net"
	"srvchecker/srv"
	"strings"
	"time"
)



type PortsResult struct {
	Ip    string
	Fqdn  string
	ports map[string]bool
}

type PortsResults []PortsResult


func (p *PortsResults) Connectivity(srvresults srv.SRVResults){
	for k,v := range srvresults {
		if k=="mra" {
			for _, entry := range v {
				for _, ip := range entry.Ips {
					if strings.Contains(ip, ".") {
						pconn := new(PortsResult)
						pconn.Ip = ip
						pconn.Fqdn = entry.Fqdn
						for _, port := range []string{entry.Port, "5060", "5061", "5222"} {
							log.Println(ip,":",port)

							timeout := time.Second
							conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
							
							if err != nil {
								pconn.ports = map[string]bool{port:false}
								
							}
							if conn != nil {
								defer conn.Close()
								pconn.ports = map[string]bool{port:true}
							}


						}
						*p = append((*p), *pconn)
					}
				}
			}
		} else {
			for _, entry := range v {
				for _, ip := range entry.Ips {
					if strings.Contains(ip, ".") {
						pconn := new(PortsResult)
						pconn.Ip = ip
						pconn.Fqdn = entry.Fqdn
						timeout := time.Second
						conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, entry.Port), timeout)
						
						if err != nil {
							pconn.ports = map[string]bool{entry.Port:false}
							
						}
						if conn != nil {
							defer conn.Close()
							pconn.ports = map[string]bool{entry.Port:true}
						}
						log.Println(ip,":",entry.Port)
						*p = append((*p), *pconn)
					}
				}
			}
		}
	}
}

