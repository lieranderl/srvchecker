package webserver

import (
	"log"
	"net/http"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"time"

	"github.com/gin-gonic/gin"
)


func Srv_process(c *gin.Context) {

	type inputRequest struct {
		Domain     string `form:"domain" json:"domain" binding:"required"`
		DnsServer  string `form:"dnsServer" json:"dnsServer"`
	}
	
	var json inputRequest

	if c.BindJSON(&json) == nil {
		startTime := time.Now()
		srvresults := new(srv.DiscoveredSrvTable)
		srvresults.ForDomain(json.Domain)
		tcpConnectivityTable := make(portconnectivity.TcpConnectivityTable, 0)
		tcpConnectivityTable.FetchFromSrv(*srvresults)
		tcpConnectivityTable.Connectivity()

		elapsedTime := time.Since(startTime)
		stime := cut_time(elapsedTime)
	
		log.Println("All process took: ", stime)
		c.JSON(http.StatusOK, gin.H{ 
			"code" : http.StatusOK, 
			"elapsedTime": stime,
			"srv":  srvresults,
			"connectivity": tcpConnectivityTable, 
		})
		
	} else {
		c.JSON(405, gin.H{ 
			"code" : http.ErrBodyNotAllowed, 
			"elapsedTime": "", 
			"srv":  "",
			"connectivity": "", 
		})
	}
}