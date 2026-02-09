package request

import (
	"bytes"
	"io"
	"net/http"
)

// mockTransport is a mock implementation of http.RoundTripper for testing purposes.
// It captures the last request made and returns a predefined response and error.
type mockTransport struct {
	LastRequest *http.Request
	Response    *http.Response
	Err         error
}

// RoundTrip executes a mock HTTP request, storing the request for inspection and returning
// the predefined response and error. It respects context cancellation.
func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.LastRequest = req

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	default:
	}

	return m.Response, m.Err
}

func newMockHTTPClient(mt *mockTransport) http.Client {
	return http.Client{Transport: mt}
}

func withMockTransport(opts RequestOptions, mt *mockTransport) RequestOptions {
	return opts.WithHTTPClientFactory(func() http.Client {
		return newMockHTTPClient(mt)
	})
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
