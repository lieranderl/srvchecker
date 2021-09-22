package srv

import (
	// "log"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSrv(t *testing.T) {

	srvresults := new(SRVResults)
	srvresults.ForDomain("tp.ciscotac.net")

	fmt.Println(srvresults)


	if val, ok := (*srvresults)["_h323ls._udp.vodafone.com"]; ok {
		
		// assert.Equal(t, "udp", val.Proto)
		assert.Equal(t, "b2b", val.Sname)
		// assert.Equal(t, "1719", val.Port)

		fqdns := make([]Fqdn, 0) 
		ips := make([]string, 0)
		priority := make([]string,0)
		weight := make([]string,0)
		ports := make([]Port, 0)
		for k,v := range val.Fqdn {
			fqdns = append(fqdns, k)
			priority = append(priority, v.Priority)
			weight = append(weight, v.Weight)
			
			for ip, port := range v.Ips {
				ips = append(ips, string(ip))
				ports = append(ports, *port)
			}
		}
		assert.Equal(t, []Fqdn([]Fqdn{"vcs.vodafone.com.", "bc.vodafone.com."}), fqdns)
		assert.Equal(t, []string{"10", "1"}, priority)
		assert.Equal(t, []string{"50", "0"}, weight)
		assert.Equal(t, []string([]string{"A record not configured", "195.232.251.6"}), ips)
		assert.Equal(t, []string{}, ports)


	} else {
		t.Fail()
	}

	assert.Equal(t, "b2b", val.Sname)
		// assert.Equal(t, "1719", val.Port)

	fqdns := make([]Fqdn, 0) 
	ips := make([]string, 0)
	priority := make([]string,0)
	weight := make([]string,0)
	ports := make([]Port, 0)
	for k,v := range val.Fqdn {
		fqdns = append(fqdns, k)
		priority = append(priority, v.Priority)
		weight = append(weight, v.Weight)
		
		for ip, port := range v.Ips {
			ips = append(ips, string(ip))
			ports = append(ports, *port)
		}
	}
	fmt.Println(fqdns)
	fmt.Println(priority)
	fmt.Println(weight)
	fmt.Println(ips)
	fmt.Println(ports)



	// for k, res := range *srvresults {
	// 	log.Println("=================")
	// 	log.Println(k)
	// 	for _,r := range res {
	// 		log.Println(r.Fqdn, r.Ips, r.Port , r.Priority, r.Weight)

	// 	}
	// }
	// t.Fail()

}



