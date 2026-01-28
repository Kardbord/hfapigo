package request

import (
	"bytes"
	"io"
	"net/http"
)

// Transport is an interface for executing HTTP requests.
// It abstracts the HTTP client implementation to allow for custom transports,
// middleware injection, and easier testing with mock implementations.
type Transport interface {
	Do(*http.Request) (*http.Response, error)
}

// httpTransport is a thin wrapper around http.Client.
// This is a good place to inject middleware later for
// logging, metrics, tracing, etc. if desired.
type httpTransport struct {
	client *http.Client
}

// NewHTTPTransport creates a new Transport implementation that wraps the provided HTTP client.
// The returned transport can be used to execute HTTP requests with the given client configuration.
func NewHTTPTransport(c *http.Client) Transport {
	return &httpTransport{client: c}
}

// Do executes the HTTP request using the underlying HTTP client.
func (t *httpTransport) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}

// mockTransport is a mock implementation of Transport for testing purposes.
// It captures the last request made and returns a predefined response and error.
type mockTransport struct {
	LastRequest *http.Request
	Response    *http.Response
	Err         error
}

// Do executes a mock HTTP request, storing the request for inspection and returning
// the predefined response and error. It respects context cancellation.
func (m *mockTransport) Do(req *http.Request) (*http.Response, error) {
	m.LastRequest = req

	// Transport implementations should respect req.Context() cancellation.
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	default:
	}

	return m.Response, m.Err
}

// newMockTransport creates a new mock transport with the specified response status,
// body, and error. This is useful for testing HTTP request handling without making
// actual network calls.
func newMockTransport(respStatus int, respBody string, err error) *mockTransport {
	return &mockTransport{
		Response: &http.Response{
			StatusCode: respStatus,
			Body:       io.NopCloser(bytes.NewBufferString(respBody)),
			Header:     make(http.Header),
		},
		Err: err,
	}
}
