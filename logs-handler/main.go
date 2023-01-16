package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var firetailApiUrl string
var firetailApiToken string

func main() {
	firetailApiUrl = "https://api.logging.eu-west-1.prod.firetail.app/logs/aws/appsync"
	firetailApiToken = os.Getenv("FIRETAIL_API_TOKEN")
	lambda.Start(Handler)
}
