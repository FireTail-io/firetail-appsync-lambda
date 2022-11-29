package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type LogMessageType string

const (
	BeginRequest      LogMessageType = "BeginRequest"
	GraphQLQuery      LogMessageType = "GraphQLQuery"
	BeginExecution    LogMessageType = "BeginExecution"
	RequestMapping    LogMessageType = "RequestMapping"
	ResponseMapping   LogMessageType = "ResponseMapping"
	EndFieldExecution LogMessageType = "EndFieldExecution"
	ExecutionSummary  LogMessageType = "ExecutionSummary"
	BeginTracing      LogMessageType = "BeginTracing"
	EndTracing        LogMessageType = "EndTracing"
	RequestSummary    LogMessageType = "RequestSummary"
	RequestHeaders    LogMessageType = "RequestHeaders"
	ResponseHeaders   LogMessageType = "ResponseHeaders"
	TokensConsumed    LogMessageType = "TokensConsumed"
	EndRequest        LogMessageType = "EndRequest"
	Plaintext         LogMessageType = "Plaintext"
)

type FiretailLog struct {
	ExecutionSummary *json.RawMessage   `json:"executionSummary,omitempty"`
	Query            *string            `json:"query,omitempty"`
	RequestID        string             `json:"request_id"`
	RequestHeaders   *json.RawMessage   `json:"requestHeaders,omitempty"`
	RequestMappings  *[]json.RawMessage `json:"requestMappings,omitempty"`
	RequestSummary   *json.RawMessage   `json:"requestSummary,omitempty"`
	ResponseHeaders  *json.RawMessage   `json:"responseHeaders,omitempty"`
	ResponseMappings *[]json.RawMessage `json:"responseMappings,omitempty"`
}

func (f *FiretailLog) AddEventMessage(logType LogMessageType, logEvent *events.CloudwatchLogsLogEvent) error {
	var rawMessage json.RawMessage
	if logType != Plaintext {
		rawMessage = json.RawMessage(logEvent.Message)
	} else {
		return f.addPlaintextEventMessage(logEvent)
	}

	switch logType {
	case RequestMapping:
		if f.RequestMappings == nil {
			f.RequestMappings = &[]json.RawMessage{}
		}
		*f.RequestMappings = append(*f.RequestMappings, rawMessage)
		break

	case ResponseMapping:
		if f.ResponseMappings == nil {
			f.ResponseMappings = &[]json.RawMessage{}
		}
		*f.ResponseMappings = append(*f.ResponseMappings, rawMessage)
		break

	case ExecutionSummary:
		f.ExecutionSummary = &rawMessage
		break

	case RequestSummary:
		f.RequestSummary = &rawMessage
		break
	}

	return nil
}

func (f *FiretailLog) addPlaintextEventMessage(logEvent *events.CloudwatchLogsLogEvent) error {
	// Make sure it has enough parts
	logParts := strings.SplitN(logEvent.Message, " ", 2)
	if len(logParts) < 2 {
		return fmt.Errorf("plaintext logEventMessage had %d parts when split by ' ' but needs >= 2", len(logParts))
	}

	// Determine its type from its prefix & extract its payload
	plaintextLogPrefixes := map[string]LogMessageType{
		"Begin Request":       BeginRequest,
		"GraphQL Query: ":     GraphQLQuery,
		"Begin Execution - ":  BeginExecution,
		"End Field Execution": EndFieldExecution,
		"Begin Tracing":       BeginTracing,
		"End Tracing":         EndTracing,
		"Request Headers: ":   RequestHeaders,
		"Response Headers: ":  ResponseHeaders,
		"Tokens Consumed: ":   TokensConsumed,
		"End Request":         EndRequest,
	}
	var logType LogMessageType
	var logPayload string
	for logPrefix, potentialLogType := range plaintextLogPrefixes {
		if strings.HasPrefix(logParts[1], logPrefix) {
			logType = potentialLogType
			logParts = strings.SplitN(logEvent.Message, logPrefix, 2)
			if len(logParts) < 2 {
				logPayload = ""
			} else {
				logPayload = logParts[1]
			}
			break
		}
	}
	if logType == "" {
		return fmt.Errorf("plaintext logEventMessage matched no plaintext log prefixes: %s", logEvent.Message)
	}

	switch logType {
	case GraphQLQuery:
		f.Query = &logParts[1]
		break
	case RequestHeaders:
		jsonPayload, err := parseMultivalueHeaders(logPayload)
		if err != nil {
			return err
		}
		jsonPayloadBytes, err := json.Marshal(jsonPayload)
		if err != nil {
			return err
		}
		rawJson := json.RawMessage(jsonPayloadBytes)
		f.RequestHeaders = &rawJson
		break
	case ResponseHeaders:
		jsonPayload, err := parseHeaders(logPayload)
		if err != nil {
			return err
		}
		jsonPayloadBytes, err := json.Marshal(jsonPayload)
		if err != nil {
			return err
		}
		rawJson := json.RawMessage(jsonPayloadBytes)
		f.ResponseHeaders = &rawJson
		break
	}

	return nil
}
