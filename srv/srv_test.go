package srv

import (
	// "log"
	"fmt"
	"testing"

)

func TestSrv(t *testing.T) {

	srvresults := new(DiscoveredSrvTable)
	srvresults.ForDomain("akbank.com")

	for _, res := range *srvresults {
		
		if res.Children != nil {
			fmt.Println("=================")
			fmt.Println(res.Children)
			if res.Children[0].Children != nil {
				fmt.Println("=================")
				fmt.Println(res.Children[0].Children[0])
				if res.Children[0].Children[0].Children != nil {
					fmt.Println("=================")
					fmt.Println(res.Children[0].Children[0].Children[0])
				}
			}
			
		}
		
		
	}
}
