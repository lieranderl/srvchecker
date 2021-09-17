package portconnectivity

import (
	// "log"
	"log"
	"srvchecker/srv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortconnectivity(t *testing.T) {

	srvresults := new(srv.SRVResults)
	srvresults.GetForDomain("verizon.com")
	var portsResults PortsResults
	portsResults.Connectivity(*srvresults, )


	for _, res := range portsResults {
		log.Println(res.Fqdn, res.Ip, res.Ports)
	}
	assert.Equal(t, "8443", portsResults[0].Ports)
	t.Fail()

}



