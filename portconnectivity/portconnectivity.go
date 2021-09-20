package portconnectivity

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"srvchecker/srv"
	"strings"
	"sync"
	"time"
)

var admin_known_ports = []string{"443", "80", "22", "7001", "2222"}
var turn_ports = []string{"443:tcp", "3478:tcp", "3478:udp"}


type PortsResult struct {
	Ip    		string
	Fqdn  		string
	Port 		map[string]bool
	Udp 		bool 
	ServName 	string
	Certs 	    []*x509.Certificate
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
		p.Port = map[string]bool{port:false}	
	}
	if conn != nil {
		defer conn.Close()
		p.Port = map[string]bool{port:true}
		if (port == "8443" || port== "5061") {
			p.GetCert(ip , port)
		}
	}
	result <- *p
}

func (p *PortsResult) GetCert(ip string, port string) {
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
		p.Certs = conn.ConnectionState().PeerCertificates
	}
    
    // for _, cert := range *p.Certs {
    //     log.Printf("Issuer Name: %s\n", cert.Issuer)
    //     log.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
    //     log.Printf("Common Name: %s \n", cert.Issuer.CommonName)
    // }
}

func (p *PortsResult) RunTurnCheck(ip string, port string, udp bool, result chan PortsResult) {
	allocation_request := []byte{0, 3, 0, 36, 33, 18, 164, 66, 153, 147, 70, 130, 126, 38, 40, 41, 228, 206, 31, 174, 0, 25, 0, 4, 17, 0, 0, 0, 0, 13, 0, 4, 0, 0, 2, 88, 128, 34, 0, 5, 65, 99, 97, 110, 111, 0, 0, 0, 0, 23, 0, 4, 1, 0, 0, 0}
    buf := make([]byte, 16)
	
	if udp {
		conn, err := net.DialTimeout("udp", ip+":"+port, 1 * time.Second)
		if err != nil {
			p.Port = map[string]bool{port:false}
			p.Udp = true
		} else {
			conn.Write(allocation_request)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			conn.Read(buf)
			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
				p.Port = map[string]bool{port:true}
				p.Udp = true
			} else {
				p.Port = map[string]bool{port:false}
				p.Udp = true
			}
			conn.Close()
		}
		
		result <- *p
	} else {
		var err error

		conn, err := net.DialTimeout("tcp", ip+":"+port, 1 * time.Second)
		if err != nil {
			p.Port = map[string]bool{port:false}
		} else {
			defer conn.Close()
			conn.Write(allocation_request)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			conn.Read(buf)
			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
				p.Port = map[string]bool{port:true}
			} else {
				p.Port = map[string]bool{port:false}
			}
		}
		result <- *p
	}
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
						for _, port := range []string{entry.Port, "5061", "5222"} {
							pconn := new(PortsResult)
							pconn.Init(ip, entry.Fqdn, entry.ServName)
							wg.Add(1)
							go pconn.Run(ip, port, input)
						}
						for _, port := range admin_known_ports {
							pconn := new(PortsResult)
							pconn.Init(ip, entry.Fqdn, "admin")
							wg.Add(1)
							go pconn.Run(ip, port, input)
						}
						for _, turnport := range turn_ports {
							udp := false
							tl := strings.Split(turnport, ":")
							port := tl[0]
							if tl[1] == "udp" {
								udp = true
							}
							pconn := new(PortsResult)
							pconn.Init(ip, entry.Fqdn, "turn")
							wg.Add(1)
							go pconn.RunTurnCheck(ip, port, udp, input)
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
						for _, port := range admin_known_ports {
							pconn := new(PortsResult)
							pconn.Init(ip, entry.Fqdn, "admin")
							wg.Add(1)
							go pconn.Run(ip, port, input)
						}
						for _, turnport := range turn_ports {
							udp := false
							tl := strings.Split(turnport, ":")
							port := tl[0]
							if tl[1] == "udp" {
								udp = true
							}
							pconn := new(PortsResult)
							pconn.Init(ip, entry.Fqdn, "turn")
							wg.Add(1)
							go pconn.RunTurnCheck(ip, port, udp, input)
						}
					}
				}
			}
		}
	}
	wg.Wait()
	close(input)
}

