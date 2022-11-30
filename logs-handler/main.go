package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var firetailApiUrl string
var firetailApiKey string

func main() {
	firetailApiUrl = "https://api.logging.eu-west-1.sandbox.firetail.app/logs/aws/appsync"
	firetailApiKey = os.Getenv("FIRETAIL_API_KEY")
	lambda.Start(Handler)
}
