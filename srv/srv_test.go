package srv

import (
	// "log"
	"fmt"
	"testing"

)

func TestSrv(t *testing.T) {

	srvresults := new(DiscoveredSrvTable)
	srvresults.ForDomain("strayer.edu")

	fmt.Println(srvresults)

	for _, res := range *srvresults {
		fmt.Println("=================")
		fmt.Println(res)
	}
	t.Fail()

}
