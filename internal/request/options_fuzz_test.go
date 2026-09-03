//go:build !integration

package request

import (
	"testing"
)

func FuzzOptionsValidate(f *testing.F) {
	f.Add("https://example.com")
	f.Fuzz(func(_ *testing.T, baseURL string) {
		opts := NewOptions().WithBaseURL(baseURL)
		_ = opts.Validate()
	})
}
