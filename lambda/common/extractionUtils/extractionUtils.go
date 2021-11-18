package extractionUtils

import (
	"encoding/json"
	"log"
	"textract-api/lambda/common/extractionDynamoDB"
)

type LogInfo struct {
	IPAddress string
	DateTime  string
}

// JSONLog pretty prints the data to AWS CloudWatch
func JSONLog(header string, object interface{}) {
	if objectJSON, err := json.Marshal(object); err != nil {
		log.Printf("JSONLog : Error trying to marshal and log: %s : %+v\n", header, err)
		return
	} else {
		log.Println(header, string(objectJSON))
	}
}

// CreateIPEvent creates an event log for a given IP address in DynamoDB.
func CreateIPEvent(IPAddress string) error {
	svc := extractionDynamoDB.Login()
	return extractionDynamoDB.CreateIPAddressLog(svc, IPAddress)
}

// CountIPEvents reads all DynamoDB events within the last 24 hours for a given IP address.
func CountIPEvents(IPAddress string) (int, error) {
	svc := extractionDynamoDB.Login()
	logs, err := extractionDynamoDB.GetAllIPAddressLogs(svc, IPAddress)
	return len(logs), err
}
