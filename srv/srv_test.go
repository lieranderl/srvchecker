package srv

import (
	// "log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSrv(t *testing.T) {

	srvresults := new(SRVResults)
	srvresults.GetForDomain("verizon.com")

	if val, ok := (*srvresults)["_collab-edge._tls.verizon.com"]; ok {
		for _, v:= range val {
			if v.Fqdn == "ohtwbgcolec12p2.verizon.com." {
				assert.Equal(t, "8443", v.Port)
				assert.Equal(t, "1", v.Priority)
				assert.Equal(t, "10", v.Weight)
				assert.Equal(t, []string{"137.188.103.13"}, v.Ips)
			}
		}
	} else {
		t.Fail()
	}

	if val, ok := (*srvresults)["_sips._tcp.verizon.com"]; ok {
		assert.Equal(t, "SRV record not configured.", val[0].Fqdn)
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



