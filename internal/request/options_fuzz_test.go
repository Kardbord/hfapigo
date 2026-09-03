//go:build !integration

package request

import (
	"testing"
)

func FuzzOptionsValidate(f *testing.F) {
	f.Add("https://example.com")
	f.Add("")
	f.Add("example.com")
	f.Add("https://example.com/api?token=secret")
	f.Add("https://example.com/api#section")
	f.Add("https://example.com/api?token=secret#section")
	f.Add("http://[::1")
	f.Fuzz(func(_ *testing.T, baseURL string) {
		opts := NewOptions().WithBaseURL(baseURL)
		_ = opts.Validate()
	})
}
