package request

import "net/http"

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
	if c == nil {
		return nil
	}
	return &httpTransport{client: c}
}

// Do executes the HTTP request using the underlying HTTP client.
func (t *httpTransport) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}
