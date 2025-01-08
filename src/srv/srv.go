package srv

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

// DiscoveredSrvRow represents a single SRV record and associated details
type DiscoveredSrvRow struct {
	Srv         string
	Fqdn        string
	Ip          string
	Priority    string
	Weight      string
	Port        uint16
	Proto       string
	IsOpened    bool
	CertsChain
	ServiceName string
}

// DiscoveredSrvTable is a collection of DiscoveredSrvRow
type DiscoveredSrvTable []*DiscoveredSrvRow

// inputSRV represents a parsed SRV input entry
type inputSRV struct {
	service  string
	proto    string
	domain   string
	servName string
}

// inputSRVlist is a list of inputSRV entries
type inputSRVlist []inputSRV

// SRVTextList contains predefined SRV records by service type
var SRVTextList = map[string][]string{
	"mra":          {"_collab-edge._tls", "_cuplogin._tcp", "_cisco-uds._tcp"},
	"b2b":          {"_h323cs._tcp", "_sip._tcp", "_sips._tcp", "_sip._udp", "_h323ls._udp"},
	"xmpp_fed":     {"_xmpp-server._tcp"},
	"cma":          {"_xmpp-client._tcp"},
	"spark":        {"_sips._tcp.sipmtls"},
	"mssip":        {"_sipfederationtls._tcp"},
	"webexmessage": {"_webexconnect._tcp"},
	"mail":         {"_autodiscover._tcp", "_smtp._tcp", "_imaps._tcp", "_pop3s._tcp", "_submission._tcp"},
	"ftps":         {"_ftps._tcp"},
}

// init initializes the inputSRVlist for a given domain
func (s *inputSRVlist) init(domain string) {
	*s = make(inputSRVlist, 0, calculateTotalEntries())

	for serviceName, srvEntries := range SRVTextList {
		for _, srvEntry := range srvEntries {
			parts := strings.Split(srvEntry, ".")
			if len(parts) < 2 {
				continue
			}
			*s = append(*s, inputSRV{
				service:  strings.TrimPrefix(parts[0], "_"),
				proto:    strings.TrimPrefix(parts[1], "_"),
				domain:   domain,
				servName: serviceName,
			})
		}
	}
}

// calculateTotalEntries calculates the total number of SRV entries for preallocation
func calculateTotalEntries() int {
	total := 0
	for _, entries := range SRVTextList {
		total += len(entries)
	}
	return total
}

// init initializes a DiscoveredSrvRow
func (d *DiscoveredSrvRow) init(cname, servName, priority, weight, fqdn string, port uint16, ip, proto string) {
	d.Srv = cname
	d.ServiceName = servName
	d.Priority = priority
	d.Weight = weight
	d.Fqdn = fqdn
	d.Port = port
	d.Ip = ip
	d.Proto = proto
}

// fetchIps resolves IPs for an SRV record and updates the DiscoveredSrvTable
func (d *DiscoveredSrvTable) fetchIps(servName, cname string, fqdn *net.SRV, proto string, wg *sync.WaitGroup) {
	defer wg.Done()

	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", fqdn.Target)
	if err != nil {
		*d = append(*d, createDiscoveredSrvRow(cname, servName, fqdn, "A record not configured", proto))
		return
	}

	for _, ip := range ips {
		row := createDiscoveredSrvRow(cname, servName, fqdn, ip.To4().String(), proto)
		if proto == "tcp" {
			row.Connect_cert(ip.To4().String(), fmt.Sprint(fqdn.Port))
		}
		*d = append(*d, row)
	}
}

// createDiscoveredSrvRow creates a new DiscoveredSrvRow
func createDiscoveredSrvRow(cname, servName string, fqdn *net.SRV, ip, proto string) *DiscoveredSrvRow {
	return &DiscoveredSrvRow{
		Srv:         cname,
		ServiceName: servName,
		Priority:    fmt.Sprint(fqdn.Priority),
		Weight:      fmt.Sprint(fqdn.Weight),
		Fqdn:        fqdn.Target,
		Port:        fqdn.Port,
		Ip:          ip,
		Proto:       proto,
	}
}

// ForDomain discovers SRV records and populates the DiscoveredSrvTable
func (s *DiscoveredSrvTable) ForDomain(domain string) {
	mysrvs := inputSRVlist{}
	mysrvs.init(domain)

	var wg sync.WaitGroup

	for _, srv := range mysrvs {
		proto := "udp"
		if strings.HasPrefix(srv.proto, "t") {
			proto = "tcp"
		}
		cname := fmt.Sprintf("_%s._%s.%s", srv.service, srv.proto, srv.domain)
		_, fqdns, err := net.LookupSRV(srv.service, srv.proto, srv.domain)
		if err != nil {
			*s = append(*s, &DiscoveredSrvRow{
				ServiceName: srv.servName,
				Srv:         cname,
				Fqdn:        "SRV record not configured",
			})
			continue
		}

		for _, fqdn := range fqdns {
			wg.Add(1)
			go s.fetchIps(srv.servName, cname, fqdn, proto, &wg)
		}
	}

	wg.Wait()

	// Sort the table by SRV name
	sort.Slice(*s, func(i, j int) bool {
		return (*s)[i].Srv <= (*s)[j].Srv
	})
}
