package api

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// Response is a re-usable method for returning data from each request
func Response(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

