package tests

import (

	"testing"

	"github.com/lieranderl/srvchecker/srv"
	"github.com/lieranderl/srvchecker/portconnectivity"
)

func TestPortconnectivity(t *testing.T) {
	
	srvresults := new(srv.DiscoveredSrvTable)
	srvresults.ForDomain("cisco.com")
	
	for _, res := range *srvresults {
		t.Log("=================")
		t.Log(res)
	}
	
	tcpConnectivityTable := make(portconnectivity.TcpConnectivityTable, 0)
	tcpConnectivityTable.FetchFromSrv(*srvresults).Connectivity()


	for _, row := range tcpConnectivityTable {
		t.Log("=========TCP========")
		t.Log(row.Fqdn)
		t.Log(row.Ip)
		for _, port := range row.Ports {
			t.Log(port)
		}
	}
	

}

// func TestTurn(t *testing.T) {
// 	PortsResult := new(PortsResult)
// 	PortsResult.RunTurnCheck("173.38.154.85", "3478", true)
// 	t.Fail()
// }


// func TestSyncMap(t *testing.T) {
// 	// var ss testPizda
// 	mapsrv := SyncMap()

// 	m := map[string]sync.Map{}
// 	mapsrv.Range(func(key, value interface{}) bool {
// 		m[fmt.Sprint(key)] = value.(sync.Map)
// 		return true
// 	})

// 	b, err := json.MarshalIndent(mapsrv, "", " ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(b))


// 	t.Fail()
// }