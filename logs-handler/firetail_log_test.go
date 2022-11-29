package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
