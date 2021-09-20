package output

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"strconv"
	"strings"
)

type Srv struct {
	Service string
	Cname string
	Fqdns []Fqdn
}

type Fqdn struct {
	Service string
	Name string
	Ips []Ip
}

type Ip struct {
	Service string
	Ip string
	Priority string
	Weight string
	SrvPort Port
	AdminPorts []Port
	AdditionalServicePorts []Port
	TurnPorts []Port
}

type Port struct {
	Service string
	Num string
	Open string
	Proto string
	Certs []*x509.Certificate
}


func checkportinlist(ports []Port, port map[string]bool) bool {
	for _, p := range ports {
		for k := range port {
			if p.Num == k {
				return true 
			}
		}
	}
	return false
}


func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}



func Output(srvresults *srv.SRVResults, portsresults *portconnectivity.PortsResults) {

	var discoveredsrv []Srv

	for cname, srvres:=range *srvresults {
		mysrv := new(Srv)
		mysrv.Cname = cname
		fqdns := make([]Fqdn,0)
		for _, sr:= range srvres {
			if sr.Cname == cname {
				mysrv.Service = sr.ServName
				myfqdn := new(Fqdn)
				ips := make([]Ip,0)
				
				for _, ip := range sr.Ips {
					myIp := new(Ip)
					myIp.Ip = ip
					myIp.Priority = sr.Priority
					myIp.Weight = sr.Weight

					if strings.Contains(ip, ".") {
						for _, pres := range *portsresults{
							myPort := new(Port)
							
							if pres.Ip == ip {
								if (pres.ServName == sr.ServName) {
									for p, v := range pres.Port {
										if p==sr.Port {
											myPort.Num = p
											myPort.Open = strconv.FormatBool(v)
											if pres.Udp {
												myPort.Proto = "UDP"
											} else {
												myPort.Proto = "TCP"
											}
											myPort.Certs = pres.Certs
											myPort.Service = pres.ServName
											myIp.SrvPort = *myPort	
											break
										} else {
											myPort.Num = p
											myPort.Open = strconv.FormatBool(v)
											if pres.Udp {
												myPort.Proto = "UDP"
											} else {
												myPort.Proto = "TCP"
											}
											myPort.Certs = pres.Certs
											myPort.Service = pres.ServName
											myIp.AdditionalServicePorts = append(myIp.AdditionalServicePorts, *myPort)
											break
										} 
									}

								}
								if (pres.ServName == "turn") {
									if !checkportinlist(myIp.TurnPorts, pres.Port) {
										for p, v := range pres.Port {
											myPort.Num = p
											myPort.Open = strconv.FormatBool(v)
											if pres.Udp {
												myPort.Proto = "UDP"
											} else {
												myPort.Proto = "TCP"
											}
											myPort.Certs = pres.Certs
											myPort.Service = pres.ServName
											myIp.TurnPorts = append(myIp.TurnPorts, *myPort)
											break
										}
									}
									
								}
								if (pres.ServName == "admin"){
									if !checkportinlist(myIp.AdminPorts, pres.Port) {
										for p, v := range pres.Port {
											myPort.Num = p
											myPort.Open = strconv.FormatBool(v)
											if pres.Udp {
												myPort.Proto = "UDP"
											} else {
												myPort.Proto = "TCP"
											}
											myPort.Certs = pres.Certs
											myPort.Service = pres.ServName
											myIp.AdminPorts = append(myIp.AdminPorts, *myPort)
											break
										}
									}
								}

							}
						}
						ips = append(ips, *myIp)
					}
					
				}

				myfqdn.Service = sr.ServName
				myfqdn.Name = sr.Fqdn
				myfqdn.Ips = ips

				fqdns = append(fqdns, *myfqdn)

			}
		}
		mysrv.Fqdns = fqdns
		discoveredsrv = append(discoveredsrv, *mysrv)
	}


	DiscoveredSRVrecordsMap := MakeDiscoveredSRVrecordsMap(discoveredsrv)

	// fmt.Println("=====================")
	// fmt.Println("SRV records that should not resolve")
	// for _, srv := range discoveredsrv {
	// 	if strings.HasPrefix(srv.Cname, "_cisco-uds") || strings.HasPrefix(srv.Cname, "_cuplogin") {
	// 		for _, fqdn := range srv.Fqdns {
	// 			if fqdn.Name != "SRV record not configured" {
	// 				fmt.Println(fqdn.Service, srv.Cname, "Pizda")
	// 			} else {
	// 				fmt.Println(fqdn.Service, srv.Cname, "Not resolvable")
	// 			}
	// 		}
	// 	}
	// }

	// fmt.Println("===================")
	// fmt.Println("TCP Connectivity:")
	// fqdn_list := make([]string, 0)
	// for _, srv := range discoveredsrv {
	// 	for _, fqdn := range srv.Fqdns {
	// 		for _, ip := range fqdn.Ips {
	// 			if !stringInSlice(fqdn.Name, fqdn_list) {
	// 				fmt.Println(fqdn.Service, ip.Service, fqdn.Name, ip.Ip) 
	// 				fmt.Println(ip.SrvPort.Num, ip.SrvPort.Open)
	// 				for _, p := range ip.AdditionalServicePorts {
	// 					fmt.Println(p.Num, p.Open)
	// 				}
					
	// 			}
	// 		}
	// 		fqdn_list = append(fqdn_list, fqdn.Name)
	// 	}
	// }

	// fmt.Println("===================")

	// fmt.Println("Admin ports:")
	// fqdn_list = make([]string, 0)
	// for _, srv := range discoveredsrv {
	// 	for _, fqdn := range srv.Fqdns {
	// 		for _, ip := range fqdn.Ips {
	// 			if !stringInSlice(fqdn.Name, fqdn_list) {
	// 				fmt.Println(fqdn.Service, ip.Service, fqdn.Name, ip.Ip) 
	// 				for _, p := range ip.AdminPorts {
	// 					fmt.Println(p.Num, p.Open)
	// 				}
	// 			}
	// 		}
	// 		fqdn_list = append(fqdn_list, fqdn.Name)
	// 	}
	// }
	// fmt.Println("===================")

	// fmt.Println("TURN connectivity:")
	// fqdn_list = make([]string, 0)
	// for _, srv := range discoveredsrv {
	// 	for _, fqdn := range srv.Fqdns {
	// 		for _, ip := range fqdn.Ips {
	// 			if !stringInSlice(fqdn.Name, fqdn_list) {
	// 				fmt.Println(fqdn.Service, fqdn.Name, ip.Ip) 
	// 				for _, p := range ip.TurnPorts {
	// 					fmt.Println(p.Num, p.Open, p.Proto)
	// 				}
	// 			}
	// 		}
	// 		fqdn_list = append(fqdn_list, fqdn.Name)
	// 	}
	// }


	fmt.Println("JJJJJJJJSSSOOOOOOONNNN")

	b, err := json.Marshal(DiscoveredSRVrecordsMap)
	if err != nil {
        fmt.Printf("Error: %s", err)
        return;
    }
    fmt.Println(string(b))


	


}
