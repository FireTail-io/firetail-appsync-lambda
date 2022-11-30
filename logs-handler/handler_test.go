package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	testData := "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"
	expectedPayload := "{\"executionSummary\":{\"duration\":62176672,\"logType\":\"ExecutionSummary\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"startTime\":\"2022-11-30T11:03:56.035126Z\",\"endTime\":\"2022-11-30T11:03:56.097302Z\",\"parsing\":{\"startOffset\":56801,\"duration\":49156},\"version\":1,\"validation\":{\"startOffset\":132790,\"duration\":73757},\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\"},\"query\":\"mutation MyMutation {\\n  createPost(input: {title: \\\"A Test Post\\\"}) {\\n    id\\n  }\\n}\\n\\nquery MyQuery {\\n  getPost(id: \\\"a5422778-e5bd-4b29-8c26-0e44a921aba1\\\") {\\n    id\\n    title\\n  }\\n  listPosts(limit: 10) {\\n    items {\\n      id\\n    }\\n  }\\n}\\n, Operation: MyQuery, Variables: {}\",\"request_id\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"requestHeaders\":{\"accept\":[\"application/json\",\"text/plain\",\"*/*\"],\"accept-encoding\":[\"gzip\",\"deflate\",\"br\"],\"accept-language\":[\"en-GB,en-US;q=0.9,en;q=0.8\"],\"cloudfront-forwarded-proto\":[\"https\"],\"cloudfront-is-desktop-viewer\":[\"true\"],\"cloudfront-is-mobile-viewer\":[\"false\"],\"cloudfront-is-smarttv-viewer\":[\"false\"],\"cloudfront-is-tablet-viewer\":[\"false\"],\"cloudfront-viewer-asn\":[\"1136\"],\"cloudfront-viewer-country\":[\"NL\"],\"content-length\":[\"309\"],\"content-type\":[\"application/json\"],\"host\":[\"c5dz3eobtjce7p4emob3koivpa.appsync-api.eu-west-1.amazonaws.com\"],\"origin\":[\"https://eu-west-1.console.aws.amazon.com\"],\"referer\":[\"https://eu-west-1.console.aws.amazon.com/\"],\"sec-ch-ua\":[\"\\\"Google Chrome\\\";v=\\\"107\\\"\",\"\\\"Chromium\\\";v=\\\"107\\\"\",\"\\\"Not=A?Brand\\\";v=\\\"24\\\"\"],\"sec-ch-ua-mobile\":[\"?0\"],\"sec-ch-ua-platform\":[\"\\\"macOS\\\"\"],\"sec-fetch-dest\":[\"empty]\"],\"sec-fetch-mode\":[\"cors\"],\"sec-fetch-site\":[\"cross-site\"],\"user-agent\":[\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\",\"like Gecko) Chrome/107.0.0.0 Safari/537.36\"],\"via\":[\"2.0 a8b68315e1e2575143f97748ffbb29a0.cloudfront.net (CloudFront)\"],\"x-amz-cf-id\":[\"hv4XlmXAktT6V5z2wpKEwjzcrXrKATjdDtYjltpLcIG_HDmBcdxzSw==\"],\"x-amz-user-agent\":[\"AWS-Console-AppSync/\"],\"x-amzn-requestid\":[\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\"],\"x-amzn-trace-id\":[\"Root=1-6387389b-73d623b74d5d9f671677aa8a\"],\"x-api-key\":[\"****mgidri\"],\"x-forwarded-for\":[\"77.173.29.29\",\"15.158.40.17\"],\"x-forwarded-port\":[\"443\"],\"x-forwarded-proto\":[\"https\"]},\"requestMappings\":[{\"logType\":\"RequestMapping\",\"path\":[\"getPost\"],\"fieldName\":\"getPost\",\"resolverArn\":\"arn:aws:appsync:eu-west-1:453671210445:apis/lcyxyv2rungh7gdht7ylrudsjy/types/Query/resolvers/getPost\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"context\":{\"arguments\":{\"id\":\"a5422778-e5bd-4b29-8c26-0e44a921aba1\"},\"stash\":{},\"outErrors\":[]},\"fieldInError\":false,\"errors\":[],\"parentType\":\"Query\",\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\",\"transformedTemplate\":\"{\\n  \\\"version\\\": \\\"2017-02-28\\\",\\n  \\\"operation\\\": \\\"GetItem\\\",\\n  \\\"key\\\": {\\n    \\\"id\\\": {\\\"S\\\":\\\"a5422778-e5bd-4b29-8c26-0e44a921aba1\\\"},\\n  },\\n}\"},{\"logType\":\"RequestMapping\",\"path\":[\"listPosts\"],\"fieldName\":\"listPosts\",\"resolverArn\":\"arn:aws:appsync:eu-west-1:453671210445:apis/lcyxyv2rungh7gdht7ylrudsjy/types/Query/resolvers/listPosts\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"context\":{\"arguments\":{\"limit\":10},\"stash\":{},\"outErrors\":[]},\"fieldInError\":false,\"errors\":[],\"parentType\":\"Query\",\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\",\"transformedTemplate\":\"{\\n  \\\"version\\\": \\\"2017-02-28\\\",\\n  \\\"operation\\\": \\\"Scan\\\",\\n  \\\"filter\\\":  null ,\\n  \\\"limit\\\": 10,\\n  \\\"nextToken\\\": null,\\n}\"}],\"requestSummary\":{\"logType\":\"RequestSummary\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\",\"statusCode\":200,\"latency\":89000000},\"responseHeaders\":{\"Content-Type\":\"application/json; charset=UTF-8\"},\"responseMappings\":[{\"logType\":\"ResponseMapping\",\"path\":[\"getPost\"],\"fieldName\":\"getPost\",\"resolverArn\":\"arn:aws:appsync:eu-west-1:453671210445:apis/lcyxyv2rungh7gdht7ylrudsjy/types/Query/resolvers/getPost\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"context\":{\"arguments\":{\"id\":\"a5422778-e5bd-4b29-8c26-0e44a921aba1\"},\"result\":{\"id\":\"a5422778-e5bd-4b29-8c26-0e44a921aba1\",\"title\":\"A Test Post\"},\"stash\":{},\"outErrors\":[]},\"fieldInError\":false,\"errors\":[],\"parentType\":\"Query\",\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\",\"transformedTemplate\":\"{id=a5422778-e5bd-4b29-8c26-0e44a921aba1, title=A Test Post}\"},{\"logType\":\"ResponseMapping\",\"path\":[\"listPosts\"],\"fieldName\":\"listPosts\",\"resolverArn\":\"arn:aws:appsync:eu-west-1:453671210445:apis/lcyxyv2rungh7gdht7ylrudsjy/types/Query/resolvers/listPosts\",\"requestId\":\"0eff56b0-5caf-47a8-9e8d-7a174e93a8db\",\"context\":{\"arguments\":{\"limit\":10},\"result\":{\"items\":[{\"id\":\"a5422778-e5bd-4b29-8c26-0e44a921aba1\",\"title\":\"A Test Post\"},{\"id\":\"0daae63f-46ab-4631-8db4-e2c36a7f00eb\",\"title\":\"A Second Test Post\"}],\"scannedCount\":2},\"stash\":{},\"outErrors\":[]},\"fieldInError\":false,\"errors\":[],\"parentType\":\"Query\",\"graphQLAPIId\":\"lcyxyv2rungh7gdht7ylrudsjy\",\"transformedTemplate\":\"{items=[{id=a5422778-e5bd-4b29-8c26-0e44a921aba1, title=A Test Post}, {id=0daae63f-46ab-4631-8db4-e2c36a7f00eb, title=A Second Test Post}], nextToken=null, scannedCount=2, startedAt=null}\"}]}\n"

	wg := &sync.WaitGroup{}
	wg.Add(1)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)
		assert.Equal(t, expectedPayload, string(bodyBytes))
		w.Write([]byte(`{"message":"success"}`))
		wg.Done()
	}))
	firetailApiUrl = testServer.URL

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := Handler(ctx, events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: testData,
		},
	})
	require.Nil(t, err)

	wg.Wait()
}

