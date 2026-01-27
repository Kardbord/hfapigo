package request

import (
	"errors"
	"net/http"
	"testing"
)

func TestDoJSON_Success(t *testing.T) {
	mt := newMockTransport(200, `{"generated_text":"hello"}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	type req struct {
		Inputs string `json:"inputs"`
	}
	type resp struct {
		GeneratedText string `json:"generated_text"`
	}

	out, err := DoJSON[req, resp](
		opts,
		http.MethodPost,
		"/chat",
		req{Inputs: "hi"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.GeneratedText != "hello" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestDoJSON_ErrorStatus(t *testing.T) {
	mt := newMockTransport(401, `unauthorized`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	_, err := DoJSON[struct{}, struct{}](
		opts,
		http.MethodGet,
		"/fail",
		struct{}{},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDoJSON_MarshalError(t *testing.T) {
	opts := NewRequestOptions()

	// Channels cannot be marshaled to JSON
	type badReq struct {
		C chan int `json:"c"`
	}

	_, err := DoJSON[badReq, struct{}](
		opts,
		http.MethodPost,
		"/test",
		badReq{C: make(chan int)},
	)

	if err == nil {
		t.Fatal("expected marshal error, got nil")
	}
}

func TestDoJSON_TransportError(t *testing.T) {
	mt := &mockTransport{
		Err: errors.New("network down"),
	}

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	_, err := DoJSON[struct{}, struct{}](
		opts,
		http.MethodGet,
		"/test",
		struct{}{},
	)

	if err == nil {
		t.Fatal("expected transport error, got nil")
	}
}

func TestDoJSON_InvalidJSONResponse(t *testing.T) {
	mt := newMockTransport(200, `{not valid json}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	_, err := DoJSON[struct{}, struct{}](
		opts,
		http.MethodGet,
		"/test",
		struct{}{},
	)

	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}

func TestDoJSON_EmptyBody(t *testing.T) {
	mt := newMockTransport(200, ``, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	_, err := DoJSON[struct{}, struct{}](
		opts,
		http.MethodGet,
		"/test",
		struct{}{},
	)

	if err == nil {
		t.Fatal("expected decode error on empty body")
	}
}

func TestDoJSON_SetsContentType(t *testing.T) {
	mt := newMockTransport(200, `{}`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = mt
	})

	_, err := DoJSON[struct{}, struct{}](
		opts,
		http.MethodPost,
		"/test",
		struct{}{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := mt.LastRequest.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected Content-Type: %q", got)
	}
}

func TestDoJSON_ReturnsZeroValueOnError(t *testing.T) {
	ft := newMockTransport(500, `boom`, nil)

	opts := NewRequestOptions().With(func(o *RequestOptions) {
		o.Transport = ft
	})

	type resp struct {
		Value string
	}

	out, err := DoJSON[struct{}, resp](
		opts,
		http.MethodGet,
		"/test",
		struct{}{},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if out != (resp{}) {
		t.Fatalf("expected zero value response, got %+v", out)
	}
}
