package lambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var ErrorMethodNotAllowed = "method Not allowed"

// ErrorBody is a re-usable response structure sent when an error occurs
type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// RequestBody is the data structure sent in a POST request
type RequestBody struct {
	Data string `json:"data"`
}

// APIResponse is a re-usable method for returning data from each request
func APIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
