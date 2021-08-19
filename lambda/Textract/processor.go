package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"textract-api/lambda"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

// Textract Client
var textractSession *textract.Textract

// WordData is the AWS information returned for a specific word
type WordData struct {
	Word       string
	Confidence float64
}

// Login uses the AWS Secret Keys to create an AWS Textract client
func Login() error {

	// Create a new AWS credentials instance for use when making API calls
	awsCredentials := credentials.NewStaticCredentials(
		os.Getenv("ACCESS_KEY"),
		os.Getenv("SECRET_KEY"),
		"",
	)

	// Initiate a new AWS Textract session with the user info provided
	textractSession = textract.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-2"), // London
		Credentials: awsCredentials,
	})))

	// TODO: Perform user credential check here. Return error if user rejected
	return nil
}

// Submit sends the encoded image data to Textract and returns the data from AWS.
func Submit(imageData []byte) ([][]WordData, error) {

	resp, err := textractSession.DetectDocumentText(
		&textract.DetectDocumentTextInput{
			Document: &textract.Document{Bytes: imageData},
		},
	)
	if err != nil {
		return nil, err
	}

	// Process every word in the document, and store in a dictionary
	wordDictionary := make(map[string]WordData)
	for _, block := range resp.Blocks {
		if *block.BlockType == "WORD" {
			wordDictionary[*block.Id] = WordData{
				Word:       *block.Text,
				Confidence: *block.Confidence,
			}
		}
	}

	// Process each line in the document, and return an array of document lines
	documentStructure := make([][]WordData, 0)
	for _, block := range resp.Blocks {
		if *block.BlockType == "LINE" {

			// Get a list of ID's on the line, and reconstruct the sentence using the dictionary
			lineStructure := make([]WordData, 0)
			for _, id := range block.Relationships[0].Ids {
				if word, exists := wordDictionary[*id]; exists {
					lineStructure = append(lineStructure, word)
				}
			}
			documentStructure = append(documentStructure, lineStructure)
		}
	}

	return documentStructure, err
}

// ProcessImage processes the request body, and submits the data to AWS Textract
func ProcessImage(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Attempt to create an AWS Textract client for image processing
	if err := Login(); err != nil {
		return lambda.APIResponse(http.StatusInternalServerError, lambda.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	// Decode the request body
	var r = lambda.RequestBody{}
	if err := json.Unmarshal([]byte(req.Body), &r); err != nil {
		return lambda.APIResponse(http.StatusInternalServerError, lambda.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	// If the data contains unnecessary formatting info, eg: "data:image/jpeg;base64,/9j/4AA...", remove it
	// and decode the Base64 string into a byte array
	imageData := strings.Split(r.Data, ",")
	decoded, err := base64.StdEncoding.DecodeString(imageData[len(imageData)-1])
	if err != nil {
		return lambda.APIResponse(http.StatusInternalServerError, lambda.ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	// Submit the image data to AWS Textract
	if wordData, err := Submit(decoded); err != nil {
		return lambda.APIResponse(http.StatusInternalServerError, lambda.ErrorBody{ErrorMsg: aws.String(err.Error())})
	} else {
		return lambda.APIResponse(http.StatusOK, wordData)
	}
}

// UnhandledMethod is the default return type for unknown API request types
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return lambda.APIResponse(http.StatusMethodNotAllowed, lambda.ErrorMethodNotAllowed)
}
