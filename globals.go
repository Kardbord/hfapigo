package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

func SetAuthorizationHeader(req *http.Request) *http.Request {
	if req == nil {
		return req
	}
	if APIKey() != "" {
		req.Header.Set(AuthHeaderKey, AuthHeaderPrefix+APIKey())
	}
	return req
}

func BuildHFAPIRequest(jsonBody []byte, url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	req.Header.Set("Content-Type", "application/json")
	SetAuthorizationHeader(req)

	return req, nil
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
