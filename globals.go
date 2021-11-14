package hfapigo

import (
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
