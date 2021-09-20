package main

import (
	"log"
	"srvchecker/output"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"time"
)




func main(){
	startTime := time.Now()
	srvresults := new(srv.SRVResults)
	srvresults.ForDomain("tp.ciscotac.net")
	portsresults := new(portconnectivity.PortsResults)
	portsresults.Connectivity(*srvresults)
	output.Output(srvresults, portsresults)
	elapsedTime := time.Since(startTime)
	log.Println("All process took: ", elapsedTime)
}
