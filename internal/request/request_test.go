package request

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/testutils"
	"github.com/Kardbord/hfapigo/v4/internal/version"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name        string
		setupOpts   func() RequestOptions
		method      string
		path        string
		body        io.Reader
		wantErr     bool
		validateReq func(t *testing.T, req *http.Request)
		validateErr func(t *testing.T, err error)
	}{
		{
			name: "builds request correctly",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com").
					WithToken("abc123").
					WithHeaders(http.Header{"X-Test": []string{"yes"}})
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				testutils.AssertURL(t, req.URL.String(), &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/test",
				})
				if got := req.Header.Get("Authorization"); got != "Bearer abc123" {
					t.Errorf("unexpected Authorization header: %q", got)
				}
				if got := req.Header.Get("X-Test"); got != "yes" {
					t.Errorf("unexpected X-Test header: %q", got)
				}
				if got := req.Header.Get("User-Agent"); got != version.UserAgent() {
					t.Errorf("unexpected User-Agent header: %q, want %q", got, version.UserAgent())
				}
			},
		},
		{
			name: "joins base URL path with relative path",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com/api")
			},
			method:  http.MethodGet,
			path:    "v1/chat/completions",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				testutils.AssertURL(t, req.URL.String(), &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/api/v1/chat/completions",
				})
			},
		},
		{
			name: "preserves query string and fragment",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com/api")
			},
			method:  http.MethodGet,
			path:    "/v1/chat/completions?model=foo#section",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				testutils.AssertURL(t, req.URL.String(), &url.URL{
					Scheme:   "https",
					Host:     "example.com",
					Path:     "/api/v1/chat/completions",
					RawQuery: "model=foo",
					Fragment: "section",
				})
			},
		},
		{
			name: "context canceled",
			setupOpts: func() RequestOptions {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithContext(ctx)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected context cancellation error")
				}
			},
		},
		{
			name: "nil context uses background",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithContext(nil)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Context() == nil {
					t.Fatal("expected non-nil request context")
				}
				if err := req.Context().Err(); err != nil {
					t.Fatalf("unexpected context error: %v", err)
				}
			},
		},
		{
			name: "returns API error on non-2xx response",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusUnauthorized, `unauthorized`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				apiErr := testutils.AssertAPIErrorStatus(t, err, http.StatusUnauthorized)
				if apiErr.Message == "" {
					t.Errorf("expected non-empty API error message")
				}
			},
		},
		{
			name: "header override",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithToken("default").
					WithHeaders(http.Header{"Authorization": []string{"Bearer override"}})
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				if got := req.Header.Get("Authorization"); got != "Bearer override" {
					t.Errorf("expected override auth header, got %q", got)
				}
			},
		},
		{
			name: "returns configuration SDKError on bad base URL",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("http://[::1")
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
			},
		},
		{
			name: "returns configuration SDKError on base URL with query",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com/api?token=abc")
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
			},
		},
		{
			name: "returns configuration SDKError on base URL with fragment",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com/api#section")
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
			},
		},
		{
			name: "returns configuration SDKError on base URL with query and fragment",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithBaseURL("https://example.com/api?token=abc#section")
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
			},
		},
		{
			name: "returns internal SDKError on invalid method",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  "GET\n",
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindInternal)
			},
		},
		{
			name: "returns internal SDKError when reading error response fails",
			setupOpts: func() RequestOptions {
				mt := &testutils.MockTransport{
					Response: &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       testutils.ErrorReadCloser{},
						Header:     make(http.Header),
					},
				}
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindInternal)
			},
		},
		{
			name: "returns APIError on nil error response body",
			setupOpts: func() RequestOptions {
				mt := &testutils.MockTransport{
					Response: &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       nil,
						Header:     make(http.Header),
					},
				}
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertAPIErrorStatus(t, err, http.StatusBadRequest)
			},
		},
		{
			name: "returns APIError on http.NoBody error response",
			setupOpts: func() RequestOptions {
				mt := &testutils.MockTransport{
					Response: &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       http.NoBody,
						Header:     make(http.Header),
					},
				}
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertAPIErrorStatus(t, err, http.StatusBadRequest)
			},
		},
		{
			name: "returns configuration SDKError when http client is nil",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithHTTPClientFactory(nil)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
			},
		},
		{
			name: "returns API error with truncated body on oversized error response",
			setupOpts: func() RequestOptions {
				mt := testutils.NewMockTransport(http.StatusTooManyRequests, strings.Repeat("x", 10), nil)
				return NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
					WithMaxResponseBodyBytes(5)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				apiErr := testutils.AssertAPIErrorStatus(t, err, http.StatusTooManyRequests)
				if apiErr.Message != strings.Repeat("x", 5)+" [truncated]" {
					t.Errorf("unexpected message: %q", apiErr.Message)
				}
				if len(apiErr.Body) != 5 {
					t.Errorf("expected truncated body length 5, got %d", len(apiErr.Body))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.setupOpts()

			resp, err := Do(opts, tt.method, tt.path, tt.body)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Fatalf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Validate error if custom validation provided
			if tt.validateErr != nil {
				tt.validateErr(t, err)
			}

			// Validate request if no error and validation provided
			if !tt.wantErr && tt.validateReq != nil {
				// Get the mock transport to access the last request
				if opts.HTTPClient == nil {
					t.Fatal("expected http client")
				}
				if mt, ok := opts.HTTPClient.Transport.(*testutils.MockTransport); ok && mt.LastRequest != nil {
					tt.validateReq(t, mt.LastRequest)
				} else {
					t.Fatal("expected mock transport with LastRequest")
				}
			}

			// Close response body if present
			if resp != nil {
				resp.Body.Close()
			}
		})
	}
}

