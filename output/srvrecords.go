package output

import (
	"fmt"
	"strings"
)



type DiscoveredCert struct {
	CN string
	ExtendedKeyUsage string
	Issuer string
	KeyUsage string
	SAN string
	Subject string
	NotAfter string
	NotBefore string
	Txt string
}

type DiscoveredIp struct {
	Ip string
	Priority string
	Weight string
	PortNum string
	PortOpen string
	PortType string
	Cert DiscoveredCert
}

type DiscoveredFqdn struct {
	Name string
	Ips []DiscoveredIp
}

type DiscoveredSRVrecords struct {
	Cname string
	Service string
	Fqdns []DiscoveredFqdn 
}




func MakeDiscoveredSRVrecordsMap(discoveredsrv []Srv) map[string]DiscoveredSRVrecords {
	DiscoveredSRVrecordsMap := make(map[string]DiscoveredSRVrecords)

	for _, srv := range discoveredsrv {
		discoveredSRVrecords := new(DiscoveredSRVrecords)
		fmt.Println(srv.Service, srv.Cname)
		discoveredSRVrecords.Cname = srv.Cname
		discoveredSRVrecords.Service = srv.Service
		for _, fqdn := range srv.Fqdns {
			fmt.Println(fqdn.Name)
			discFqdn := new(DiscoveredFqdn)
			discFqdn.Name = fqdn.Name
			for _, ip := range fqdn.Ips {
				discIp := new(DiscoveredIp)
				if len(ip.SrvPort.Certs) > 0 {
					fmt.Println(ip.Ip, ip.Priority , ip.Weight, ip.SrvPort.Num, ip.SrvPort.Open, "Cert:", ip.SrvPort.Certs[0].Subject.CommonName)
					discIp.Ip = ip.Ip
					discIp.Priority = ip.Priority
					discIp.Weight = ip.Weight
					discIp.PortNum = ip.SrvPort.Num
					discIp.PortOpen = ip.SrvPort.Open
					discIp.PortType = ip.SrvPort.Proto
					
					discIp.Cert = DiscoveredCert{CN: ip.SrvPort.Certs[0].Subject.CommonName,
						Issuer: ip.SrvPort.Certs[0].Issuer.CommonName,
						SAN: strings.Join(ip.SrvPort.Certs[0].DNSNames, ", "),
						NotAfter: ip.SrvPort.Certs[0].NotAfter.String(),
						NotBefore: ip.SrvPort.Certs[0].NotBefore.String(),
					}
				} else {
					discIp.Ip = ip.Ip
					discIp.Priority = ip.Priority
					discIp.Weight = ip.Weight
					discIp.PortNum = ip.SrvPort.Num
					discIp.PortOpen = ip.SrvPort.Open
					discIp.PortType = ip.SrvPort.Proto
					
					fmt.Println(ip.Ip, ip.Priority , ip.Weight, ip.SrvPort.Num, ip.SrvPort.Open)
				}
				discFqdn.Ips = append(discFqdn.Ips, *discIp)
			}
			discoveredSRVrecords.Fqdns = append(discoveredSRVrecords.Fqdns, *discFqdn)
		}
		DiscoveredSRVrecordsMap[srv.Cname] = *discoveredSRVrecords
	}
	return DiscoveredSRVrecordsMap
	
}
