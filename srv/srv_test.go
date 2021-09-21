package srv

import (
	// "log"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSrv(t *testing.T) {

	srvresults := new(SRVResults)
	srvresults.ForDomain("vodafone.com")

	fmt.Println(srvresults)


	if val, ok := (*srvresults)["_h323ls._udp.vodafone.com"]; ok {
		
		assert.Equal(t, "udp", val.Proto)
		assert.Equal(t, "b2b", val.Sname)
		assert.Equal(t, "1719", val.Port)

		fqdns := make([]string, 0) 
		ips := make([][]string, 0)
		priority := make([]string,0)
		weight := make([]string,0)
		for k,v := range val.Fqdn {
			fqdns = append(fqdns, k)
			ips = append(ips, v.Ips)
			priority = append(priority, v.Priority)
			weight = append(weight, v.Weight)
		}
		assert.Equal(t, []string([]string{"vcs.vodafone.com.", "bc.vodafone.com."}), fqdns)
		assert.Equal(t, []string{"10", "1"}, priority)
		assert.Equal(t, []string{"50", "0"}, weight)
		assert.Equal(t, [][]string([][]string{[]string{"A record not configured"}, []string{"195.232.251.6"}}), ips)


	} else {
		t.Fail()
	}


	// for k, res := range *srvresults {
	// 	log.Println("=================")
	// 	log.Println(k)
	// 	for _,r := range res {
	// 		log.Println(r.Fqdn, r.Ips, r.Port , r.Priority, r.Weight)

	// 	}
	// }
	// t.Fail()

}



