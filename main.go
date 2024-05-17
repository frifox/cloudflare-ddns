package main

import (
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var AppUser = os.Getenv("APP_USER")
var AppPass = os.Getenv("APP_PASS")
var AuthHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(AppUser+":"+AppPass))

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader, ok := request.Headers["authorization"]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "no auth", StatusCode: 400}, nil
	}
	if authHeader != AuthHeader {
		return events.APIGatewayProxyResponse{Body: "bad auth", StatusCode: 400}, nil
	}

	hostname := request.QueryStringParameters["hostname"]
	if hostname == "" {
		return events.APIGatewayProxyResponse{Body: "hostname empty", StatusCode: 400}, nil
	}
	ip := request.QueryStringParameters["ip"]
	if ip == "" {
		return events.APIGatewayProxyResponse{Body: "ip empty", StatusCode: 400}, nil
	}

	return cloudflare(hostname, ip)
}
