package request

import (
	"bytes"
	"io"
	"net/http"
)

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
