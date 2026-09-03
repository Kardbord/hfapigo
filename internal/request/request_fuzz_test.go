//go:build !integration

package request

import (
	"bytes"
	"testing"
)

func FuzzJoinURL(f *testing.F) {
	f.Add("https://example.com", "/api")
	f.Fuzz(func(_ *testing.T, baseURL, path string) {
		_, _ = joinURL(baseURL, path)
	})
}

func FuzzBuildHTTPRequest(f *testing.F) {
	f.Add("GET", "/test", []byte("hello"))
	f.Fuzz(func(_ *testing.T, method, path string, body []byte) {
		opts := NewOptions()
		_, _ = buildHTTPRequest(opts, method, path, bytes.NewReader(body))
	})
}
