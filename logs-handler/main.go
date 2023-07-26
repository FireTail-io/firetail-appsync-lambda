package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

const DefaultFiretailApiUrl string = "https://api.logging.eu-west-1.prod.firetail.app/logs/aws/appsync"

var firetailApiUrl string
var firetailApiToken string

func loadEnvVars() {
	var firetailApiUrlSet bool
	firetailApiUrl, firetailApiUrlSet = os.LookupEnv("FIRETAIL_API_URL")
	if !firetailApiUrlSet {
		firetailApiUrl = DefaultFiretailApiUrl
	}
	firetailApiToken = os.Getenv("FIRETAIL_API_TOKEN")	
}

func main() {
	loadEnvVars()
	lambda.Start(Handler)
}
