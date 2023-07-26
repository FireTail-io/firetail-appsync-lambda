package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVars(t *testing.T) {
	const MockFiretailApiUrl = "MOCK_FIRETAIL_API_URL"
	const MockFiretailApiToken = "MOCK_FIRETAIL_API_TOKEN"

	t.Setenv("FIRETAIL_API_URL", MockFiretailApiUrl)
	t.Setenv("FIRETAIL_API_TOKEN", MockFiretailApiToken)

	loadEnvVars()

	assert.Equal(t, firetailApiUrl, MockFiretailApiUrl)
	assert.Equal(t, firetailApiToken, MockFiretailApiToken)
}

func TestLoadEnvVarsApiUrlUnset(t *testing.T) {
	const MockFiretailApiToken = "MOCK_FIRETAIL_API_TOKEN"

	t.Setenv("FIRETAIL_API_TOKEN", MockFiretailApiToken)

	loadEnvVars()

	assert.Equal(t, firetailApiUrl, DefaultFiretailApiUrl)
	assert.Equal(t, firetailApiToken, MockFiretailApiToken)
}