package output

import "strings"

// type DiscoveredPort struct {
// 	Fqdn string
// 	Ip string
// 	Ports []map[string]bool
// }

type DiscoveredPort map[string]map[string]map[string]string

func MakeTcpConnectivity(discoveredsrv []Srv, channel chan DiscoveredPort){
	
	discoveredPort := make(DiscoveredPort)
	for _, srv := range discoveredsrv {
		for _, fqdn := range srv.Fqdns {
			if strings.Contains(fqdn.Name, ".") {
				if _, ok := discoveredPort[fqdn.Name]; !ok {
					discoveredPort[fqdn.Name] = make(map[string]map[string]string)
				}
				for _, ip := range fqdn.Ips {
					if ip.SrvPort.Proto == "tcp" {
						if _, ok := discoveredPort[fqdn.Name][ip.Ip]; !ok {
							discoveredPort[fqdn.Name][ip.Ip] = make(map[string]string)
						}
						discoveredPort[fqdn.Name][ip.Ip][ip.SrvPort.Num] = ip.SrvPort.Open
					}
				}
			}
		}
	}

	channel <- discoveredPort
}


