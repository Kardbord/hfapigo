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

	type testCase struct {
		name        string
		setupOpts   func() RequestOptions
		method      string
		path        string
		reqBody     req
		wantErr     bool
		wantResp    *resp
		validateErr func(t *testing.T, err error)
		validateReq func(t *testing.T, req *http.Request)
	}

	tests := []testCase{
		{
			name: "successful request",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, `{"generated_text":"hello"}`, nil)
				})
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
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusUnauthorized, `unauthorized`, nil)
				})
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
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusInternalServerError, `internal server error`, nil)
				})
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
			setupOpts: func() RequestOptions {
				mt := &mockTransport{
					Err: errors.New("network down"),
				}
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = mt
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error from transport")
				}

				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindTransport {
					t.Errorf("expected transport SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "invalid JSON response",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, `{not valid json}`, nil)
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error when decoding invalid JSON")
				}

				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindSerialization {
					t.Errorf("expected serialization SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "empty response body",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, ``, nil)
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindSerialization {
					t.Errorf("expected serialization SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "response body too large",
			setupOpts: func() RequestOptions {
				large := `{"generated_text":"` + strings.Repeat("a", int(DefaultMaxResponseBodyBytes)) + `"}`
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, large, nil)
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindInternal {
					t.Errorf("expected internal SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "custom response limit allows larger body",
			setupOpts: func() RequestOptions {
				large := `{"generated_text":"` + strings.Repeat("a", int(DefaultMaxResponseBodyBytes)) + `"}`
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, large, nil)
					o.MaxResponseBodyBytes = DefaultMaxResponseBodyBytes + 64
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: false,
			wantResp: &resp{
				GeneratedText: strings.Repeat("a", int(DefaultMaxResponseBodyBytes)),
			},
		},
		{
			name: "sets Content-Type header",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusOK, `{}`, nil)
				})
			},
			method:  http.MethodPost,
			path:    "/test",
			reqBody: req{},
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				if got := req.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("expected Content-Type 'application/json', got %q", got)
				}
			},
		},
		{
			name: "returns zero value on error",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().With(func(o *RequestOptions) {
					o.Transport = newMockTransport(http.StatusInternalServerError, `boom`, nil)
				})
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
			opts := tt.setupOpts()

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
			if tt.wantResp != nil {
				if out != *tt.wantResp {
					t.Errorf("DoJSON() response = %+v, want %+v", out, *tt.wantResp)
				}
			}

			// Validate request if custom validation provided
			if tt.validateReq != nil {
				if mt, ok := opts.Transport.(*mockTransport); ok && mt.LastRequest != nil {
					tt.validateReq(t, mt.LastRequest)
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
	var sdkErr *internalErrors.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}
	if sdkErr.Kind != internalErrors.SDKErrorKindSerialization {
		t.Fatalf("expected serialization SDKError, got %q", sdkErr.Kind)
	}
}
