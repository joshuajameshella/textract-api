package textract

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

// RequestBody is the data structure sent in the API request.
type RequestBody struct {
	Data   string `json:"data"`
	UserID string `json:"userID"`
}

// WordData is the data associated with each word in a document.
type WordData struct {
	ID         string
	Page       int64
	Word       string
	Confidence float64
}

// Login uses the AWS Secret Keys to create an AWS Textract client.
func Login() *textract.Textract {
	awsCredentials := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("ACCESS_KEY"),
			os.Getenv("SECRET_KEY"),
			"",
		),
		Region: aws.String(os.Getenv("REGION")),
	}
	awsSession := session.Must(session.NewSession())
	return textract.New(awsSession, awsCredentials)
}

// StartTextractProcess takes a document in S3, and begins processing it - Returns a JobID used to check progress.
func StartTextractProcess(client *textract.Textract, documentName string) (*string, error) {
	response, err := client.StartDocumentTextDetection(&textract.StartDocumentTextDetectionInput{
		DocumentLocation: &textract.DocumentLocation{
			S3Object: &textract.S3Object{
				Bucket:  aws.String(os.Getenv("S3_BUCKET_NAME")),
				Name:    aws.String(documentName),
				Version: nil,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return response.JobId, err
}

// CheckJobStatus takes a JobID, and queries AWS for the status of the extraction process.
func CheckJobStatus(client *textract.Textract, jobID *string) (string, error) {
	result, err := client.GetDocumentTextDetection(&textract.GetDocumentTextDetectionInput{
		JobId: jobID,
	})
	return *result.JobStatus, err
}

// IsJobComplete periodically checks the job status of a given JobID.
func IsJobComplete(client *textract.Textract, jobID *string) (bool, error) {
	jobCompleted := false
	status, err := CheckJobStatus(client, jobID)
	if err != nil {
		return false, err
	}

	for !jobCompleted {
		time.Sleep(time.Millisecond * 500)
		status, _ = CheckJobStatus(client, jobID)
		jobCompleted = status == "SUCCEEDED"
	}

	return jobCompleted, nil
}

// GetJobResults retrieves the data of a completed Textract process.
func GetJobResults(client *textract.Textract, jobID *string) ([][]WordData, error) {
	result, err := client.GetDocumentTextDetection(&textract.GetDocumentTextDetectionInput{
		JobId: jobID,
	})

	pageCount := 1
	data := make([]WordData, 0)
	for _, block := range result.Blocks {
		if *block.BlockType == "WORD" {
			data = append(data, WordData{
				ID:         *block.Id,
				Page:       *block.Page,
				Word:       *block.Text,
				Confidence: *block.Confidence,
			})
		}
		if int(*block.Page) > pageCount {
			pageCount = int(*block.Page)
		}
	}

	documentStructure := make([][]WordData, pageCount)
	for _, words := range data {
		index := int(words.Page) - 1
		tempPageData := documentStructure[index]
		tempPageData = append(tempPageData, words)
		documentStructure[index] = tempPageData
	}

	return documentStructure, err
}
