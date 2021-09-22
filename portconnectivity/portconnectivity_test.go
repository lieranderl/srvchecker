package portconnectivity

import (
	"encoding/json"
	"fmt"
	"srvchecker/srv"
	"testing"
)

func TestPortconnectivity(t *testing.T) {
	srvresults := new(srv.SRVResults)
	srvresults.ForDomain("globits.de")

	var portsResults PortsResults
	portsResults.FetchFromSrvResults(srvresults)


	nosrv, err := json.Marshal(srvresults)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(string(nosrv))

	nosrv, err = json.Marshal(portsResults)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(string(nosrv))
	
	t.Fail()

}

// func TestTurn(t *testing.T) {
// 	PortsResult := new(PortsResult)
// 	PortsResult.RunTurnCheck("173.38.154.85", "3478", true)
// 	t.Fail()
// }
