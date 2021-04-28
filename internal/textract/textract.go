package textract

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

// WordData is the AWS information returned for a specific word
type WordData struct {
	Word       string
	Confidence float64
}

// Textract Client
var textractSession *textract.Textract

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
