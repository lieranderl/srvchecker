package portconnectivity

import (
	// "bytes"
	// "crypto/tls"

	// "log"
	// "net"
	"net"
	"srvchecker/srv"
	"strings"
	"sync"
	"time"
	// "strings"
	// "sync"
	// "time"
)

var admin_known_ports = []string{"443", "80", "22", "7001", "2222"}
var turn_ports = []string{"443:tcp", "3478:tcp", "3478:udp"}


// type PortsResult struct {
// 	Sname 		string
// 	Fqdn		map[string]
// }

// type 
// Ip    		string
// Port 		map[string]bool
// Proto		string
// Certs 	    []*x509.Certificate



type Sname string
type Fqdn string
type Ip string
type Port struct {
	IsOpened bool
	Sname string
}

type PortsResults map[Fqdn]map[Ip]map[string]*Port

func (p *PortsResults)fetchFromSrvResults(srvres *srv.SRVResults) {
	var wg sync.WaitGroup
	(*p) = make(map[Fqdn]map[Ip]map[string]*Port)
	for _, srvresult := range *srvres {
		for fqdn, ips := range srvresult.Fqdn {
			if strings.Contains(string(fqdn), ".") {
				for ip, port := range ips.Ips {
					if strings.Contains(string(ip), ".") {
						wg.Add(1)
						go p.FetchPorts(fqdn, ip, port, srvresult.Sname, &wg)

					// 	if _, ok := (*p)[Fqdn(fqdn)]; !ok {
					// 		(*p)[Fqdn(fqdn)] = make(map[Ip]map[string]*Port)
					// 	}
					// 	if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)] ; !ok {
					// 		(*p)[Fqdn(fqdn)][Ip(ip)] = make(map[string]*Port)
					// 	}
					// 	for potuniq, pp := range *port {
					// 		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)] ; !ok {
					// 			(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)] = new(Port)
					// 		}
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)].IsOpened = pp.IsOpen
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)].Sname = srvresult.Sname
					// 	}
						
					// 	if srvresult.Sname == "mra" {
					// 		for _, port := range []string{"5061", "5222"} {
					// 			if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port] ; !ok {
					// 				(*p)[Fqdn(fqdn)][Ip(ip)][port] = new(Port)
					// 			}
					// 			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].IsOpened = false
					// 			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].Sname = srvresult.Sname
	
					// 		}
					// 	}
					// 	for _, port := range admin_known_ports {
					// 		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"] ; !ok {
					// 			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"] = new(Port)
					// 		}
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].IsOpened = false
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].Sname = srvresult.Sname
					// 	}
					// 	for _, port := range turn_ports {
					// 		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port] ; !ok {
					// 			(*p)[Fqdn(fqdn)][Ip(ip)][port] = new(Port)
					// 		}
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][port].IsOpened = false
					// 		(*p)[Fqdn(fqdn)][Ip(ip)][port].Sname = srvresult.Sname
					// 	}
					}
				}
			}
			
		}
	}
	wg.Wait()
}


// func (p *PortsResult) Init(ip string, fqdn string, servname string) {
// 	p.Ip = ip
// 	p.Fqdn = fqdn
// 	p.Sname = servname
// }

func (p *PortsResults)FetchPorts(fqdn srv.Fqdn, ip srv.Ip, port *srv.Port, sname string, wg *sync.WaitGroup){
	
	defer wg.Done()
	if _, ok := (*p)[Fqdn(fqdn)]; !ok {
		(*p)[Fqdn(fqdn)] = make(map[Ip]map[string]*Port)
	}
	if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)] ; !ok {
		(*p)[Fqdn(fqdn)][Ip(ip)] = make(map[string]*Port)
	}
	for potuniq, pp := range *port {
		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)] ; !ok {
			(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)] = new(Port)
		}
		(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)].IsOpened = pp.IsOpen
		(*p)[Fqdn(fqdn)][Ip(ip)][string(potuniq)].Sname = sname
	}
	
	if sname == "mra" {
		for _, port := range []string{"5061", "5222"} {
			if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port] ; !ok {
				(*p)[Fqdn(fqdn)][Ip(ip)][port] = new(Port)
			}
			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].IsOpened = CheckConnection(string(ip), port)
			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].Sname = sname

		}
	}
	for _, port := range admin_known_ports {
		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"] ; !ok {
			(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"] = new(Port)
		}
		(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].IsOpened = CheckConnection(string(ip), port)
		(*p)[Fqdn(fqdn)][Ip(ip)][port+":tcp"].Sname = sname
	}
	for _, port := range turn_ports {
		if _, ok := (*p)[Fqdn(fqdn)][Ip(ip)][port] ; !ok {
			(*p)[Fqdn(fqdn)][Ip(ip)][port] = new(Port)
		}
		(*p)[Fqdn(fqdn)][Ip(ip)][port].IsOpened = false
		(*p)[Fqdn(fqdn)][Ip(ip)][port].Sname = sname
	}
}



