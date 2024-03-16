package tests

import (
	"testing"

	"github.com/lieranderl/srvchecker/srv"
)

func TestCerts(t *testing.T) {

	dr := srv.DiscoveredSrvRow{}
	dr.Connect_cert("cisco.com", "443")

	if dr.CertsChain == nil {
		t.Error("CertsChain is nil")
	}

	if dr.CertsChain[0].Cn == "www.cisco.com" {
		t.Log("CommonName is www.cisco.com")
	} else {
		t.Error("CommonName is not www.cisco.com")
	}

}
