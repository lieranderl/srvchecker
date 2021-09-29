package main

import (
	"fmt"
	"log"
	"net/http"

	// "srvchecker/portconnectivity"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"time"

	"github.com/gin-gonic/gin"
)


func main(){
	router := gin.Default()
	router.Use(CORSMiddleware())
	
	router.POST("/srv_process", srv_process)
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()

}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func srv_process(c *gin.Context) {

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
		log.Println("All process took: ", elapsedTime)
		c.JSON(http.StatusOK, gin.H{ 
			"code" : http.StatusOK, 
			"elapsedTime": fmt.Sprint(elapsedTime), 
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