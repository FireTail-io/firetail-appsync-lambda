package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendToFiretail(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)
		assert.Equal(t, "{\"query\":\"TEST_QUERY\",\"request_id\":\"TEST_ID\"}\n", string(requestBody))
		wg.Done()
		w.Write([]byte(`{"message":"success"}`))
	}))

	testQuery := "TEST_QUERY"
	err := SendToFiretail(map[string]*FiretailLog{
		"TEST_ID": {
			Query:     &testQuery,
			RequestID: "TEST_ID",
		},
	}, testServer.URL, "TEST_KEY")
	require.Nil(t, err)

	wg.Wait()
}

func TestSendToFiretailBadServer(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)
		assert.Equal(t, "{\"query\":\"TEST_QUERY\",\"request_id\":\"TEST_ID\"}\n", string(requestBody))
		wg.Done()
		w.Write([]byte(`{"message":"fail"}`))
	}))

	testQuery := "TEST_QUERY"
	err := SendToFiretail(map[string]*FiretailLog{
		"TEST_ID": {
			Query:     &testQuery,
			RequestID: "TEST_ID",
		},
	}, testServer.URL, "TEST_KEY")
	require.NotNil(t, err)
	assert.Equal(t, "got err response from firetail api: map[message:fail]", err.Error())

	wg.Wait()
}

func TestSendToFiretailNoServer(t *testing.T) {
	testQuery := "TEST_QUERY"
	err := SendToFiretail(map[string]*FiretailLog{
		"TEST_ID": {
			Query:     &testQuery,
			RequestID: "TEST_ID",
		},
	}, "http://127.0.0.1:0", "TEST_KEY")
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Post \"http://127.0.0.1:0\": dial tcp 127.0.0.1:0")
}

func TestSendToFiretailBadUrl(t *testing.T) {
	testQuery := "TEST_QUERY"
	err := SendToFiretail(map[string]*FiretailLog{
		"TEST_ID": {
			Query:     &testQuery,
			RequestID: "TEST_ID",
		},
	}, "\n", "TEST_KEY")
	require.NotNil(t, err)
	assert.Equal(t, "parse \"\\n\": net/url: invalid control character in URL", err.Error())
}
