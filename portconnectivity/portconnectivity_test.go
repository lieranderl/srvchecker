package portconnectivity

import (
	"log"
	"srvchecker/srv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortconnectivity(t *testing.T) {

	srvresults := new(srv.SRVResults)
	srvresults.GetForDomain("verizon.com")
	var portsResults PortsResults
	portsResults.Connectivity(*srvresults)


	for _, res := range portsResults {
		if res.Fqdn == "ohtwbgcolec14p2.verizon.com." {
			assert.Equal(t, map[string]bool{"5061":true}, res.Ports)
			assert.Equal(t, "mssip", res.ServName)
			assert.Equal(t, "137.188.103.17", res.Ip)
		}
		log.Println(res.Fqdn, res.Ip, res.Ports)
	}
	

}



