package hfapigo

import (
	"io"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/api"
)

// NOTE: This file keeps small, stable interfaces in the root package so users
// interact with hfapigo directly while implementations stay internal.

// ChatService provides methods for interacting with chat completion endpoints.
// This interface is kept stable; new functionality is added via extended interfaces.
//
// Example:
//
//	type ChatServiceV2 interface {
//		ChatService
//		Stream(prompt string, opts ...RequestOption) (<-chan api.ChatResponse, error)
//	}
//
// For tests, it is recommended to use WithHTTPClient to inject a custom transport
// instead of mocking this interface.
type ChatService interface {
	Complete(prompt string, opts ...api.RequestOption) (api.ChatResponse, error)
}

// RawService provides methods for sending raw HTTP requests, with optional SDK error interpretation.
// This interface is kept stable; new functionality is added via extended interfaces.
//
// Example:
//
//	type RawServiceV2 interface {
//		RawService
//		DoWithRetry(requestBody []byte, method string, path string, opts ...RequestOption) (*http.Response, error)
//	}
//
// For tests, it is recommended to use WithHTTPClient to inject a custom transport
// instead of mocking this interface.
type RawService interface {
	// Do performs a raw HTTP request with a byte slice body and applies SDK error interpretation on non-2xx responses.
	// The caller must close resp.Body on success.
	Do(requestBody []byte, method string, path string, opts ...api.RequestOption) (*http.Response, error)
	// DoRaw performs a raw HTTP request with a byte slice body without translating non-2xx responses into SDK errors.
	// The caller must close resp.Body on success.
	DoRaw(requestBody []byte, method string, path string, opts ...api.RequestOption) (*http.Response, error)
	// DoReader performs a raw HTTP request with an io.Reader body and applies SDK error interpretation on non-2xx responses.
	// The caller must close resp.Body on success.
	DoReader(requestBody io.Reader, method string, path string, opts ...api.RequestOption) (*http.Response, error)
	// DoRawReader performs a raw HTTP request with an io.Reader body without translating non-2xx responses into SDK errors.
	// The caller must close resp.Body on success.
	DoRawReader(requestBody io.Reader, method string, path string, opts ...api.RequestOption) (*http.Response, error)
}
