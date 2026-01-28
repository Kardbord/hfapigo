package request

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestHTTPTransport_Delegates(t *testing.T) {
	called := false
	client := &http.Client{
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{}")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	tp := NewHTTPTransport(client)

	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_, err := tp.Do(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected underlying client to be called")
	}
}
