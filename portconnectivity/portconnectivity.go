package portconnectivity

import (
	"bytes"
	"net"
	"sort"
	"srvchecker/srv"
	"strings"
	"time"

	"sync"
)

var admin_known_ports = []string{"443", "80", "22"}
var traversal_ports = []string{"7001", "2222"}
var turn_ports = []string{"443:tcp", "3478:tcp", "3478:udp"}


type Port struct {
	IsOpened bool
	Num string
	Proto string
	Type string
}

type TcpConnectivityRow struct {
	ServiceName string
	Fqdn string
	Ip string
	Ports []*Port
}



type TcpConnectivityTable []*TcpConnectivityRow



func containsPorts(s []*Port, port, proto, t string) bool {
    for _, a := range s {
        if (a.Num == port && a.Proto == proto && a.Type == t) {
            return true
        }
    }
    return false
}

func containsTcpConnectivity(s TcpConnectivityTable, ip string) bool {
    for _, a := range s {
        if (a.Ip == ip) {
            return true
        }
    }
    return false
}


func (t *TcpConnectivityTable)FetchFromSrv(srvres srv.DiscoveredSrvTable)  {
	
	// var wg sync.WaitGroup
	//Service port check
	for _, srv := range srvres {
		if strings.Contains(srv.Ip, ".") {
			if !containsTcpConnectivity(*t, srv.Ip) {
				tcpConnectivityRow := new(TcpConnectivityRow)
				tcpConnectivityRow.Fqdn = srv.Fqdn
				tcpConnectivityRow.Ip = srv.Ip
				
				if srv.Proto == "tcp" {
					if !containsPorts(tcpConnectivityRow.Ports, srv.Port, "tcp", "service") {
						tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: srv.Port, IsOpened: srv.IsOpened, Type: "service", Proto: "tcp"})
					}	
				}
			
				if srv.ServiceName == "mra" {
					for _, port := range []string{"5061", "5222"} {
						if !containsPorts(tcpConnectivityRow.Ports, port, "tcp", "service") {
							tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: port, IsOpened: false, Type: "service", Proto: "tcp"})
						}
					}
				}		
				
				for _, port := range admin_known_ports {
					if !containsPorts(tcpConnectivityRow.Ports, port, "tcp", "admin") {
						tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: port, IsOpened: false, Type: "admin", Proto: "tcp"})
					}
				}
				for _, port := range traversal_ports {
					if !containsPorts(tcpConnectivityRow.Ports, port, "tcp", "traversal") {
						tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: port, IsOpened: false, Type: "traversal", Proto: "tcp"})
					}
				}
				for _, port := range turn_ports {
					pp := strings.Split(port, ":")
					port = pp[0]
					proto := pp[1]
					if !containsPorts(tcpConnectivityRow.Ports, port, proto, "turn") {
						tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: port, IsOpened: false, Type: "turn", Proto: proto})
					}
				}
				*t = append(*t, tcpConnectivityRow)
			} else {
				if srv.Proto == "tcp" {
					for _, tcpConnectivityRow := range *t {
						if (tcpConnectivityRow.Ip == srv.Ip) {
							if !containsPorts(tcpConnectivityRow.Ports, srv.Port, "tcp", "service") {
								tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: srv.Port, IsOpened: srv.IsOpened, Type: "service", Proto: "tcp"})
							}	
						}
					}
				}
				
			}
		}
	}
	sort.Slice((*t)[:], func(i, j int) bool {
		return (*t)[i].Fqdn < (*t)[j].Fqdn
	})


	// wg.Wait()
}

func (t *TcpConnectivityTable)Connectivity()  {
	var wg sync.WaitGroup
	for _, i := range *t {
		for _, port := range i.Ports {
			wg.Add(1)
			go port.connection(i.Ip, &wg)
		}
	}
	wg.Wait()
}

func (p *Port)connection(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	if p.Type == "turn" {
		p.IsOpened = RunTurnCheck(ip, p.Num, p.Proto)
	} else {
		p.IsOpened = CheckConnection(ip, p.Num)
	}
}


func CheckConnection(ip string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false	
	} else {
		defer conn.Close()
		return true
	}
}

func RunTurnCheck(ip string, port string, proto string) bool {
	allocation_request := []byte{0, 3, 0, 36, 33, 18, 164, 66, 153, 147, 70, 130, 126, 38, 40, 41, 228, 206, 31, 174, 0, 25, 0, 4, 17, 0, 0, 0, 0, 13, 0, 4, 0, 0, 2, 88, 128, 34, 0, 5, 65, 99, 97, 110, 111, 0, 0, 0, 0, 23, 0, 4, 1, 0, 0, 0}
    buf := make([]byte, 16)
	
	if proto == "udp" {
		conn, err := net.DialTimeout("udp", ip+":"+port, 2 * time.Second)
		if err != nil {
			return false
		} else {
			defer conn.Close()
			conn.Write(allocation_request)
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			conn.Read(buf)
			if bytes.HasPrefix(buf, []byte{1, 19, 0}) {
				return true
			} else {
				return false
			}
			
		}
		
	} else {
		var err error
		conn, err := net.DialTimeout("tcp", ip+":"+port, 2 * time.Second)
		if err != nil {
			return false
		} else {
			defer conn.Close()
			conn.Write(allocation_request)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			conn.Read(buf)
			if bytes.HasPrefix(buf, []byte{1, 19, 0}) {
				return true
			} else {
				return false
			}
		}
	}
}