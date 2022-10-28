package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const routeError = "/error"

type errorData struct {
	HVal        string `json:"h_val"`
	PeerId      string `json:"peer_id"`
	ErrorStatus string `json:"error_status"`
}

// function to send error message to status server
func SendError(errMsg, onlineServerDomain, peerId, hVal string) error {
	errData := new(errorData)
	errData.ErrorStatus = errMsg
	errData.PeerId = peerId
	errData.HVal = hVal
	errDataMarshaled, err := json.Marshal(errData)

	// reports to status server by making HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", onlineServerDomain, routeError), bytes.NewReader(errDataMarshaled))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to make new http request, reason: %v", err))
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to perform http.DefaultClient.Do(), reason: %v", err))
	}
	defer res.Body.Close()

	return nil
}
