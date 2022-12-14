package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrimHeadersStringNoOpenBrace(t *testing.T) {
	testString := "{Content-Type=application/json; charset=UTF-8"
	result, err := trimHeadersString(testString)
	assert.Zero(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "headers string should end with '}'", err.Error())
}

func TestTrimHeadersStringNoCloseBrace(t *testing.T) {
	testString := "{Content-Type=application/json; charset=UTF-8"
	result, err := trimHeadersString(testString)
	assert.Zero(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "headers string should end with '}'", err.Error())
}

func TestParseHeaders(t *testing.T) {
	testString := "{Content-Type=application/json; charset=UTF-8}"

	result, err := parseHeaders(testString)
	require.Nil(t, err)

	assert.Equal(t, result, map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
	})
}

func TestParseHeadersMalformed(t *testing.T) {
	testString := "{Content-Type:application/json; charset:UTF-8}"
	result, err := parseHeaders(testString)
	assert.Nil(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "header had !=2 subparts when split by first '=': Content-Type:application/json; charset:UTF-8", err.Error())
}

func TestParseMultiValueHeadersMalformed(t *testing.T) {
	testString := "{Content-Type:[application/json; charset:UTF-8]}"
	result, err := parseMultivalueHeaders(testString)
	assert.Nil(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "multivalue header had !=2 subparts when split by first '=[': Content-Type:[application/json; charset:UTF-8]", err.Error())
}

func TestParseHeadersNoOpenBrace(t *testing.T) {
	testString := "Content-Type=application/json; charset=UTF-8}"
	result, err := parseHeaders(testString)
	assert.Nil(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "headers string should start with '{'", err.Error())
}

func TestParseMultiValueHeadersNoOpenBrace(t *testing.T) {
	testString := "Content-Type=application/json; charset=UTF-8}"
	result, err := parseMultivalueHeaders(testString)
	assert.Nil(t, result)
	require.NotNil(t, err)
	assert.Equal(t, "headers string should start with '{'", err.Error())
}

func TestParseMultivalueHeaders(t *testing.T) {
	testString := `{content-length=[322], referer=[https://eu-west-1.console.aws.amazon.com/], cloudfront-viewer-country=[NL], sec-fetch-site=[cross-site], x-amzn-requestid=[832cf953-06db-4b07-9e4f-8d5f8a7691e2], origin=[https://eu-west-1.console.aws.amazon.com], x-amz-user-agent=[AWS-Console-AppSync/], x-forwarded-port=[443], via=[2.0 00f66bc6263192200d1a0cdb83e969f8.cloudfront.net (CloudFront)], sec-ch-ua-mobile=[?0], cloudfront-viewer-asn=[1136], cloudfront-is-desktop-viewer=[true], host=[c5dz3eobtjce7p4emob3koivpa.appsync-api.eu-west-1.amazonaws.com], content-type=[application/json], sec-fetch-mode=[cors], x-forwarded-proto=[https], accept-language=[en-GB,en-US;q=0.9,en;q=0.8], x-forwarded-for=[77.173.29.29, 15.158.40.15], accept=[application/json, text/plain, */*], cloudfront-is-smarttv-viewer=[false], sec-ch-ua=["Google Chrome";v="107", "Chromium";v="107", "Not=A?Brand";v="24"], x-amzn-trace-id=[Root=1-6384e16b-3d08b227276904141dcd192b], cloudfront-is-tablet-viewer=[false], x-api-key=[****mgidri], sec-ch-ua-platform=["macOS"], cloudfront-forwarded-proto=[https], accept-encoding=[gzip, deflate, br], x-amz-cf-id=[uNT3aeWVqWr12Wz6b94neGhKRLPPa_o65dR40NZx4KQc08Lotbd_3g==], user-agent=[Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36], cloudfront-is-mobile-viewer=[false], sec-fetch-dest=[empty]}`

	result, err := parseMultivalueHeaders(testString)
	require.Nil(t, err)

	assert.Equal(
		t,
		map[string][]string{
			"accept":                       {"application/json, text/plain, */*"},
			"accept-encoding":              {"gzip, deflate, br"},
			"accept-language":              {"en-GB,en-US;q=0.9,en;q=0.8"},
			"cloudfront-forwarded-proto":   {"https"},
			"cloudfront-is-desktop-viewer": {"true"},
			"cloudfront-is-mobile-viewer":  {"false"},
			"cloudfront-is-smarttv-viewer": {"false"},
			"cloudfront-is-tablet-viewer":  {"false"},
			"cloudfront-viewer-asn":        {"1136"},
			"cloudfront-viewer-country":    {"NL"},
			"content-type":                 {"application/json"},
			"host":                         {"c5dz3eobtjce7p4emob3koivpa.appsync-api.eu-west-1.amazonaws.com"},
			"origin":                       {"https://eu-west-1.console.aws.amazon.com"},
			"referer":                      {"https://eu-west-1.console.aws.amazon.com/"},
			"sec-ch-ua":                    {"\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\""},
			"sec-ch-ua-mobile":             {"?0"},
			"sec-ch-ua-platform":           {"\"macOS\""},
			"sec-fetch-dest":               {"empty"},
			"sec-fetch-mode":               {"cors"},
			"sec-fetch-site":               {"cross-site"},
			"user-agent":                   {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
			"via":                          {"2.0 00f66bc6263192200d1a0cdb83e969f8.cloudfront.net (CloudFront)"},
			"x-amz-cf-id":                  {"uNT3aeWVqWr12Wz6b94neGhKRLPPa_o65dR40NZx4KQc08Lotbd_3g=="},
			"x-amz-user-agent":             {"AWS-Console-AppSync/"},
			"x-amzn-requestid":             {"832cf953-06db-4b07-9e4f-8d5f8a7691e2"},
			"x-amzn-trace-id":              {"Root=1-6384e16b-3d08b227276904141dcd192b"},
			"x-api-key":                    {"****mgidri"},
			"x-forwarded-for":              {"77.173.29.29, 15.158.40.15"},
			"x-forwarded-port":             {"443"},
			"x-forwarded-proto":            {"https"},
			"content-length":               {"322"},
		},
		result,
	)
}
