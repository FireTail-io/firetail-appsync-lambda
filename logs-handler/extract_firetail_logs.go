package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func ExtractFiretailLogs(logsData *events.CloudwatchLogsData) (map[string]*FiretailLog, error) {
	firetailLogs := map[string]*FiretailLog{}
	var errs error

	for _, logEvent := range logsData.LogEvents {
		// All of the logs that we care about in JSON format have a `logType` and `requestId` field.
		// We don't care about any of their other values.
		type JsonLog struct {
			LogType   string `json:"logType"`
			RequestID string `json:"requestId"`
		}
		var jsonLog JsonLog
		err := json.Unmarshal([]byte(logEvent.Message), &jsonLog)

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
		}

		// Add this event message to the firetailLog for the corresponding request ID!
		err = firetailLog.AddEventMessage(logType, &logEvent)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessage(err, "err adding event message to firetail log"))
			continue
		}

		firetailLogs[requestID] = firetailLog
	}

	return firetailLogs, errs
}
