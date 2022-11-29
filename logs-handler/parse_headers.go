package main

import (
	"errors"
	"fmt"
	"strings"
)

// Trims the open and closing curly braces from the headers string we get from Cloudwatch
func trimHeadersString(headersString string) (string, error) {
	if headersString[0] == '{' {
		headersString = headersString[1:]
	} else {
		return "", errors.New("headers string should start with '{'")
	}
	if headersString[len(headersString)-1] == '}' {
		headersString = headersString[:len(headersString)-1]
	} else {
		return "", errors.New("headers string should end with '}'")
	}
	return headersString, nil
}

// (ー_ー)!! I am ashamed of this, please do not talk to me about it
// Parses the headers string we get from CloudWatch into a map[string]string.
func parseHeaders(headersString string) (map[string][]string, error) {
	headersString, err := trimHeadersString(headersString)
	if err != nil {
		return nil, err
	}
	headers := map[string][]string{}
	parts := strings.Split(headersString, ",")
	for _, part := range parts {
		subparts := strings.SplitN(part, "=", 2)
		if len(subparts) != 2 {
			return nil, fmt.Errorf("header had !=2 subparts when split by first '=': %s", part)
		}
		headers[strings.Trim(subparts[0], " ")] = []string{subparts[1]}
	}
	return headers, nil
}

// (ー_ー)!! I am ashamed of this, please do not talk to me about it
// Parses the multivalue headers string we get from CloudWatch into a map[string][]string.
func parseMultivalueHeaders(headersString string) (map[string][]string, error) {
	headersString, err := trimHeadersString(headersString)
	if err != nil {
		return nil, err
	}
	headers := map[string][]string{}
	for _, part := range strings.Split(headersString, "],") {
		subparts := strings.SplitN(part, "=[", 2)
		if len(subparts) != 2 {
			return nil, fmt.Errorf("multivalue header had !=2 subparts when split by first '=[': %s", part)
		}
		headers[strings.Trim(subparts[0], " ")] = strings.Split(subparts[1], ", ")
	}
	return headers, nil
}
