package request

import (
	"bytes"
	"io"
	"net/http"
)

type Transport interface {
	Do(*http.Request) (*http.Response, error)
}

// Thin wrapper around http.Client for now. This
// is a good place to inject middleware later for
// logging, metrics, tracing, etc. if desired.
type httpTransport struct {
	client *http.Client
}

func NewHTTPTransport(c *http.Client) Transport {
	return &httpTransport{client: c}
}

func (t *httpTransport) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}

type mockTransport struct {
	LastRequest *http.Request
	Response    *http.Response
	Err         error
}

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
