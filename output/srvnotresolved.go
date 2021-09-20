package output

import "strings"

func MakeUndiscoveredSrv(discoveredsrv []Srv, channel chan []DiscoveredSRVrecords) {
	undslist := make([]DiscoveredSRVrecords, 0)	
	for _, srv := range discoveredsrv {
		if strings.HasPrefix(srv.Cname, "_cisco-uds") || strings.HasPrefix(srv.Cname, "_cuplogin") {
			discoveredSRVrecords := new(DiscoveredSRVrecords)
			discoveredSRVrecords.Cname = srv.Cname
			discoveredSRVrecords.Service = srv.Service
			for _, fqdn := range srv.Fqdns {
				discFqdn := new(DiscoveredFqdn)
				discFqdn.Name = fqdn.Name
				discoveredSRVrecords.Fqdns = append(discoveredSRVrecords.Fqdns, *discFqdn)
			}
			undslist = append(undslist, *discoveredSRVrecords)
		}
	}
	channel <- undslist
}
	