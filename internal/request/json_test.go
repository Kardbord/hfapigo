package request

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
)

func TestDoJSON(t *testing.T) {
	type req struct {
		Inputs string `json:"inputs"`
	}
	type resp struct {
		GeneratedText string `json:"generated_text"`
	}

	tests := []struct {
		name           string
		setupTransport func() Transport
		method         string
		path           string
		reqBody        req
		wantErr        bool
		wantResp       *resp
		validateErr    func(t *testing.T, err error)
		validateReq    func(t *testing.T, mt *mockTransport)
	}{
		{
			name: "successful request",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusOK, `{"generated_text":"hello"}`, nil)
			},
			method:  http.MethodPost,
			path:    "/chat",
			reqBody: req{Inputs: "hi"},
			wantErr: false,
			wantResp: &resp{
				GeneratedText: "hello",
			},
		},
		{
			name: "401 error status",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusUnauthorized, `unauthorized`, nil)
			},
			method:  http.MethodGet,
			path:    "/fail",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var apiErr *internalErrors.APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("expected *errors.APIError, got %T", err)
					return
				}
				if apiErr.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status code 401, got %d", apiErr.StatusCode)
				}
				if !apiErr.IsAuthenticationError() {
					t.Error("expected IsAuthenticationError() to return true")
				}
			},
		},
		{
			name: "500 error status",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusInternalServerError, `internal server error`, nil)
			},
			method:  http.MethodGet,
			path:    "/fail",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var apiErr *internalErrors.APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("expected *errors.APIError, got %T", err)
					return
				}
				if apiErr.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code 500, got %d", apiErr.StatusCode)
				}
				if !apiErr.IsServerError() {
					t.Error("expected IsServerError() to return true")
				}
			},
		},
		{
			name: "transport error",
			setupTransport: func() Transport {
				return &mockTransport{
					Err: errors.New("network down"),
				}
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "network down") {
					t.Errorf("expected error to contain 'network down', got: %v", err)
				}
			},
		},
		{
			name: "invalid JSON response",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusOK, `{not valid json}`, nil)
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "failed to decode response body") {
					t.Errorf("expected decode error, got: %v", err)
				}
			},
		},
		{
			name: "empty response body",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusOK, ``, nil)
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "failed to decode response body") {
					t.Errorf("expected decode error on empty body, got: %v", err)
				}
			},
		},
		{
			name: "sets Content-Type header",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusOK, `{}`, nil)
			},
			method:  http.MethodPost,
			path:    "/test",
			reqBody: req{},
			wantErr: false,
			validateReq: func(t *testing.T, mt *mockTransport) {
				if got := mt.LastRequest.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("expected Content-Type 'application/json', got %q", got)
				}
			},
		},
		{
			name: "returns zero value on error",
			setupTransport: func() Transport {
				return newMockTransport(http.StatusInternalServerError, `boom`, nil)
			},
			method:   http.MethodGet,
			path:     "/test",
			reqBody:  req{},
			wantErr:  true,
			wantResp: &resp{}, // zero value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := tt.setupTransport()
			opts := NewRequestOptions().With(func(o *RequestOptions) {
				o.Transport = transport
			})

			out, err := DoJSON[req, resp](
				opts,
				tt.method,
				tt.path,
				tt.reqBody,
			)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Fatalf("DoJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Validate error if custom validation provided
			if err != nil && tt.validateErr != nil {
				tt.validateErr(t, err)
			}

			// Validate response if expected
			if tt.wantResp != nil && !tt.wantErr {
				if out != *tt.wantResp {
					t.Errorf("DoJSON() response = %+v, want %+v", out, *tt.wantResp)
				}
			}

			// Validate request if custom validation provided
			if tt.validateReq != nil {
				if mt, ok := transport.(*mockTransport); ok {
					tt.validateReq(t, mt)
				}
			}
		})
	}
}

func TestDoJSON_MarshalError(t *testing.T) {
	// This test is separate because it uses an unmarshalable type
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

	if !strings.Contains(err.Error(), "failed to marshal request body") {
		t.Errorf("expected marshal error message, got: %v", err)
	}
}
