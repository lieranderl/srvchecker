package main

import (
	"log"
	"srvchecker/srv"
)



func main(){
	
	srvresults := new(srv.SRVResults)
	srvresults.GetForDomain("verizon.com")


	for k, res := range *srvresults {
		log.Println("=================")
		log.Println(k)
		for _,r := range res {
			log.Println(r.Fqdn, r.Ips, r.Port , r.Priority, r.Weight)
		}
	}
}