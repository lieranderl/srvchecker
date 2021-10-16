package srv

import (
	// "log"
	"fmt"
	"testing"

)

func TestSrv(t *testing.T) {

	srvresults := new(DiscoveredSrvTable)
	srvresults.ForDomain("mofa.gov.sa")

	fmt.Println(srvresults)

	for _, res := range *srvresults {
		fmt.Println("=================")
		for _, cert := range res.Certs {
			fmt.Println("=========CERT========")
			fmt.Println(cert)
			fmt.Println("=========END========")
		}
	}
	t.Fail()

}
