package request

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestDo_BuildsRequestCorrectly(t *testing.T) {
	mt := newMockTransport(200, `{}`, nil)

	opts := NewRequestOptions().With(
		func(o *RequestOptions) {
			o.BaseURL = "https://example.com"
			o.Token = "abc123"
			o.Transport = mt
		},
	)

	_, err := Do(
		opts,
		http.MethodGet,
		"/test",
		nil,
		map[string]string{"X-Test": "yes"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := mt.LastRequest
	if req == nil {
		t.Fatal("expected request to be sent")
	}

	if req.URL.String() != "https://example.com/test" {
		t.Fatalf("unexpected URL: %s", req.URL)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer abc123" {
		t.Fatalf("unexpected Authorization header: %q", got)
	}

	if got := req.Header.Get("X-Test"); got != "yes" {
		t.Fatalf("unexpected X-Test header: %q", got)
	}
}

func TestDo_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mt := newMockTransport(200, `{}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Ctx = ctx
		o.Transport = mt
	})

	_, err := Do(
		opts,
		http.MethodGet,
		"/test",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestDoBytes_BodyIsCorrect(t *testing.T) {
	mt := newMockTransport(200, `{}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	data := []byte("hello world")

	_, err := DoBytes(
		opts,
		http.MethodPost,
		"/test",
		data,
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body, _ := io.ReadAll(mt.LastRequest.Body)
	if string(body) != "hello world" {
		t.Fatalf("unexpected body: %q", string(body))
	}
}

func TestDo_HeaderOverride(t *testing.T) {
	mt := newMockTransport(200, `{}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Token = "default"
		o.Transport = mt
	})

	_, err := Do(
		opts,
		http.MethodGet,
		"/test",
		nil,
		map[string]string{
			"Authorization": "Bearer override",
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := mt.LastRequest.Header.Get("Authorization"); got != "Bearer override" {
		t.Fatalf("expected override auth header, got %q", got)
	}
}
