package output

import (
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

func (discIp *DiscoveredIp) Init(ip Ip, discoveredcert []Cert, ) {
	discIp.Ip = ip.Ip
	discIp.Priority = ip.Priority
	discIp.Weight = ip.Weight
	discIp.PortNum = ip.SrvPort.Num
	discIp.PortOpen = ip.SrvPort.Open
	discIp.PortType = ip.SrvPort.Proto

	for _, cert := range discoveredcert {
		if (cert.Ip == ip.Ip && cert.Port == ip.SrvPort.Num) {
			if len(cert.Certs) > 0 {
				discIp.Cert = DiscoveredCert{CN: cert.Certs[0].Subject.CommonName,
					Issuer: cert.Certs[0].Issuer.CommonName,
					SAN: strings.Join(cert.Certs[0].DNSNames, ", "),
					NotAfter: cert.Certs[0].NotAfter.String(),
					NotBefore: cert.Certs[0].NotBefore.String(),
				}
			}
		}
	}
}

type DiscoveredFqdn struct {
	Name string
	Service string
	Ips []DiscoveredIp
}

type DiscoveredSRVrecords struct {
	Cname string
	Service string
	Fqdns []DiscoveredFqdn 
}


func MakeDiscoveredSRVrecordsMap(discoveredsrv []Srv, discoveredcert []Cert, channel chan []DiscoveredSRVrecords)  {
	dslist := make([]DiscoveredSRVrecords, 0)

	for _, srv := range discoveredsrv {
		if !(strings.HasPrefix(srv.Cname, "_cisco-uds") || strings.HasPrefix(srv.Cname, "_cuplogin")) {
			discoveredSRVrecords := new(DiscoveredSRVrecords)
			discoveredSRVrecords.Cname = srv.Cname
			discoveredSRVrecords.Service = srv.Service
			for _, fqdn := range srv.Fqdns {
				discFqdn := new(DiscoveredFqdn)
				discFqdn.Name = fqdn.Name
				for _, ip := range fqdn.Ips {
					discIp := new(DiscoveredIp)
					discIp.Init(ip, discoveredcert)
					discFqdn.Ips = append(discFqdn.Ips, *discIp)
				}
				discoveredSRVrecords.Fqdns = append(discoveredSRVrecords.Fqdns, *discFqdn)
			}
			dslist = append(dslist, *discoveredSRVrecords) 
		}
	}
	channel <- dslist
}
