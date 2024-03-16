package portconnectivity

import (
	"bytes"
	"fmt"
	"net"
	"sort"

	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/lieranderl/srvchecker/srv"
)

var admin_known_ports = []string{"443", "80", "22"}
var traversal_ports = []string{"7001", "2222"}
var turn_ports = []string{"443:tcp", "3478:tcp", "3478:udp"}

type Port struct {
	IsOpened    bool
	Num         uint16
	Proto       string
	Type        string
	ServiceName string
}

type Ports []*Port

type TcpConnectivityRow struct {
	Fqdn string
	Ip   string
	Ports
}

type TcpConnectivityTable []*TcpConnectivityRow

// func containsPorts(ports []*Port, port uint16, proto, t, serv string) bool {
// 	for _, p := range ports {
// 		if p.Num == port && p.Proto == proto && p.Type == t && p.ServiceName == serv {
// 			return true
// 		}
// 	}
// 	return false
// }

// func containsTcpConnectivity(s TcpConnectivityTable, ip string) bool {
// 	for _, a := range s {
// 		if a.Ip == ip {
// 			return true
// 		}
// 	}
// 	return false
// }

func (t *TcpConnectivityTable) FetchFromSrv(srvres srv.DiscoveredSrvTable) *TcpConnectivityTable {
	ipExists := make(map[string]*TcpConnectivityRow)
	for _, tcpConnectivityRow := range *t {
		ipExists[tcpConnectivityRow.Ip] = tcpConnectivityRow
	}

	parsePort := func(portStr string) uint16 {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			fmt.Printf("Invalid port number: " + err.Error()) // Replace with proper error handling
			return 0
		}
		return uint16(port)
	}

	adminPorts := make([]*Port, len(admin_known_ports))
	for i, portStr := range admin_known_ports {
		adminPorts[i] = &Port{Num: parsePort(portStr), IsOpened: false, Type: "admin", Proto: "tcp"}
	}

	traversalPorts := make([]*Port, len(traversal_ports))
	for i, portStr := range traversal_ports {
		traversalPorts[i] = &Port{Num: parsePort(portStr), IsOpened: false, Type: "traversal", Proto: "tcp"}
	}

	turnPorts := make([]*Port, len(turn_ports))
	for i, portStr := range turn_ports {
		pp := strings.Split(portStr, ":")
		turnPorts[i] = &Port{Num: parsePort(pp[0]), IsOpened: false, Type: "turn", Proto: pp[1]}
	}

	for _, srv := range srvres {
		if !strings.Contains(srv.Ip, ".") {
			continue
		}

		tcpConnectivityRow, exists := ipExists[srv.Ip]
		if !exists {
			tcpConnectivityRow = &TcpConnectivityRow{Fqdn: srv.Fqdn, Ip: srv.Ip, Ports: Ports{}}
			*t = append(*t, tcpConnectivityRow)
			ipExists[srv.Ip] = tcpConnectivityRow
		}

		addPort := func(port uint16, portType, proto, serviceName string) {
			for _, existingPort := range tcpConnectivityRow.Ports {
				if existingPort.Num == port && existingPort.Proto == proto && existingPort.Type == portType {
					fmt.Println("Port already exists: ", port, portType, proto, serviceName)
					return
				}
			}
			tcpConnectivityRow.Ports = append(tcpConnectivityRow.Ports, &Port{Num: port, IsOpened: false, Type: portType, Proto: proto, ServiceName: serviceName})

			fmt.Println("Added port: ", port, portType, proto, serviceName)
		}

		if srv.Proto == "tcp" {
			addPort(srv.Port, "service", "tcp", srv.ServiceName)
		}

		if srv.ServiceName == "mra" {
			for _, port := range []string{"5061", "5222"} {
				addPort(parsePort(port), "service", "tcp", srv.ServiceName)
			}
		}

		for _, port := range adminPorts {
			addPort(port.Num, port.Type, port.Proto, srv.ServiceName)
		}
		for _, port := range traversalPorts {
			addPort(port.Num, port.Type, port.Proto, srv.ServiceName)
		}
		for _, port := range turnPorts {
			addPort(port.Num, port.Type, port.Proto, srv.ServiceName)
		}
	}

	// Consider sorting after all operations if necessary.
	sort.Slice((*t)[:], func(i, j int) bool {
		return (*t)[i].Fqdn < (*t)[j].Fqdn
	})
	return t
}

func (t *TcpConnectivityTable) Connectivity() {
	var wg sync.WaitGroup
	for _, i := range *t {
		for _, port := range i.Ports {
			wg.Add(1)
			go port.connection(i.Ip, &wg)
		}
	}
	wg.Wait()
}

func (p *Port) connection(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	if p.Type == "turn" {
		p.IsOpened = checkTurnConnection(ip, fmt.Sprint(p.Num))
	} else {
		p.IsOpened = checkTcpConnection(ip, fmt.Sprint(p.Num))
	}
}

func checkTcpConnection(ip string, port string) bool {
    timeout := time.Second
    address := net.JoinHostPort(ip, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

func checkTurnConnection(ip string, port string) bool {
    timeout := 2 * time.Second
    allocationRequest := []byte{
        0, 3, 0, 36, 33, 18, 164, 66, 153, 147, 70, 130, 126, 38, 40, 41, 228, 206, 31, 174,
        0, 25, 0, 4, 17, 0, 0, 0, 0, 13, 0, 4, 0, 0, 2, 88, 128, 34, 0, 5, 65, 99, 97, 110,
        111, 0, 0, 0, 0, 23, 0, 4, 1, 0, 0, 0,
    }
    expectedResponsePrefix := []byte{1, 19, 0}
    buf := make([]byte, 16)

    address := net.JoinHostPort(ip, port)
    conn, err := net.DialTimeout("udp", address, timeout)
    if err != nil {
        return false
    }
    defer conn.Close()

    if _, err = conn.Write(allocationRequest); err != nil {
        return false
    }

    if err = conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
        return false
    }

    if _, err = conn.Read(buf); err != nil {
        return false
    }

    return bytes.HasPrefix(buf, expectedResponsePrefix)
}
