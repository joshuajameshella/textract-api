package textract

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"os"
)

// RequestBody is the data structure sent in the API request
type RequestBody struct {
	Data string `json:"data"`
	UserID string `json:"userID"`
}

// Login uses the AWS Secret Keys to create an AWS Textract client
func Login() (*textract.Textract, error) {

	// Create a new AWS credentials instance for use when making API calls
	awsCredentials := credentials.NewStaticCredentials(
		os.Getenv("ACCESS_KEY"),
		os.Getenv("SECRET_KEY"),
		"",
	)

	// Initiate a new AWS Textract session with the user info provided
	client := textract.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-2"), // London
		Credentials: awsCredentials,
	})))

	return client, nil
}

// Submit sends the encoded image data to Textract and returns the data from AWS.
func Submit(textractClient *textract.Textract, imageData []byte) ([]string, error) {

	resp, err := textractClient.DetectDocumentText(
		&textract.DetectDocumentTextInput{
			Document: &textract.Document{Bytes: imageData},
		},
	)
	if err != nil {
		return nil, err
	}

	// Process every word in the document, and store in a dictionary
	words := make([]string, 0)
	for _, block := range resp.Blocks {
		if *block.BlockType == "WORD" {
			words = append(words, *block.Text)
		}
	}

	return words, err
}

// CalculateFileSize calculates the file size based on the length of data provided.
func CalculateFileSize(imageData string) string {
	fileSizeBytes := float64(len(imageData)) * (float64(3) / float64(4))
	if fileSizeBytes > 1000 {
		return fmt.Sprintf("%v KB", fileSizeBytes / 1000)
	} else if fileSizeBytes > 1000000 {
		return fmt.Sprintf("%v MB", fileSizeBytes / 1000000)
	} else if fileSizeBytes > 1000000000 {
		return fmt.Sprintf("%v GB", fileSizeBytes / 1000000000)
	}
	return fmt.Sprintf("%f bytes", fileSizeBytes)
}