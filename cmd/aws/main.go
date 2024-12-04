package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log the incoming event
	log.Printf("Received event: %+v\n", event)

	// Parse the JSON payload
	var input inputRequest
	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		log.Printf("Failed to parse JSON payload: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body:       "Invalid JSON payload",
		}, nil
	}

	// Validate the input
	if input.Domain == "" {
		log.Println("Domain is required but missing")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Domain is required",
		}, nil
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
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to process response",
		}, nil
	}

	// Return the API Gateway response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	// Start the Lambda function
	lambda.Start(handler)
}
