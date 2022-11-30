package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

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
		logBytes, err := json.Marshal(firetailLog)
		if err != nil {
			log.Printf("Err marshalling firetail log for request ID %s: %s", requestID, err.Error())
			continue
		}
		log.Println(string(logBytes))
	}

	return SendToFiretail(firetailLogs, firetailApiUrl, firetailApiKey)
}
