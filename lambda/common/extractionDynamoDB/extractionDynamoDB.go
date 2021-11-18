package extractionDynamoDB

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type EventLog struct {
	IPAddress string `json:"REQUEST-IP-ADDRESS"`
	DateTime  string `json:"DATETIME"`
}

func Login() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_KEY"),
				os.Getenv("SECRET_KEY"),
				"",
			),
			Region: aws.String(os.Getenv("REGION")),
		},
	}))

	return dynamodb.New(sess)
}

func GetAllIPAddressLogs(svc *dynamodb.DynamoDB, IPAddress string) ([]EventLog, error) {

	var eventLogs []EventLog
	startTime := time.Now().AddDate(0, 0, -1).Unix()
	endTime := time.Now().Unix()

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("ExtractionEngineRequestLogs"),
		KeyConditions: map[string]*dynamodb.Condition{
			"REQUEST-IP-ADDRESS": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(IPAddress),
					},
				},
			},
			"DATETIME": {
				ComparisonOperator: aws.String("BETWEEN"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(fmt.Sprintf("%v", startTime)),
					},
					{
						S: aws.String(fmt.Sprintf("%v", endTime)),
					},
				},
			},
		},
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		fmt.Printf("Error querying DynamoDB: %v\n", err)
		return eventLogs, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &eventLogs); err != nil {
		fmt.Printf("Error unmarshalling DynamoDB response: %v\n", err)
		return eventLogs, err
	}

	fmt.Printf("Found %v requests from the IP address %v in the past 24 hours\n", len(eventLogs), IPAddress)

	return eventLogs, nil
}

func CreateIPAddressLog(svc *dynamodb.DynamoDB, IPAddress string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"REQUEST-IP-ADDRESS": {S: aws.String(IPAddress)},
			"DATETIME":           {S: aws.String(fmt.Sprintf("%v", time.Now().Unix()))},
		},
		TableName: aws.String("ExtractionEngineRequestLogs"),
	}

	if _, err := svc.PutItem(input); err != nil {
		fmt.Printf("Error inserting event log into DynamoDB Table")
		return err
	}

	return nil
}
