package portconnectivity

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lieranderl/srvchecker/srv"
)

// Constants for predefined ports
var predefinedPorts = map[string][]string{
	"admin":     {"443", "80", "22"},
	"traversal": {"7001", "2222"},
	"turn":      {"443:tcp", "3478:tcp", "3478:udp"},
	"mra":       {"5061", "5222"},
}

// Port represents a network port
type Port struct {
	IsOpened    bool
	Num         uint16
	Proto       string
	Type        string
	ServiceName string
}

// Ports is a collection of Port pointers
type Ports []*Port

// TcpConnectivityRow represents a row in the connectivity table
type TcpConnectivityRow struct {
	Fqdn  string
	Ip    string
	Ports Ports
}

// TcpConnectivityTable is a collection of connectivity rows
type TcpConnectivityTable []*TcpConnectivityRow

// FetchFromSrv populates the connectivity table based on SRV records
func (t *TcpConnectivityTable) FetchFromSrv(srvres srv.DiscoveredSrvTable) *TcpConnectivityTable {
	ipExists := make(map[string]*TcpConnectivityRow)

	// Parse predefined ports
	parsedPorts := parsePredefinedPorts()

	for _, srv := range srvres {
		if !strings.Contains(srv.Ip, ".") {
			continue
		}

		// Ensure row exists for each IP
		tcpConnectivityRow := ensureRowExists(t, ipExists, srv.Fqdn, srv.Ip)

		// Add service-specific ports
		if srv.Proto == "tcp" {
			tcpConnectivityRow.Ports.AddPort(srv.Port, "service", "tcp", srv.ServiceName)
		}
		if srv.ServiceName == "mra" {
			for _, port := range predefinedPorts["mra"] {
				tcpConnectivityRow.Ports.AddPort(parsePort(port), "service", "tcp", srv.ServiceName)
			}
		}

		// Add predefined ports
		for _, ports := range parsedPorts {
			for _, port := range ports {
				tcpConnectivityRow.Ports.AddPort(port.Num, port.Type, port.Proto, srv.ServiceName)
			}
		}
	}

	// Sort the table by FQDN
	sort.Slice(*t, func(i, j int) bool {
		return (*t)[i].Fqdn < (*t)[j].Fqdn
	})
	return t
}

// Connectivity checks the connectivity for each port
func (t *TcpConnectivityTable) Connectivity() {
	var wg sync.WaitGroup
	for _, row := range *t {
		for _, port := range row.Ports {
			wg.Add(1)
			go port.checkConnection(row.Ip, &wg)
		}
	}
	wg.Wait()
}

// AddPort adds a port to the Ports collection if not already present
func (p *Ports) AddPort(num uint16, portType, proto, serviceName string) {
	for _, existingPort := range *p {
		if existingPort.Num == num && existingPort.Proto == proto && existingPort.Type == portType {
			return
		}
	}
	*p = append(*p, &Port{Num: num, IsOpened: false, Type: portType, Proto: proto, ServiceName: serviceName})
}

// checkConnection checks if a port is open
func (p *Port) checkConnection(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	if p.Type == "turn" {
		p.IsOpened = checkTurnConnection(ip, fmt.Sprint(p.Num))
	} else {
		p.IsOpened = checkTcpConnection(ip, fmt.Sprint(p.Num))
	}
}

// parsePredefinedPorts parses the predefined ports into structured data
func parsePredefinedPorts() map[string][]*Port {
	result := make(map[string][]*Port)
	for portType, portStrings := range predefinedPorts {
		var ports []*Port
		for _, portStr := range portStrings {
			pp := strings.Split(portStr, ":")
			proto := "tcp"
			if len(pp) > 1 {
				proto = pp[1]
			}
			ports = append(ports, &Port{Num: parsePort(pp[0]), IsOpened: false, Type: portType, Proto: proto})
		}
		result[portType] = ports
	}
	return result
}

// ensureRowExists ensures that a row exists for the given IP
func ensureRowExists(t *TcpConnectivityTable, ipExists map[string]*TcpConnectivityRow, fqdn, ip string) *TcpConnectivityRow {
	if row, exists := ipExists[ip]; exists {
		return row
	}
	row := &TcpConnectivityRow{Fqdn: fqdn, Ip: ip, Ports: Ports{}}
	*t = append(*t, row)
	ipExists[ip] = row
	return row
}

// parsePort safely parses a port number
func parsePort(portStr string) uint16 {
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		fmt.Printf("Invalid port number: %v\n", err) // Replace with proper logging
		return 0
	}
	return uint16(port)
}

// checkTcpConnection checks TCP connectivity to a port
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

// checkTurnConnection checks TURN connectivity to a port
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
