package extractionS3

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Session uses the AWS Secret Keys to create an AWS S3 session.
func S3Session() *session.Session {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_KEY"),
				os.Getenv("SECRET_KEY"),
				"",
			),
			Region: aws.String(os.Getenv("REGION")),
		},
	}))

	return awsSession
}

// UploadFile takes the file contents and file name, and uploads to S3.
func UploadFile(data []byte, filename string) error {
	uploader := s3manager.NewUploader(S3Session())
	reader := strings.NewReader(string(data))

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   reader,
	})

	return err
}
