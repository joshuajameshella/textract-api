package extractionUtils

import (
	"encoding/json"
	"log"
)

// JSONLog pretty prints the data to AWS CloudWatch
func JSONLog(header string, object interface{}) {
	if objectJSON, err := json.Marshal(object); err != nil {
		log.Printf("JSONLog : Error trying to marshal and log: %s : %+v\n", header, err)
		return
	} else {
		log.Println(header, string(objectJSON))
	}
}
