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
	srvresults.GetForDomain("ciscotac.net")
	portsresults := new(portconnectivity.PortsResults)
	portsresults.Connectivity(*srvresults)



	// for k, res := range *srvresults {
	// 	log.Println("=================")
	// 	log.Println(k)
	// 	for _,r := range res {
	// 		log.Println(r.ServName, r.Fqdn, r.Ips, r.Port, r.Priority, r.Weight)
	// 	}
	// }
	// for _, res := range *portsresults {
	// 	log.Println(res.Fqdn, res.Ip, res.ServName, res.Port)	
	// 	 for _, cert := range res.Certs {
	// 		log.Printf("Issuer Name: %s\n", cert.Issuer)
	// 		log.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
	// 		log.Printf("Common Name: %s \n", cert.Issuer.CommonName)
	// 	}
	// }

	
	output.Output(srvresults, portsresults)



	// b, err := json.Marshal(discoveredsrv)
	// if err != nil {
    //     fmt.Printf("Error: %s", err)
    //     return;
    // }
    // fmt.Println(string(b))



	elapsedTime := time.Since(startTime)
	log.Println("All process took: ", elapsedTime)
}
