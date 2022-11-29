package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var firetailApiKey string

func Handler(ctx context.Context, event events.CloudwatchLogsEvent) error {
	logsData, err := event.AWSLogs.Parse()
	if err != nil {
		return errors.WithMessage(err, "err parsing CloudwatchLogsEvent")
	}

	firetailLogs, err := ExtractFiretailLogs(&logsData)
	if err != nil {
		log.Println("Errs extracting Firetail logs:", err.Error())
	}
	if firetailLogs == nil || len(firetailLogs) == 0 {
		log.Println("Generated no Firetail logs from this batch. Exiting...")
		return nil
	}

	for requestID, firetailLog := range firetailLogs {
		if !firetailLog.IsPopulated() {
			log.Printf("No useful information was extracted from this Cloudwatch logs batch for request ID %s", requestID)
			continue
		}
		logBytes, err := json.Marshal(firetailLog)
		if err != nil {
			log.Printf("Err marshalling firetail log for request ID %s: %s", requestID, err.Error())
			continue
		}
		log.Println(string(logBytes))
	}

	return SendToFiretail(firetailLogs, firetailApiKey)
}

func main() {
	firetailApiKey = os.Getenv("FIRETAIL_API_KEY")
	lambda.Start(Handler)
}
