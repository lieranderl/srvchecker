package srv

import (
	// "log"
	"fmt"
	"testing"

)

func TestSrv(t *testing.T) {

	srvresults := new(DiscoveredSrvTable)
	srvresults.ForDomain("mofa.gov.sa")


	for _, res := range *srvresults {
		
		if res.Certs != nil {
			fmt.Println("=================")
			fmt.Println(res.Certs)
			if res.Certs[0].Child != nil {
				fmt.Println("=================")
				fmt.Println(res.Certs[0].Child[0])
				if res.Certs[0].Child[0].Child != nil {
					fmt.Println("=================")
					fmt.Println(res.Certs[0].Child[0].Child[0])
				}
			}
			
		}
		
		
	}
	t.Fail()

}
