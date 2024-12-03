package srvprocess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// Srvprocess handles the SRV processing request
func Srvprocess(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse and validate the input request
	var input inputRequest
	if err := parseRequestBody(r.Body, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	log.Println("Processing SRV and connectivity data for domain:", input.Domain)
	startTime := time.Now()

	// Perform SRV discovery and connectivity checks
	srvResults, connectivityResults := processDomain(input.Domain)

	// Calculate elapsed time
	elapsedTime := time.Since(startTime).Round(time.Millisecond).String()
	log.Printf("Processing completed in %s\n", elapsedTime)

	// Respond with the results
	respondWithJSON(w, http.StatusOK, responseData{
		Code:         http.StatusOK,
		ElapsedTime:  elapsedTime,
		Srv:          srvResults,
		Connectivity: connectivityResults,
	})
}

// setCORSHeaders sets the CORS headers for the response
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
}

// parseRequestBody parses and validates the JSON body of the incoming request
func parseRequestBody(body io.ReadCloser, input *inputRequest) error {
	defer body.Close()
	if err := json.NewDecoder(body).Decode(input); err != nil {
		if err == io.EOF {
			return fmt.Errorf("request body is empty")
		}
		log.Printf("Failed to parse request body: %v\n", err)
		return fmt.Errorf("invalid request body")
	}
	return nil
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

// respondWithJSON encodes a response as JSON and writes it to the response writer
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write JSON response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
