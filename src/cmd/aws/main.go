package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lieranderl/srvchecker/portconnectivity"
	"github.com/lieranderl/srvchecker/srv"
)

// inputRequest represents the structure of the incoming JSON payload
type inputRequest struct {
	Domain    string `json:"domain" binding:"required"`
	DnsServer string `json:"dnsServer"`
}

// responseData defines the structure for the JSON response
type responseData struct {
	Code         int                                   `json:"code"`
	ElapsedTime  string                                `json:"elapsedTime"`
	Srv          *srv.DiscoveredSrvTable               `json:"srv"`
	Connectivity portconnectivity.TcpConnectivityTable `json:"connectivity"`
}

// processDomain performs SRV discovery and connectivity checks for a domain
func processDomain(domain string) (*srv.DiscoveredSrvTable, portconnectivity.TcpConnectivityTable) {
	// Perform SRV discovery
	srvResults := new(srv.DiscoveredSrvTable)
	srvResults.ForDomain(domain)

	// Perform connectivity checks
	connectivityTable := make(portconnectivity.TcpConnectivityTable, 0)
	connectivityTable.FetchFromSrv(*srvResults).Connectivity()

	return srvResults, connectivityTable
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON payload
	var input inputRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Printf("Failed to parse JSON payload: %v\n", err)
		return
	}
	defer r.Body.Close()

	// Validate the input
	if input.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		log.Println("Domain is required but missing")
		return
	}

	// Start processing
	log.Println("Processing SRV and connectivity data for domain:", input.Domain)
	startTime := time.Now()

	// Perform SRV discovery and connectivity checks
	srvResults, connectivityResults := processDomain(input.Domain)

	// Calculate elapsed time
	elapsedTime := time.Since(startTime).Round(time.Millisecond).String()
	log.Printf("Processing completed in %s\n", elapsedTime)

	// Prepare the response
	response := responseData{
		Code:         200,
		ElapsedTime:  elapsedTime,
		Srv:          srvResults,
		Connectivity: connectivityResults,
	}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to process response", http.StatusInternalServerError)
		log.Printf("Failed to encode JSON response: %v\n", err)
		return
	}
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Register the POST / route
	mux.HandleFunc("POST /", mainHandler)

	// get port from environment variable or use default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server is listening on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
