package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"net/http"
	"os"
)

var CFEmail = os.Getenv("CF_EMAIL")
var CFKey = os.Getenv("CF_KEY")
var CFZoneID = os.Getenv("CF_ZONE_ID")
var CFRecordID = os.Getenv("CF_RECORD_ID")

type CFRequest struct {
	Type    string
	Content string
	Name    string
	Proxied bool
	TTL     int
}

func cloudflare(hostname string, ip string) (events.APIGatewayProxyResponse, error) {
	// prepare request
	cfReq, _ := json.Marshal(CFRequest{
		Type:    "A",
		Name:    hostname,
		Content: ip,
		Proxied: false,
		TTL:     3600,
	})
	uri := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", CFZoneID, CFRecordID)
	req, _ := http.NewRequest("PUT", uri, bytes.NewReader(cfReq))
	req.Header.Set("X-Auth-Email", CFEmail)
	req.Header.Set("X-Auth-Key", CFKey)

	// call cloudflare
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	// read response
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}
	err = response.Body.Close()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	// fwd response to client
	return events.APIGatewayProxyResponse{Body: string(respBody), StatusCode: response.StatusCode}, nil
}
