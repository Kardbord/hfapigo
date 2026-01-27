package request

import (
	"testing"
)

func TestRequestOptions_With_IsImmutable(t *testing.T) {
	orig := NewRequestOptions()

	updated := orig.With(func(o *RequestOptions) {
		o.Token = "secret"
	})

	if orig.Token != "" {
		t.Fatalf("expected original Token to be empty, got %q", orig.Token)
	}

	if updated.Token != "secret" {
		t.Fatalf("expected updated Token to be set, got %q", updated.Token)
	}
}

func TestNewRequestOptions_Defaults(t *testing.T) {
	opts := NewRequestOptions()

	if opts.Ctx == nil {
		t.Fatal("expected default context")
	}
	if opts.BaseURL != DefaultBaseURL {
		t.Fatalf("unexpected BaseURL: %q", opts.BaseURL)
	}
	if opts.Token != DefaultToken {
		t.Fatalf("unexpected Token: %q", opts.Token)
	}
	if opts.Model != DefaultModel {
		t.Fatalf("unexpected Model: %q", opts.Model)
	}
	if opts.Provider != DefaultProvider {
		t.Fatalf("unexpected Provider: %q", opts.Provider)
	}
	if opts.Transport == nil {
		t.Fatal("expected default transport")
	}
}

func TestRequestOptions_With_DuplicateOptions(t *testing.T) {
	opts := NewRequestOptions().With(
		func(o *RequestOptions) { o.Token = "first" },
		func(o *RequestOptions) { o.Token = "second" },
	)

	if opts.Token != "second" {
		t.Fatalf("expected last option to win, got %q", opts.Token)
	}
}
