package tests

import (
	"github.com/lieranderl/srvchecker/srv"
	"testing"
)

func TestSrv(t *testing.T) {

	srvresults := new(srv.DiscoveredSrvTable)
	srvresults.ForDomain("akbank.com")

	for _, res := range *srvresults {
		t.Log("res: ", res.Fqdn)
		for i, cert := range res.CertsChain {
			t.Log("cert ", i, cert)
		}

	}
}
