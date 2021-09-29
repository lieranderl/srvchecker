package portconnectivity

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"srvchecker/srv"
	"strings"
	"sync"
	"time"
)

var admin_known_ports = []string{"443", "80", "22", "7001", "2222"}
var turn_ports = []string{"443:tcp", "3478:tcp", "3478:udp"}
var mutex = &sync.RWMutex{}
type Sname string
type Fqdn string
type Ip string
type Port struct {
	IsOpened bool
	Num string
	Proto string
}

type TcpConnectivityRow struct {
	ServiceName string
	Fqdn string
	Ip string
	AdminPorts []Port
	TraversalPorts []Port
}

type TurnConnectivityRow struct {
	ServiceName string
	Fqdn string
	Ip string
	TcpPorts []Port
	UdpPorts []Port
}


type TcpConnectivityTable []TcpConnectivityRow
type TurnConnectivityTable []TurnConnectivityRow




type PortsResults map[Fqdn]map[Ip]map[string]*Port

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}


func FetchFromSrv(srvres srv.DiscoveredSrvTable)  {
	
	var wg sync.WaitGroup

	// fqdn_buffer := make([]string,0)
	//Service port check
	for _, srv := range srvres {
		if (strings.Contains(srv.Fqdn, ".") && strings.Contains(srv.Ip, ".") && srv.Proto == "tcp") { 
			wg.Add(1)
			go srv.Connect_cert(srv.Ip, srv.Port, &wg)
		}
	}
	wg.Wait()
}


// func (p *PortsResults)FetchPorts(fqdn srv.Fqdn, ip srv.Ip, port *srv.Port, sname string, wg *sync.WaitGroup){
// 	defer wg.Done()
// 	current_fqdn := (*p)[Fqdn(fqdn)]
// 	currrent_ip := current_fqdn[Ip(ip)]

// 	if _, ok := (*p)[Fqdn(fqdn)]; !ok {
// 		current_fqdn = make(map[Ip]map[string]*Port)
// 	}
// 	if _, ok := current_fqdn[Ip(ip)] ; !ok {
// 		currrent_ip = make(map[string]*Port)
// 	}
// 	for potuniq, pp := range *port {
// 		if _, ok := currrent_ip[string(potuniq)] ; !ok {
// 			currrent_ip[string(potuniq)] = new(Port)
// 		}
// 		currrent_ip[string(potuniq)].IsOpened = pp.IsOpen
// 		currrent_ip[string(potuniq)].Sname = sname
// 	}

// 	if sname == "mra" {
// 		for _, port := range []string{"5061", "5222"} {
// 			if _, ok := currrent_ip[port+":tcp"] ; !ok {
// 				currrent_ip[port+":tcp"] = new(Port)
// 			}
// 			currrent_ip[port+":tcp"].IsOpened = CheckConnection(string(ip), port)
// 			currrent_ip[port+":tcp"].Sname = sname

// 		}
// 	}
// 	for _, port := range admin_known_ports {
// 		if _, ok := currrent_ip[port+":tcp"] ; !ok {
// 			currrent_ip[port+":tcp"] = new(Port)
// 		}
// 		currrent_ip[port+":tcp"].IsOpened = CheckConnection(string(ip), port)
// 		currrent_ip[port+":tcp"].Sname = sname
// 	}
// 	for _, port := range turn_ports {
// 		if _, ok := currrent_ip[port+"turn"] ; !ok {
// 			currrent_ip[port+"turn"] = new(Port)
// 		}
// 		currrent_ip[port+"turn"].IsOpened = RunTurnCheck(string(ip), port)
// 		currrent_ip[port+"turn"].Sname = sname
// 	}

// 	mutex.Lock()
// 	current_fqdn[Ip(ip)] = currrent_ip
// 	(*p)[Fqdn(fqdn)]= current_fqdn
// 	mutex.Unlock()

// }

// func CheckConnection(ip string, port string) bool {
// 	timeout := time.Second
// 	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
// 	if err != nil {
// 		return false	
// 	} else {
// 		defer conn.Close()
// 		return true
// 	}
// }

// func RunTurnCheck(ip string, port string) bool {
// 	allocation_request := []byte{0, 3, 0, 36, 33, 18, 164, 66, 153, 147, 70, 130, 126, 38, 40, 41, 228, 206, 31, 174, 0, 25, 0, 4, 17, 0, 0, 0, 0, 13, 0, 4, 0, 0, 2, 88, 128, 34, 0, 5, 65, 99, 97, 110, 111, 0, 0, 0, 0, 23, 0, 4, 1, 0, 0, 0}
//     buf := make([]byte, 16)
	
// 	pp := strings.Split(port, ":")
// 	port = pp[0]
// 	proto := pp[1]

// 	if proto == "udp" {
// 		conn, err := net.DialTimeout("udp", ip+":"+port, 1 * time.Second)
// 		if err != nil {
// 			return false
// 		} else {
// 			defer conn.Close()
// 			conn.Write(allocation_request)
// 			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
// 			conn.Read(buf)
// 			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
// 				return true
// 			} else {
// 				return false
// 			}
			
// 		}
		
// 	} else {
// 		var err error
// 		conn, err := net.DialTimeout("tcp", ip+":"+port, 1 * time.Second)
// 		if err != nil {
// 			return false
// 		} else {
// 			defer conn.Close()
// 			conn.Write(allocation_request)
// 			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
// 			conn.Read(buf)
// 			if bytes.HasPrefix(buf, []byte{1, 19, 0,160}) {
// 				return true
// 			} else {
// 				return false
// 			}
// 		}
// 	}
// }