func TestDo_DrainsErrorResponseBody(t *testing.T) {
	data := strings.Repeat("a", 10)
	tracker := &testutils.ReadTracker{Data: []byte(data)}
	mt := &testutils.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       tracker,
			Header:     make(http.Header),
		},
	}
	opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithMaxResponseBodyBytes(4)

	_, err := Do(opts, http.MethodGet, "/test", nil)
	testutils.RequireError(t, err)
	if tracker.ReadBytes != len(data) {
		t.Fatalf("expected body to be drained, read %d bytes, want %d", tracker.ReadBytes, len(data))
	}
	if !tracker.Closed {
		t.Fatal("expected response body to be closed")
	}
}

func TestDoBytes(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		validateReq func(t *testing.T, req *http.Request)
	}{
		{
			name: "sends body correctly",
			data: []byte("hello world"),
			validateReq: func(t *testing.T, req *http.Request) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				if string(body) != "hello world" {
					t.Errorf("unexpected body: %q, want %q", string(body), "hello world")
				}
			},
		},
		{
			name: "sends empty body",
			data: []byte(""),
			validateReq: func(t *testing.T, req *http.Request) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				if string(body) != "" {
					t.Errorf("unexpected body: %q, want empty", string(body))
				}
			},
		},
		{
			name: "sends JSON data",
			data: []byte(`{"key":"value"}`),
			validateReq: func(t *testing.T, req *http.Request) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				if !strings.Contains(string(body), "key") {
					t.Errorf("unexpected body: %q", string(body))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := testutils.NewMockTransport(http.StatusOK, `{}`, nil)
			opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })

			_, err := DoBytes(opts, http.MethodPost, "/test", tt.data)
			testutils.RequireNoError(t, err)

			if tt.validateReq != nil && mt.LastRequest != nil {
				tt.validateReq(t, mt.LastRequest)
			}
		})
	}
}

func TestDoRaw(t *testing.T) {
	t.Run("returns response on non-2xx without closing body", func(t *testing.T) {
		tracker := &testutils.CloseTracker{}
		mt := &testutils.MockTransport{
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       tracker,
				Header:     make(http.Header),
			},
		}
		opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })

		resp, err := DoRaw(opts, http.MethodGet, "/test", nil)
		testutils.RequireNoError(t, err)
		if resp == nil || resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status 401 response, got %#v", resp)
		}
		if tracker.Closed {
			t.Fatal("expected response body to remain open")
		}
	})
	t.Run("normalizes nil body to non-nil response body", func(t *testing.T) {
		mt := &testutils.MockTransport{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       nil,
				Header:     make(http.Header),
			},
		}
		opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })

		resp, err := DoRaw(opts, http.MethodGet, "/test", nil)
		testutils.RequireNoError(t, err)
		if resp.Body == nil {
			t.Fatal("expected non-nil response body")
		}
	})
	t.Run("returns error when client transport returns nil response without error", func(t *testing.T) {
		mt := &testutils.MockTransport{}
		opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })

		_, err := DoRaw(opts, http.MethodGet, "/test", nil)
		testutils.RequireError(t, err)
		testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindTransport)
	})
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		path    string
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "empty path returns base URL",
			baseURL: "https://example.com/api",
			path:    "",
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/api",
			},
		},
		{
			name:    "path only joins with base",
			baseURL: "https://example.com/api",
			path:    "v1/chat",
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/api/v1/chat",
			},
		},
		{
			name:    "query only preserves base path",
			baseURL: "https://example.com/api",
			path:    "?model=foo",
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/api",
				RawQuery: "model=foo",
			},
		},
		{
			name:    "fragment only preserves base path",
			baseURL: "https://example.com/api",
			path:    "#section",
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/api",
				Fragment: "section",
			},
		},
		{
			name:    "path with query and fragment",
			baseURL: "https://example.com/api",
			path:    "/v1/chat?model=foo#section",
			want: &url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/api/v1/chat",
				RawQuery: "model=foo",
				Fragment: "section",
			},
		},
		{
			name:    "rejects full URL path",
			baseURL: "https://example.com/api",
			path:    "https://evil.example.com/override",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := joinURL(tt.baseURL, tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected joinURL error")
				}
				return
			}
			if err != nil {
				t.Fatalf("joinURL error: %v", err)
			}
			testutils.AssertURL(t, got, tt.want)
		})
	}
}

func TestDo_IgnoresResponseOnTransportError(t *testing.T) {
	tracker := &testutils.ReadTracker{Data: []byte("ignored")}
	mt := &testutils.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       tracker,
			Header:     make(http.Header),
		},
		Err: errors.New("boom"),
	}
	opts := NewRequestOptions().WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })

	_, err := Do(opts, http.MethodGet, "/test", nil)
	testutils.RequireError(t, err)
	if tracker.ReadBytes != 0 {
		t.Fatalf("expected response body to be ignored, read %d bytes", tracker.ReadBytes)
	}
	if tracker.Closed {
		t.Fatal("expected response body to remain open")
	}
}
