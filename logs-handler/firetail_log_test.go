package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPopulatedEmpty(t *testing.T) {
	firetailLog := FiretailLog{
		RequestID: "TEST_REQUEST",
	}
	isPopulated := firetailLog.IsPopulated()
	assert.False(t, isPopulated)
}

func TestIsPopulatedPopulated(t *testing.T) {
	testQuery := "TEST_QUERY"
	firetailLog := FiretailLog{
		RequestID: "TEST_REQUEST",
		Query:     &testQuery,
	}
	isPopulated := firetailLog.IsPopulated()
	assert.True(t, isPopulated)
}

func TestAddEventMessageRequestMapping(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(RequestMapping, &testEvent)
	require.Nil(t, err)
	assert.Equal(t, &[]json.RawMessage{json.RawMessage(testEvent.Message)}, testLog.RequestMappings)
}

func TestAddEventMessageResponseMapping(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(ResponseMapping, &testEvent)
	require.Nil(t, err)
	assert.Equal(t, &[]json.RawMessage{json.RawMessage(testEvent.Message)}, testLog.ResponseMappings)
}

func TestAddEventMessageExecutionSummary(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(ExecutionSummary, &testEvent)
	require.Nil(t, err)
	require.NotNil(t, testLog.ExecutionSummary)
	assert.Equal(t, json.RawMessage(testEvent.Message), *testLog.ExecutionSummary)
}

func TestAddEventMessageRequestSummary(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(RequestSummary, &testEvent)
	require.Nil(t, err)
	require.NotNil(t, testLog.RequestSummary)
	assert.Equal(t, json.RawMessage(testEvent.Message), *testLog.RequestSummary)
}

func TestAddEventMessageBeginTracing(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(BeginTracing, &testEvent)
	require.Nil(t, err)
	require.False(t, testLog.IsPopulated())
}

func TestAddEventMessageMalformedPlaintext(t *testing.T) {
	testEvent := events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_MESSAGE",
	}
	testLog := &FiretailLog{}
	err := testLog.AddEventMessage(Plaintext, &testEvent)
	require.NotNil(t, err)
	assert.Equal(t, "plaintext logEventMessage had 1 parts when split by ' ' but needs >= 2", err.Error())
}

func TestAddPlaintextEventMessageQuery(t *testing.T) {
	testEvent := &events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_ID GraphQL Query: TEST_QUERY",
	}
	testLog := &FiretailLog{}
	err := testLog.addPlaintextEventMessage(testEvent)
	require.Nil(t, err)
	require.NotNil(t, testLog.Query)
	assert.Equal(t, "TEST_QUERY", *testLog.Query)
}

func TestAddPlaintextEventRequestHeaders(t *testing.T) {
	testEvent := &events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_ID Request Headers: {TEST_HEADER=[TEST_VALUE_1, TEST_VALUE_2]}",
	}
	testLog := &FiretailLog{}
	err := testLog.addPlaintextEventMessage(testEvent)
	require.Nil(t, err)
	require.NotNil(t, testLog.RequestHeaders)
	assert.Equal(t, "{\"TEST_HEADER\":[\"TEST_VALUE_1\",\"TEST_VALUE_2]\"]}", string(*testLog.RequestHeaders))
}

func TestAddPlaintextEventResponseHeaders(t *testing.T) {
	testEvent := &events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_ID Response Headers: {TEST_HEADER=TEST_VALUE}",
	}
	testLog := &FiretailLog{}
	err := testLog.addPlaintextEventMessage(testEvent)
	require.Nil(t, err)
	require.NotNil(t, testLog.ResponseHeaders)
	assert.Equal(t, "{\"TEST_HEADER\":\"TEST_VALUE\"}", string(*testLog.ResponseHeaders))
}

func TestAddPlaintextEventBeginRequest(t *testing.T) {
	testEvent := &events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_ID Begin Request",
	}
	testLog := &FiretailLog{}
	err := testLog.addPlaintextEventMessage(testEvent)
	require.Nil(t, err)
	assert.False(t, testLog.IsPopulated())
}

func TestAddPlaintextEventNoMatches(t *testing.T) {
	testEvent := &events.CloudwatchLogsLogEvent{
		ID:        "TEST_ID",
		Timestamp: 3142,
		Message:   "TEST_ID Test Event",
	}
	testLog := &FiretailLog{}
	err := testLog.addPlaintextEventMessage(testEvent)
	require.NotNil(t, err)
	assert.Equal(t, "plaintext logEventMessage matched no plaintext log prefixes: TEST_ID Test Event", err.Error())
}
