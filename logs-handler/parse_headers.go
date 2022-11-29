package main

import (
	"errors"
	"fmt"
	"strings"
)

// (ー_ー)!! I am ashamed of this, please do not talk to me about it
// Parses the headers string we get from CloudWatch into a map[string][]string.
// Handles single and multivalue headers strings, but not mixed.
func parseHeaders(headersString string) (map[string][]string, error) {
	if headersString[0] == '{' {
		headersString = headersString[1:]
	} else {
		return nil, errors.New("headers string should start with '{'")
	}

	if headersString[len(headersString)-1] == '}' {
		headersString = headersString[:len(headersString)-1]
	} else {
		return nil, errors.New("headers string should end with '}'")
	}

	headers := map[string][]string{}

	// If there's a closing square bracket we'll assume it's a multivalue headers string
	if strings.Contains(headersString, "]") {
		for _, part := range strings.Split(headersString, "],") {
			subparts := strings.SplitN(part, "=[", 2)
			if len(subparts) != 2 {
				return nil, fmt.Errorf("multivalue header had !=2 subparts when split by first '=[': %s", part)
			}
			headers[strings.Trim(subparts[0], " ")] = strings.Split(subparts[1], ", ")
		}
	} else {
		// Otherwise we'll assume each header has a single value
		parts := strings.Split(headersString, ",")
		for _, part := range parts {
			subparts := strings.SplitN(part, "=", 2)
			if len(subparts) != 2 {
				return nil, fmt.Errorf("header had !=2 subparts when split by first '=': %s", part)
			}
			headers[strings.Trim(subparts[0], " ")] = []string{subparts[1]}
		}
	}

	return headers, nil
}
