package srv

import (
	// "log"
	"fmt"
	"testing"

)

func TestSrv(t *testing.T) {

	srvresults := new(DiscoveredSrvTable)
	srvresults.ForDomain("akbank.com")

	for _, res := range *srvresults {
		fmt.Println("res: ", res.Fqdn)
		for i, cert := range res.CertsChain{
			fmt.Println("cert ", i, cert)
		} 
		
	}
}
