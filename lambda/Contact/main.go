package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"textract-api/lambda/common/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/gomail.v2"
)

type requestBody struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

var (
	SenderName = "Extraction Engine"
	Sender     = os.Getenv("SenderEmailAddress")
	Recipient  = os.Getenv("Recipients")
	SmtpUser   = os.Getenv("SmtpUser")
	SmtpPass   = os.Getenv("SmtpPass")
	Host       = "email-smtp.eu-west-2.amazonaws.com"
	Port       = 587
)

// Entrypoint for the Lambda Function
func main() {
	lambda.Start(handleRequest)
}

// handleRequest takes the request body and performs the necessary commands
func handleRequest(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Decode the request body
	var data = requestBody{}
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return api.Response(http.StatusInternalServerError, err.Error())
	}

	// Send the message to Alex & I via email.
	if err := sendEmail(data); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return api.Response(http.StatusInternalServerError, "Failed to send user message.")
	}

	fmt.Printf("Successfully sent email to email to: %v\n", Recipient)
	return api.Response(http.StatusOK, "Successfully sent user message.")
}

// sendEmail takes the request body, and sends the data via email, for visibility.
func sendEmail(data requestBody) error {

	Subject := "New Message Received"

	HtmlBody := "<html>" +
		"<head><title>Extraction Engine Message Alert</title></head>" +
		"<body>" +
		"<h1>Extraction Engine Message Alert</h1><br/>&nbsp;" +
		fmt.Sprintf(`<p><b>From:  </b>%v</p>`, data.Name) +
		fmt.Sprintf(`<p><b>Email:  </b>%v</p>`, data.Email) +
		fmt.Sprintf(`<p><b>Message:  </b>%v</p>`, data.Message) +
		"</body>" +
		"</html>"

	m := gomail.NewMessage()
	m.SetBody("text/html", HtmlBody)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(Sender, SenderName)},
		"To":      {Recipient},
		"Subject": {Subject},
	})

	d := gomail.NewDialer(Host, Port, SmtpUser, SmtpPass)
	return d.DialAndSend(m)
}
