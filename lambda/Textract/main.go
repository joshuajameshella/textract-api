package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"textract-api/lambda/common/api"
	"textract-api/lambda/common/extractionUtils"
	"textract-api/lambda/common/textract"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsTextract "github.com/aws/aws-sdk-go/service/textract"
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
	handler.client = textract.Login()

	// Decode the request body
	var r = textract.RequestBody{}
	if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	extractionUtils.JSONLog("New request from user: ", r.UserID)

	for _, fileName := range r.Data {

		jobID, err := textract.StartTextractProcess(handler.client, fileName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return api.Response(http.StatusInternalServerError, err.Error())
		}

		extractionUtils.JSONLog("Job started. ID: ", *jobID)

		if jobComplete, err := textract.IsJobComplete(handler.client, jobID); err != nil {
			fmt.Printf("Error: %v\n", err)
			return api.Response(http.StatusInternalServerError, err.Error())

		} else if jobComplete {

			if _, err = textract.GetJobResults(handler.client, jobID); err != nil {
				fmt.Printf("Error: %v\n", err)
				return api.Response(http.StatusInternalServerError, err.Error())
			} else {
				// TODO: Upload data to S3, in json format
			}
		}
	}

	return api.Response(http.StatusOK, "Success")
}
