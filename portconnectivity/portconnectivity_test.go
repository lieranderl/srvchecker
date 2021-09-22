package portconnectivity

import (
	"fmt"
	"srvchecker/srv"
	"testing"
)

func TestPortconnectivity(t *testing.T) {
	srvresults := new(srv.SRVResults)
	srvresults.ForDomain("vodafone.com")

	var portsResults PortsResults
	portsResults.fetchFromSrvResults(srvresults)

	fmt.Println(srvresults)
	fmt.Println(portsResults)
	
	t.Fail()

}

// func TestTurn(t *testing.T) {
// 	PortsResult := new(PortsResult)
// 	PortsResult.RunTurnCheck("173.38.154.85", "3478", true)
// 	t.Fail()
// }
