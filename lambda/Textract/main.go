package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsTextract "github.com/aws/aws-sdk-go/service/textract"
	"textract-api/lambda/common/api"
	"textract-api/lambda/common/extractionUtils"
	"textract-api/lambda/common/textract"
)

// Handler manages requests to AWS Textract
type Handler struct {
	client *awsTextract.Textract
}

// Entrypoint for the Lambda Function
func main() {
	handler := Handler{}
	lambda.Start(handler.handleRequest)
}

// handleRequest takes the request body and performs the necessary commands
func (handler *Handler) handleRequest(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Create a new AWS Textract client used for requests
	var err error
	if handler.client, err = textract.Login(); err != nil {
		extractionUtils.JSONLog("Error creating Textract Client: ", err.Error())
		return api.Response(http.StatusInternalServerError, err)
	}

	// Decode the request body
	var r = textract.RequestBody{}
	if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	start := time.Now()
	extractionUtils.JSONLog("New request from user: ", r.UserID)
	extractionUtils.JSONLog("File Size: ", textract.CalculateFileSize(r.Data))

	// Decode the data from string to []byte, required by AWS.
	decoded, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	// Submit the image data to AWS Textract
	wordData, err := textract.Submit(handler.client, decoded)
	if err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	log.Printf("Request took %s", time.Since(start))
	return api.Response(http.StatusOK, wordData)
}
