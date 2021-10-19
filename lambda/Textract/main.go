package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"textract-api/lambda/common/api"
	"textract-api/lambda/common/extractionS3"
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

	// Submit the file for AWS processing.
	jobID, err := textract.StartTextractProcess(handler.client, r.Data)
	if err != nil {
		extractionUtils.JSONLog("Error processing file : ", fmt.Sprintf("%v", err))
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	extractionUtils.JSONLog("Job started. ID: ", *jobID)

	// Periodically check the status of the textract process running in AWS, using the jobID
	jobComplete, err := textract.IsJobComplete(handler.client, jobID)
	if err != nil {
		extractionUtils.JSONLog("Error checking job status : ", fmt.Sprintf("%v", err))
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	// When the textract process is complete, get the job results, marshall into []byte data, and upload to S3.
	if jobComplete {
		data, err := textract.GetJobResults(handler.client, jobID)
		if err != nil {
			extractionUtils.JSONLog("Error retrieving job data : ", fmt.Sprintf("%v", err))
			return api.Response(http.StatusInternalServerError, err.Error())
		}

		processedData, err := json.Marshal(data)
		if err != nil {
			extractionUtils.JSONLog("Error converting struct to json : ", fmt.Sprintf("%v", err))
			return api.Response(http.StatusInternalServerError, err.Error())
		}

		processedFileName := fmt.Sprintf("%v.json", strings.Split(r.Data, ".")[0])

		if err = extractionS3.UploadToS3(processedData, processedFileName); err != nil {
			extractionUtils.JSONLog("Error uploading json data : ", fmt.Sprintf("%v", err))
			return api.Response(http.StatusInternalServerError, err.Error())
		}

		extractionUtils.JSONLog("File successfully processed! ", processedFileName)
	}

	return api.Response(http.StatusOK, "Success")
}
