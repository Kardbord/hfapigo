package hfapigo

import (
	"net/http"
	"sync"
)

const APIBaseURL = "https://api-inference.huggingface.co/models/"

var (
	apiKeyMutex = sync.RWMutex{}

	APIKey = func() string {
		apiKeyMutex.RLock()
		defer apiKeyMutex.RUnlock()
		return ""
	}
)

func SetAPIKey(key string) {
	apiKeyMutex.Lock()
	defer apiKeyMutex.Unlock()
	APIKey = func() string {
		apiKeyMutex.RLock()
		defer apiKeyMutex.RUnlock()
		return key
	}
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
