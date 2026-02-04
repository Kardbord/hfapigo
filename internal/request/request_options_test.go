package request

import (
	"context"
	"net/http"
	"testing"

	"github.com/Kardbord/hfapigo/v4/internal/version"
)

func TestRequestOptions_With(t *testing.T) {
	tests := []struct {
		name     string
		initial  RequestOptions
		options  []RequestOption
		validate func(t *testing.T, orig, updated RequestOptions)
	}{
		{
			name:    "is immutable",
			initial: NewRequestOptions(),
			options: []RequestOption{
				WithToken("secret"),
			},
			validate: func(t *testing.T, orig, updated RequestOptions) {
				if orig.Token != "" {
					t.Errorf("expected original Token to be empty, got %q", orig.Token)
				}
				if updated.Token != "secret" {
					t.Errorf("expected updated Token to be 'secret', got %q", updated.Token)
				}
			},
		},
		{
			name:    "duplicate options - last wins",
			initial: NewRequestOptions(),
			options: []RequestOption{
				WithToken("first"),
				WithToken("second"),
			},
			validate: func(t *testing.T, orig, updated RequestOptions) {
				if updated.Token != "second" {
					t.Errorf("expected last option to win, got %q", updated.Token)
				}
			},
		},
		{
			name:    "multiple fields",
			initial: NewRequestOptions(),
			options: []RequestOption{
				WithToken("token123"),
				WithModel("llama-3"),
				WithProvider("aws"),
				WithUserAgent("myapp/1.2.3"),
			},
			validate: func(t *testing.T, orig, updated RequestOptions) {
				if updated.Token != "token123" {
					t.Errorf("expected Token 'token123', got %q", updated.Token)
				}
				if updated.Model != "llama-3" {
					t.Errorf("expected Model 'llama-3', got %q", updated.Model)
				}
				if updated.Provider != "aws" {
					t.Errorf("expected Provider 'aws', got %q", updated.Provider)
				}
				if updated.UserAgent != "myapp/1.2.3" {
					t.Errorf("expected UserAgent 'myapp/1.2.3, got %q", updated.UserAgent)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.initial
			updated := orig.With(tt.options...)

			if tt.validate != nil {
				tt.validate(t, orig, updated)
			}
		})
	}
}

func TestRequestOptions_WithHelpers(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "ok")
	mt := newMockTransport(http.StatusOK, `{}`, nil)

	opts := NewRequestOptions().
		WithBaseURL("https://example.com").
		WithToken("token").
		WithModel("model").
		WithProvider("provider").
		WithContext(ctx).
		WithMaxResponseBodyBytes(42).
		WithTransport(mt)

	if opts.BaseURL != "https://example.com" {
		t.Errorf("expected BaseURL to be set, got %q", opts.BaseURL)
	}
	if opts.Token != "token" {
		t.Errorf("expected Token to be set, got %q", opts.Token)
	}
	if opts.Model != "model" {
		t.Errorf("expected Model to be set, got %q", opts.Model)
	}
	if opts.Provider != "provider" {
		t.Errorf("expected Provider to be set, got %q", opts.Provider)
	}
	if opts.Ctx != ctx {
		t.Error("expected context to be set")
	}
	if opts.MaxResponseBodyBytes != 42 {
		t.Errorf("expected MaxResponseBodyBytes to be 42, got %d", opts.MaxResponseBodyBytes)
	}
	if opts.Transport != mt {
		t.Error("expected Transport to be set")
	}
}

func TestWithUserAgentSuffix(t *testing.T) {
	opts := NewRequestOptions().WithUserAgentSuffix("custom/1.0")
	want := version.UserAgent() + " custom/1.0"
	if opts.UserAgent != want {
		t.Errorf("expected UserAgent %q, got %q", want, opts.UserAgent)
	}
}

func TestWithHeaders_CopiesMap(t *testing.T) {
	headers := http.Header{"X-Test": []string{"one"}}
	opts := NewRequestOptions().WithHeaders(headers)
	headers.Set("X-Test", "two")

	if opts.Headers == nil || opts.Headers.Get("X-Test") != "one" {
		t.Errorf("expected headers to be copied, got %#v", opts.Headers)
	}
}

