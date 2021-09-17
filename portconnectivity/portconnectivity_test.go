package portconnectivity

import (
	// "log"
	"srvchecker/srv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortconnectivity(t *testing.T) {

	srvresults := new(srv.SRVResults)
	srvresults.GetForDomain("verizon.com")
	var portsResults PortsResults
	portsResults.Connectivity(*srvresults)


	// for k, res := range *srvresults {
	// 	log.Println("=================")
	// 	log.Println(k)
	// 	for _,r := range res {
	// 		log.Println(r.Fqdn, r.Ips, r.Port , r.Priority, r.Weight)

	// 	}
	// }
	assert.Equal(t, "8443", portsResults[0].ports)
	t.Fail()

}



