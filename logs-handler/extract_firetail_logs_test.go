package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFiretailLogsSingleQueryEvent(t *testing.T) {
	testData := &events.CloudwatchLogsData{
		LogEvents: []events.CloudwatchLogsLogEvent{{
			ID:        "TEST_ID",
			Timestamp: 0,
			Message:   "TEST_ID GraphQL Query: TEST_QUERY",
		}},
	}

	logs, err := ExtractFiretailLogs(testData)
	require.Nil(t, err)

	require.Contains(t, logs, "TEST_ID")
	require.NotNil(t, logs["TEST_ID"].Query)
	assert.Equal(t, "TEST_QUERY", *(logs["TEST_ID"].Query))
}

func TestExtractFiretailLogsSingleRequestMappingEvent(t *testing.T) {
	testData := &events.CloudwatchLogsData{
		LogEvents: []events.CloudwatchLogsLogEvent{{
			ID:        "TEST_ID",
			Timestamp: 0,
			Message:   `{"logType":"RequestMapping","requestId":"TEST_ID"}`,
		}},
	}

	logs, err := ExtractFiretailLogs(testData)
	require.Nil(t, err)

	require.Contains(t, logs, "TEST_ID")
	require.NotNil(t, logs["TEST_ID"].RequestMappings)
	require.Len(t, *logs["TEST_ID"].RequestMappings, 1)
	assert.Equal(t, "{\"logType\":\"RequestMapping\",\"requestId\":\"TEST_ID\"}", string((*logs["TEST_ID"].RequestMappings)[0]))
}

func TestExtractFiretailLogsSingleMalformedEvent(t *testing.T) {
	testData := &events.CloudwatchLogsData{
		LogEvents: []events.CloudwatchLogsLogEvent{{
			ID:        "TEST_ID",
			Timestamp: 0,
			Message:   `TEST_ID Invalid Prefix`,
		}},
	}

	logs, err := ExtractFiretailLogs(testData)
	assert.Len(t, logs, 0)
	require.NotNil(t, err)
	assert.Equal(t, "1 error occurred:\n\t* err adding event message to firetail log: plaintext logEventMessage matched no plaintext log prefixes: TEST_ID Invalid Prefix\n\n", err.Error())
}
