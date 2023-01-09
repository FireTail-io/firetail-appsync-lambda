package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func SendToFiretail(firetailLogs map[string]*FiretailLog, apiUrl, apiToken string) error {
	reqBytes := []byte{}
	for _, firetailLog := range firetailLogs {
		logBytes, err := json.Marshal(*firetailLog)
		if err != nil {
			return err
		}
		reqBytes = append(reqBytes, logBytes...)
		reqBytes = append(reqBytes, '\n')
	}

	req, err := http.NewRequest(
		"POST",
		apiUrl,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("x-ft-api-key", apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if res["message"] != "success" {
		return errors.New(fmt.Sprintf("got err response from firetail api: %v", res))
	}

	return nil
}
