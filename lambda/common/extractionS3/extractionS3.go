package extractionS3

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Client uses the AWS Secret Keys to create an AWS S3 client.
func S3Client() *s3manager.Uploader {
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
	return s3manager.NewUploader(awsSession)
}

// UploadToS3 takes the file contents and file name, and uploads to S3.
func UploadToS3(data []byte, filename string) error {

	reader := strings.NewReader(string(data))
	uploader := S3Client()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   reader,
	})

	return err
}