func CheckConnection(ip string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false	
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

// func (p *PortsResult) GetCert(ip string, port string) {
// 	conf := &tls.Config{
//         InsecureSkipVerify: true,
//     }

//     conn, err := tls.Dial("tcp", ip+":"+port, conf)
// 	conn.SetDeadline(time.Now().Add(2 * time.Second))
//     if err != nil {
//         log.Println("Error in Dial", err)

//     }
// 	if conn != nil {
// 		defer conn.Close()
// 		p.Certs = conn.ConnectionState().PeerCertificates
// 	}
    
//     // for _, cert := range *p.Certs {
//     //     log.Printf("Issuer Name: %s\n", cert.Issuer)
//     //     log.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
//     //     log.Printf("Common Name: %s \n", cert.Issuer.CommonName)
//     // }
// }

// func (p *PortsResult) RunTurnCheck(ip string, port string, udp bool, result chan PortsResult) {
// 	allocation_request := []byte{0, 3, 0, 36, 33, 18, 164, 66, 153, 147, 70, 130, 126, 38, 40, 41, 228, 206, 31, 174, 0, 25, 0, 4, 17, 0, 0, 0, 0, 13, 0, 4, 0, 0, 2, 88, 128, 34, 0, 5, 65, 99, 97, 110, 111, 0, 0, 0, 0, 23, 0, 4, 1, 0, 0, 0}
//     buf := make([]byte, 16)
	
// 	if udp {
// 		conn, err := net.DialTimeout("udp", ip+":"+port, 1 * time.Second)
// 		if err != nil {
// 			p.Port = map[string]bool{port:false}
// 			p.Udp = true
// 		} else {
// 			conn.Write(allocation_request)
// 			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
// 			conn.Read(buf)
// 			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
// 				p.Port = map[string]bool{port:true}
// 				p.Udp = true
// 			} else {
// 				p.Port = map[string]bool{port:false}
// 				p.Udp = true
// 			}
// 			conn.Close()
// 		}
		
// 		result <- *p
// 	} else {
// 		var err error

// 		conn, err := net.DialTimeout("tcp", ip+":"+port, 1 * time.Second)
// 		if err != nil {
// 			p.Port = map[string]bool{port:false}
// 		} else {
// 			defer conn.Close()
// 			conn.Write(allocation_request)
// 			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
// 			conn.Read(buf)
// 			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
// 				p.Port = map[string]bool{port:true}
// 			} else {
// 				p.Port = map[string]bool{port:false}
// 			}
// 		}
// 		result <- *p
// 	}
// }


// // type PortsResults []PortsResult

// func (p *PortsResults)handleResults(input chan PortsResult, wg *sync.WaitGroup) {
// 	for result := range input {
// 		*p = append(*p, result)
// 		wg.Done()
// 	}
// }


// func (p *PortsResults) Connectivity(srvresults srv.SRVResults){

// 	input := make(chan PortsResult)
// 	var wg sync.WaitGroup
// 	go p.handleResults(input, &wg)

// 	for _,v := range srvresults {
// 		for _, entry := range v {
// 			if entry.ServName == "mra" {
// 				for _, ip := range entry.Ips {
// 					if strings.Contains(ip, ".") {
// 						for _, port := range []string{entry.Port, "5061", "5222"} {
// 							pconn := new(PortsResult)
// 							pconn.Init(ip, entry.Fqdn, entry.ServName)
// 							wg.Add(1)
// 							go pconn.Run(ip, port, entry.Proto, input)
							
// 						}
// 						for _, port := range admin_known_ports {
// 							pconn := new(PortsResult)
// 							pconn.Init(ip, entry.Fqdn, "admin")
// 							wg.Add(1)
// 							go pconn.Run(ip, port, "tcp", input)
// 						}
// 						for _, turnport := range turn_ports {
// 							udp := false
// 							tl := strings.Split(turnport, ":")
// 							port := tl[0]
// 							if tl[1] == "udp" {
// 								udp = true
// 							}
// 							pconn := new(PortsResult)
// 							pconn.Init(ip, entry.Fqdn, "turn")
// 							wg.Add(1)
// 							go pconn.RunTurnCheck(ip, port, udp, input)
// 						}
// 					}
// 				}
// 			} else {
// 				for _, ip := range entry.Ips {
// 					if strings.Contains(ip, ".") {
// 						pconn := new(PortsResult)
// 						pconn.Init(ip, entry.Fqdn, entry.ServName)
// 						wg.Add(1)
// 						go pconn.Run(ip, entry.Port, entry.Proto, input)
// 						for _, port := range admin_known_ports {
// 							pconn := new(PortsResult)
// 							pconn.Init(ip, entry.Fqdn, "admin")
// 							wg.Add(1)
// 							go pconn.Run(ip, port, "tcp", input)
// 						}
// 						for _, turnport := range turn_ports {
// 							udp := false
// 							tl := strings.Split(turnport, ":")
// 							port := tl[0]
// 							if tl[1] == "udp" {
// 								udp = true
// 							}
// 							pconn := new(PortsResult)
// 							pconn.Init(ip, entry.Fqdn, "turn")
// 							wg.Add(1)
// 							go pconn.RunTurnCheck(ip, port, udp, input)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	wg.Wait()
// 	close(input)
// }

