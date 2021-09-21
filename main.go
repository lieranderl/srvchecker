package main

import (
	"fmt"
	"log"
	"net/http"
	"srvchecker/output"
	"srvchecker/portconnectivity"
	"srvchecker/srv"
	"time"

	"github.com/gin-gonic/gin"
)




func main(){

	router := gin.Default()

	
	router.POST("/srv_process", srv_process)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()

}

func srv_process(c *gin.Context) {

	domain := c.PostForm("domain")
	
	startTime := time.Now()
	srvresults := new(srv.SRVResults)
	srvresults.ForDomain(domain)
	portsresults := new(portconnectivity.PortsResults)
	portsresults.Connectivity(*srvresults)
	srv, nosrv, connectivity := output.Output(srvresults, portsresults)
	elapsedTime := time.Since(startTime)
	log.Println("All process took: ", elapsedTime)


	c.JSON(http.StatusOK, gin.H{ 
		"code" : http.StatusOK, 
		"elapsedTime": fmt.Sprint(elapsedTime), 
		"srv":  srv,
		"nosrv": nosrv,
		"connectivity": connectivity, 
	})
	
	
}