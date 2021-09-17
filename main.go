package main

import (
	"log"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"time"
)



func main(){
	startTime := time.Now()

	
	srvresults := new(srv.SRVResults)
	srvresults.GetForDomain("verizon.com")
	portsresults := new(portconnectivity.PortsResults)
	portsresults.Connectivity(*srvresults)



	for k, res := range *srvresults {
		log.Println("=================")
		log.Println(k)
		for _,r := range res {
			log.Println(r.ServName, r.Fqdn, r.Ips, r.Port, r.Priority, r.Weight)
		}
	}
	for _, res := range *portsresults {
		log.Println(res.Fqdn, res.Ip, res.Ports)	
	}

	elapsedTime := time.Since(startTime)
	log.Println("All process took: ", elapsedTime)
}