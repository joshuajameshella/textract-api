package internal

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"textract-api/internal/textract"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

// ProcessImage processes the request body, and submits the data to AWS Textract
func ProcessImage(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Attempt to create an AWS Textract client for image processing
	if err := textract.Login(); err != nil {
		return APIResponse(http.StatusInternalServerError, ErrorBody{aws.String(err.Error())})
	}

	// Decode the request body
	var r = RequestBody{}
	if err := json.Unmarshal([]byte(req.Body), &r); err != nil {
		return APIResponse(http.StatusInternalServerError, ErrorBody{aws.String(err.Error())})
	}

	// If the data contains unnecessary formatting info, eg: "data:image/jpeg;base64,/9j/4AA...", remove it
	// and decode the Base64 string into a byte array
	imageData := strings.Split(r.Data, ",")
	decoded, err := base64.StdEncoding.DecodeString(imageData[len(imageData)-1])
	if err != nil {
		return APIResponse(http.StatusInternalServerError, ErrorBody{aws.String(err.Error())})
	}

	// Submit the image data to AWS Textract
	if wordData, err := textract.Submit(decoded); err != nil {
		return APIResponse(http.StatusInternalServerError, ErrorBody{aws.String(err.Error())})
	} else {
		return APIResponse(http.StatusOK, wordData)
	}
}

// UnhandledMethod is the default return type for unknown API request types
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return APIResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
