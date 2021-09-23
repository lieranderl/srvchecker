package portconnectivity

import (
	"encoding/json"
	"fmt"
	"srvchecker/srv"


	// "sync"
	"testing"
)

func TestPortconnectivity(t *testing.T) {
	srvresults := new(srv.SRVResults)
	srvresults.ForDomain("verizon.com")

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