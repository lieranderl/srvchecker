package srvprocess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"strings"
	"time"

	"example.com/srvprocess/portconnectivity"
	"example.com/srvprocess/srv"
)


func Srvprocess(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")


	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}



	type inputRequest struct {
		Domain     string `form:"domain" json:"domain" binding:"required"`
		DnsServer  string `form:"dnsServer" json:"dnsServer"`
	}
	
	var json_input inputRequest


	if err := json.NewDecoder(r.Body).Decode(&json_input); err != nil {
		switch err {
		case io.EOF:
			fmt.Fprint(w, "Pizda!")
			return
		default:
			log.Printf("json.NewDecoder: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	if json_input.Domain == "" {
		fmt.Fprint(w, "no domain input")
		return
	}

	log.Println("Processing...")


	startTime := time.Now()
	srvresults := new(srv.DiscoveredSrvTable)
	srvresults.ForDomain(json_input.Domain)
	tcpConnectivityTable := make(portconnectivity.TcpConnectivityTable, 0)
	tcpConnectivityTable.FetchFromSrv(*srvresults)
	tcpConnectivityTable.Connectivity()

	elapsedTime := time.Since(startTime)
	stime := cut_time(elapsedTime)

	log.Println("All process took: ", stime)

	type H map[string]interface{}

	json.NewEncoder(w).Encode(H{ 
		"code" : http.StatusOK, 
		"elapsedTime": stime,
		"srv":  srvresults,
		"connectivity": tcpConnectivityTable, 
	})
	
}

func cut_time(elapsedTime time.Duration) string{
	s := elapsedTime.String()
	sslice:=strings.Split(s, ".") 
		s = s[:len(s)-6]
	if len(sslice) > 1 {
		sslice[1] = sslice[1][:(len(sslice[1])-(len(sslice[1])-2))]
		s = strings.Join(sslice, ".")
	}
	return s
}