func TestHandlerWithNoData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := Handler(ctx, events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: "",
		},
	})
	require.NotNil(t, err)
	assert.Equal(t, "err parsing CloudwatchLogsEvent: EOF", err.Error())
}

func TestHandlerWithNoRelevantData(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	testData := `{
		"logEvents": [{
			"id": "TEST_ID",
			"message": "TEST_ID Begin Request"
		}]
	}`

	// Gzip the data
	var gzipBytes bytes.Buffer
	gzipper := gzip.NewWriter(&gzipBytes)
	_, err := gzipper.Write([]byte(testData))
	require.Nil(t, err)
	err = gzipper.Close()
	require.Nil(t, err)

	// b64 encode the gzipped data
	b64TestData := base64.StdEncoding.EncodeToString(gzipBytes.Bytes())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Despite having no mock Firetail logs API in this test, this handler call shouldn't
	// err as there should be no logs to send to Firetail.
	err = Handler(ctx, events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: b64TestData,
		},
	})
	require.Nil(t, err)

	assert.Contains(t, string(logOutput.Bytes()), "Generated no Firetail logs from this batch. Exiting...")
}

func TestHandlerWithMalformedData(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	testData := `{
		"logEvents": [{
			"id": "TEST_ID",
			"message": "TEST_ID Request Headers: {TEST_HEADER=TEST_VALUE}"
		}]
	}`

	// Gzip the data
	var gzipBytes bytes.Buffer
	gzipper := gzip.NewWriter(&gzipBytes)
	_, err := gzipper.Write([]byte(testData))
	require.Nil(t, err)
	err = gzipper.Close()
	require.Nil(t, err)

	// b64 encode the gzipped data
	b64TestData := base64.StdEncoding.EncodeToString(gzipBytes.Bytes())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Despite having no mock Firetail logs API in this test, this handler call shouldn't
	// err as there should be no logs to send to Firetail.
	err = Handler(ctx, events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: b64TestData,
		},
	})
	require.Nil(t, err)

	assert.Contains(t, string(logOutput.Bytes()), "Errs extracting Firetail logs: 1 error occurred:\n\t* err adding event message to firetail log: multivalue header had !=2 subparts when split by first '=[': TEST_HEADER=TEST_VALUE")
	assert.Contains(t, string(logOutput.Bytes()), "Generated no Firetail logs from this batch. Exiting...")
}
