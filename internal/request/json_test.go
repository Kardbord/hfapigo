package request

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/testutils"
)

func TestDoJSON(t *testing.T) {
	type req struct {
		Inputs string `json:"inputs"`
	}
	type resp struct {
		GeneratedText string `json:"generated_text"`
	}

	type testCase struct {
		name              string
		setupOpts         func() RequestOptions
		method            string
		path              string
		reqBody           req
		wantErr           bool
		wantResp          *resp
		validateErr       func(t *testing.T, err error)
		validateReq       func(t *testing.T, req *http.Request)
		validateTransport func(t *testing.T, mt *testutils.MockTransport)
	}

	tests := []testCase{
		{
			name: "successful request",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{"generated_text":"hello"}`, nil))
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
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewMockTransport(http.StatusUnauthorized, `unauthorized`, nil))
				})
			},
			method:  http.MethodGet,
			path:    "/fail",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				apiErr := testutils.AssertAPIErrorStatus(t, err, http.StatusUnauthorized)
				if !apiErr.IsAuthenticationError() {
					t.Error("expected IsAuthenticationError() to return true")
				}
			},
		},
		{
			name: "500 error status",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewMockTransport(http.StatusInternalServerError, `internal server error`, nil))
				})
			},
			method:  http.MethodGet,
			path:    "/fail",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				apiErr := testutils.AssertAPIErrorStatus(t, err, http.StatusInternalServerError)
				if !apiErr.IsServerError() {
					t.Error("expected IsServerError() to return true")
				}
			},
		},
		{
			name: "transport error",
			setupOpts: func() RequestOptions {
				mt := &testutils.MockTransport{
					Err: errors.New("network down"),
				}
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.RequireError(t, err)
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindTransport)
			},
		},
		{
			name: "invalid JSON response",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{not valid json}`, nil))
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.RequireError(t, err)
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindSerialization)
			},
		},
		{
			name: "empty response body",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, ``, nil))
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindSerialization)
			},
		},
		{
			name: "response body too large",
			setupOpts: func() RequestOptions {
				large := `{"generated_text":"` + strings.Repeat("a", int(DefaultMaxResponseBodyBytes)) + `"}`
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, large, nil))
				})
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindInternal)
			},
		},
		{
			name: "custom response limit allows larger body",
			setupOpts: func() RequestOptions {
				large := `{"generated_text":"` + strings.Repeat("a", int(DefaultMaxResponseBodyBytes)) + `"}`
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, large, nil))
				}).WithMaxResponseBodyBytes(DefaultMaxResponseBodyBytes + 64)
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
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{}`, nil))
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
			name: "rejects non-JSON Content-Type override",
			setupOpts: func() RequestOptions {
				opts := NewRequestOptions().WithHeader("Content-Type", "text/plain")
				return opts.WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{}`, nil))
				})
			},
			method:  http.MethodPost,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindSerialization)
			},
		},
		{
			name: "fills empty Content-Type override",
			setupOpts: func() RequestOptions {
				opts := NewRequestOptions().WithHeader("Content-Type", "")
				return opts.WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{}`, nil))
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
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewMockTransport(http.StatusInternalServerError, `boom`, nil))
				})
			},
			method:   http.MethodGet,
			path:     "/test",
			reqBody:  req{},
			wantErr:  true,
			wantResp: &resp{}, // zero value
		},
		{
			name: "sets Accept header",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewJSONMockTransport(http.StatusOK, `{}`, nil))
				})
			},
			method:  http.MethodPost,
			path:    "/test",
			reqBody: req{},
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				if got := req.Header.Get("Accept"); got != "application/json" {
					t.Errorf("expected Accept 'application/json', got %q", got)
				}
			},
		},
		{
			name: "allows missing response Content-Type",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:   http.MethodGet,
			path:     "/test",
			reqBody:  req{},
			wantErr:  false,
			wantResp: &resp{},
		},
		{
			name: "errors on non-JSON response Content-Type",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				mt.Response.Header.Set("Content-Type", "text/plain")
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindSerialization)
			},
		},
		{
			name: "errors on invalid response Content-Type syntax",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				mt.Response.Header.Set("Content-Type", "application/json; charset")
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
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
			name: "accepts +json response Content-Type",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{"generated_text":"hello"}`, nil)
				mt.Response.Header.Set("Content-Type", "application/problem+json")
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: false,
			wantResp: &resp{
				GeneratedText: "hello",
			},
		},
		{
			name: "returns zero value on 204 No Content",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewMockTransport(http.StatusNoContent, ``, nil))
				})
			},
			method:   http.MethodGet,
			path:     "/test",
			reqBody:  req{},
			wantErr:  false,
			wantResp: &resp{},
		},
		{
			name: "returns zero value on 205 Reset Content",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client {
					return testutils.NewMockHTTPClient(testutils.NewMockTransport(http.StatusResetContent, ``, nil))
				})
			},
			method:   http.MethodGet,
			path:     "/test",
			reqBody:  req{},
			wantErr:  false,
			wantResp: &resp{},
		},
		{
			name: "drains response on size error",
			setupOpts: func() RequestOptions {
				data := strings.Repeat("a", 16)
				tracker := &testutils.ReadTracker{Data: []byte(data)}
				mt := &testutils.MockTransport{
					Response: &http.Response{
						StatusCode: http.StatusOK,
						Body:       tracker,
						Header:     make(http.Header),
					},
				}
				mt.Response.Header.Set("Content-Type", "application/json")
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithMaxResponseBodyBytes(4)
			},
			method:  http.MethodGet,
			path:    "/test",
			reqBody: req{},
			wantErr: true,
			validateTransport: func(t *testing.T, mt *testutils.MockTransport) {
				tracker, ok := mt.Response.Body.(*testutils.ReadTracker)
				if !ok {
					t.Fatal("expected readTracker body")
				}
				if tracker.ReadBytes != len(tracker.Data) {
					t.Fatalf("expected body to be drained, read %d bytes, want %d", tracker.ReadBytes, len(tracker.Data))
				}
				if !tracker.Closed {
					t.Fatal("expected response body to be closed")
				}
			},
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
				if opts.HTTPClient == nil {
					t.Fatal("expected http client")
				}
				if mt, ok := opts.HTTPClient.Transport.(*testutils.MockTransport); ok && mt.LastRequest != nil {
					tt.validateReq(t, mt.LastRequest)
				}
			}

			if tt.validateTransport != nil {
				if opts.HTTPClient == nil {
					t.Fatal("expected http client")
				}
				if mt, ok := opts.HTTPClient.Transport.(*testutils.MockTransport); ok {
					tt.validateTransport(t, mt)
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

type validationReq struct{}

func (validationReq) MarshalJSON() ([]byte, error) {
	return nil, &internalErrors.SDKError{
		Kind:    internalErrors.SDKErrorKindValidation,
		Message: "invalid payload",
	}
}

func TestDoJSON_MarshalValidationError(t *testing.T) {
	opts := NewRequestOptions()

	_, err := DoJSON[validationReq, struct{}](
		opts,
		http.MethodPost,
		"/test",
		validationReq{},
	)

	if err == nil {
		t.Fatal("expected marshal error, got nil")
	}
	testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindValidation)
}
