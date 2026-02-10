//go:build test
// +build test

package testutils

import (
	"net/url"
	"testing"
)

// AssertURL compares a raw URL string to the expected URL parts.
func AssertURL(t *testing.T, raw string, want *url.URL) {
	t.Helper()

	got, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("failed to parse URL %q: %v", raw, err)
	}

	if got.Scheme != want.Scheme {
		t.Errorf("unexpected scheme: %s", got.Scheme)
	}
	if got.Host != want.Host {
		t.Errorf("unexpected host: %s", got.Host)
	}
	if got.Path != want.Path {
		t.Errorf("unexpected path: %s", got.Path)
	}
	if got.RawQuery != want.RawQuery {
		t.Errorf("unexpected query: %s", got.RawQuery)
	}
	if got.Fragment != want.Fragment {
		t.Errorf("unexpected fragment: %s", got.Fragment)
	}
}
