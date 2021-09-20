package portconnectivity

import (
	"log"
	"srvchecker/srv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortconnectivity(t *testing.T) {

	srvresults := new(srv.SRVResults)
	srvresults.ForDomain("verizon.com")
	var portsResults PortsResults
	portsResults.Connectivity(*srvresults)

	for _, res := range portsResults {
		if res.Fqdn == "ohtwbgcolec14p2.verizon.com." {
			assert.Equal(t, map[string]bool{"5061": true}, res.Port)
			assert.Equal(t, "mssip", res.ServName)
			assert.Equal(t, "137.188.103.17", res.Ip)
		}
		log.Println(res.Fqdn, res.Ip, res.Port)
	}

}

// func TestTurn(t *testing.T) {
// 	PortsResult := new(PortsResult)
// 	PortsResult.RunTurnCheck("173.38.154.85", "3478", true)
// 	t.Fail()
// }