func TestRequestOptions_DefensiveHeaderClone(t *testing.T) {
	mt := newMockTransport(http.StatusOK, `{}`, nil)

	tests := []struct {
		name  string
		apply func(RequestOptions) RequestOptions
	}{
		{
			name: "With",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.With(WithToken("token"))
			},
		},
		{
			name: "WithBaseURL",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithBaseURL("https://example.com")
			},
		},
		{
			name: "WithToken",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithToken("token")
			},
		},
		{
			name: "WithModel",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithModel("model")
			},
		},
		{
			name: "WithProvider",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithProvider("provider")
			},
		},
		{
			name: "WithUserAgent",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithUserAgent("ua/1.0")
			},
		},
		{
			name: "WithUserAgentSuffix",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithUserAgentSuffix("custom/1.0")
			},
		},
		{
			name: "WithMaxResponseBodyBytes",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithMaxResponseBodyBytes(42)
			},
		},
		{
			name: "WithContext",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithContext(context.Background())
			},
		},
		{
			name: "WithHTTPClient",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithHTTPClient(http.DefaultClient)
			},
		},
		{
			name: "WithTransport",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithTransport(mt)
			},
		},
		{
			name: "WithHeaders",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithHeaders(http.Header{"X-Other": []string{"value"}})
			},
		},
		{
			name: "WithHeader",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithHeader("X-Other", "value")
			},
		},
		{
			name: "WithDefaultHeader",
			apply: func(opts RequestOptions) RequestOptions {
				return opts.WithDefaultHeader("X-Test", "default")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := NewRequestOptions().WithHeader("X-Test", "one")
			derived := tt.apply(orig)

			orig.Headers.Set("X-Test", "two")
			if got := derived.Headers.Get("X-Test"); got != "one" {
				t.Errorf("expected derived header to stay 'one', got %q", got)
			}

			derived.Headers.Set("X-Test", "three")
			if got := orig.Headers.Get("X-Test"); got != "two" {
				t.Errorf("expected original header to stay 'two', got %q", got)
			}
		})
	}
}

func TestNewRequestOptions(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, opts RequestOptions)
	}{
		{
			name: "has default context",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Ctx == nil {
					t.Error("expected default context, got nil")
				}
			},
		},
		{
			name: "has default BaseURL",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.BaseURL != DefaultBaseURL {
					t.Errorf("expected BaseURL %q, got %q", DefaultBaseURL, opts.BaseURL)
				}
			},
		},
		{
			name: "has default Token",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Token != DefaultToken {
					t.Errorf("expected Token %q, got %q", DefaultToken, opts.Token)
				}
			},
		},
		{
			name: "has default Model",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Model != DefaultModel {
					t.Errorf("expected Model %q, got %q", DefaultModel, opts.Model)
				}
			},
		},
		{
			name: "has default Provider",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Provider != DefaultProvider {
					t.Errorf("expected Provider %q, got %q", DefaultProvider, opts.Provider)
				}
			},
		},
		{
			name: "has default UserAgent",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.UserAgent != version.UserAgent() {
					t.Errorf("expected UserAgent %q, got %q", version.UserAgent(), opts.UserAgent)
				}
			},
		},
		{
			name: "has default headers",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Headers != nil {
					t.Errorf("expected default headers to be nil, got %#v", opts.Headers)
				}
			},
		},
		{
			name: "has default Transport",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.Transport == nil {
					t.Error("expected default transport, got nil")
				}
			},
		},
		{
			name: "has default max response body size",
			validate: func(t *testing.T, opts RequestOptions) {
				if opts.MaxResponseBodyBytes != DefaultMaxResponseBodyBytes {
					t.Errorf("expected %d max response body bytes, got %d", DefaultMaxResponseBodyBytes, opts.MaxResponseBodyBytes)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewRequestOptions()
			if tt.validate != nil {
				tt.validate(t, opts)
			}
		})
	}
}
