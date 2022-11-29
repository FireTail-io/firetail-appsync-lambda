package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.CloudwatchLogsEvent) error {
	logsData, err := event.AWSLogs.Parse()
	if err != nil {
		return errors.WithMessage(err, "err parsing CloudwatchLogsEvent")
	}

	firetailLogs := map[string]*FiretailLog{}

	for _, logEvent := range logsData.LogEvents {
		// All of the logs that we care about in JSON format have a `logType` and `requestId` field.
		// We don't care about any of their other values.
		type JsonLog struct {
			LogType   string `json:"logType"`
			RequestID string `json:"requestId"`
		}
		var jsonLog JsonLog
		err = json.Unmarshal([]byte(logEvent.Message), &jsonLog)

		// Extract the requestID - if it failed to marshal, it'll just be the first element when
		// the logEvent.Message is split by spaces.
		var requestID string
		if err != nil {
			requestID = strings.Split(logEvent.Message, " ")[0]
		} else {
			requestID = jsonLog.RequestID
		}

		// Extract the logType - if the logEvent failed to marshal as JSON, then it's plaintext
		var logType LogMessageType
		if err != nil {
			logType = Plaintext
		} else {
			logType = LogMessageType(jsonLog.LogType)
		}

		// Get the existing firetailLog for this requestID, and if it doesn't exist then create it
		firetailLog, firetailLogExists := firetailLogs[requestID]
		if !firetailLogExists {
			firetailLog = &FiretailLog{
				RequestID: requestID,
			}
			firetailLogs[requestID] = firetailLog
		}

		// Add this event message to the firetailLog for the corresponding request ID!
		err = firetailLog.AddEventMessage(logType, &logEvent)
		if err != nil {
			log.Println("Err adding event message to firetail log:", err.Error())
		}
	}

	// Print all the FiretailLogs we built up
	for requestID, firetailLog := range firetailLogs {
		logBytes, err := json.Marshal(firetailLog)
		if err != nil {
			log.Printf("Err marshalling firetail log for request ID %s: %s", requestID, err.Error())
			continue
		}
		log.Println(string(logBytes))
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
