package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const APIBaseURL = "https://api-inference.huggingface.co/models/"

var APIKey = func() string { return "" }

func SetAPIKey(key string) {
	APIKey = func() string { return key }
}

const (
	AuthHeaderKey    = "Authorization"
	AuthHeaderPrefix = "Bearer "
)

func setAuthorizationHeader(req *http.Request) *http.Request {
	if req == nil {
		return req
	}
	if APIKey() != "" {
		req.Header.Set(AuthHeaderKey, AuthHeaderPrefix+APIKey())
	}
	return req
}

// MakeHFAPIRequest builds and sends an HTTP POST request to the given model
// using the provided JSON body. If the request is successful, returns the
// response JSON and a nil error. If the request fails, returns an empty slice
// and an error describing the failure.
func MakeHFAPIRequest(jsonBody []byte, model string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, APIBaseURL+model, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthorizationHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = checkRespForError(respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

type apiError struct {
	Error string `json:"error,omitempty"`
}

type apiErrors struct {
	Errors []string `json:"error,omitempty"`
}

// Checks for errors in the API response and returns them if
// found.
func checkRespForError(respJSON []byte) error {
	// Check for single error
	{
		buf := make([]byte, len(respJSON))
		copy(buf, respJSON)
		apiErr := apiError{}
		json.Unmarshal(buf, &apiErr)
		if apiErr.Error != "" {
			return errors.New(string(respJSON))
		}
	}

	// Check for multiple errors
	{
		buf := make([]byte, len(respJSON))
		copy(buf, respJSON)
		apiErrs := apiErrors{}
		json.Unmarshal(buf, &apiErrs)
		if apiErrs.Errors != nil {
			return errors.New(string(respJSON))
		}
	}

	return nil
}

func MakeHFAPIRequestWithMedia(model, mediaFile string) ([]byte, error) {
	buf, err := os.ReadFile(mediaFile)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, APIBaseURL+model, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	setAuthorizationHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = checkRespForError(respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
