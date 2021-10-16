package portconnectivity

import (

	"fmt"
	"srvchecker/srv"


	// "sync"
	"testing"
)

func TestPortconnectivity(t *testing.T) {
	
	srvresults := new(srv.DiscoveredSrvTable)
	srvresults.ForDomain("cisco.com")
	
	for _, res := range *srvresults {
		fmt.Println("=================")
		fmt.Println(res)
	}
	
	tcpConnectivityTable := make(TcpConnectivityTable, 0)
	tcpConnectivityTable.FetchFromSrv(*srvresults)
	tcpConnectivityTable.Connectivity()


	for _, row := range tcpConnectivityTable {
		fmt.Println("=========TCP========")
		fmt.Println(row.Fqdn)
		fmt.Println(row.Ip)
		fmt.Println(row.ServiceName)
		for _, port := range row.Ports {
			fmt.Print(port)
		}
	}
	
	
	t.Fail()

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