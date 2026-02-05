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
	"github.com/Kardbord/hfapigo/v4/internal/version"
)

type errorReadCloser struct{}

func (e errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (e errorReadCloser) Close() error { return nil }

func assertURL(t *testing.T, raw string, want *url.URL) {
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
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithBaseURL("https://example.com").
					WithToken("abc123").
					WithTransport(mt).
					WithHeaders(http.Header{"X-Test": []string{"yes"}})
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				assertURL(t, req.URL.String(), &url.URL{
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
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithBaseURL("https://example.com/api").
					WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "v1/chat/completions",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				assertURL(t, req.URL.String(), &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/api/v1/chat/completions",
				})
			},
		},
		{
			name: "preserves query string and fragment",
			setupOpts: func() RequestOptions {
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithBaseURL("https://example.com/api").
					WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "/v1/chat/completions?model=foo#section",
			body:    nil,
			wantErr: false,
			validateReq: func(t *testing.T, req *http.Request) {
				assertURL(t, req.URL.String(), &url.URL{
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
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithContext(ctx).
					WithTransport(mt)
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
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithContext(nil).
					WithTransport(mt)
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
				mt := newMockTransport(http.StatusUnauthorized, `unauthorized`, nil)
				return NewRequestOptions().WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var apiErr *internalErrors.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got %T", err)
				}
				if apiErr.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", apiErr.StatusCode)
				}
				if apiErr.Message == "" {
					t.Errorf("expected non-empty API error message")
				}
			},
		},
		{
			name: "header override",
			setupOpts: func() RequestOptions {
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithToken("default").
					WithTransport(mt).
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
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().
					WithBaseURL("http://[::1").
					WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindConfiguration {
					t.Errorf("expected configuration SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "returns internal SDKError on invalid method",
			setupOpts: func() RequestOptions {
				mt := newMockTransport(http.StatusOK, `{}`, nil)
				return NewRequestOptions().WithTransport(mt)
			},
			method:  "GET\n",
			path:    "/test",
			body:    nil,
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
			name: "returns internal SDKError when reading error response fails",
			setupOpts: func() RequestOptions {
				mt := &mockTransport{
					Response: &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       errorReadCloser{},
						Header:     make(http.Header),
					},
				}
				return NewRequestOptions().WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
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
			name: "returns configuration SDKError when transport is nil",
			setupOpts: func() RequestOptions {
				return NewRequestOptions().WithTransport(nil)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var sdkErr *internalErrors.SDKError
				if !errors.As(err, &sdkErr) {
					t.Fatalf("expected SDKError, got %T", err)
				}
				if sdkErr.Kind != internalErrors.SDKErrorKindConfiguration {
					t.Errorf("expected configuration SDKError, got %q", sdkErr.Kind)
				}
			},
		},
		{
			name: "returns API error with truncated body on oversized error response",
			setupOpts: func() RequestOptions {
				mt := newMockTransport(http.StatusTooManyRequests, strings.Repeat("x", 10), nil)
				return NewRequestOptions().
					WithMaxResponseBodyBytes(5).
					WithTransport(mt)
			},
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				var apiErr *internalErrors.APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got %T", err)
				}
				if apiErr.StatusCode != http.StatusTooManyRequests {
					t.Errorf("expected status 429, got %d", apiErr.StatusCode)
				}
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
				if mt, ok := opts.Transport.(*mockTransport); ok && mt.LastRequest != nil {
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
			mt := newMockTransport(http.StatusOK, `{}`, nil)
			opts := NewRequestOptions().WithTransport(mt)

			_, err := DoBytes(opts, http.MethodPost, "/test", tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateReq != nil && mt.LastRequest != nil {
				tt.validateReq(t, mt.LastRequest)
			}
		})
	}
}

func TestDoRaw(t *testing.T) {
	t.Run("returns response on non-2xx without closing body", func(t *testing.T) {
		tracker := &closeTracker{}
		mt := &mockTransport{
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       tracker,
				Header:     make(http.Header),
			},
		}
		opts := NewRequestOptions().WithTransport(mt)

		resp, err := DoRaw(opts, http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status 401 response, got %#v", resp)
		}
		if tracker.closed {
			t.Fatal("expected response body to remain open")
		}
	})
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		path    string
		want    *url.URL
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := joinURL(tt.baseURL, tt.path)
			if err != nil {
				t.Fatalf("joinURL error: %v", err)
			}
			assertURL(t, got, tt.want)
		})
	}
}

type closeTracker struct {
	closed bool
}

func (c *closeTracker) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (c *closeTracker) Close() error {
	c.closed = true
	return nil
}

func TestDo_ClosesResponseOnTransportError(t *testing.T) {
	tracker := &closeTracker{}
	mt := &mockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       tracker,
			Header:     make(http.Header),
		},
		Err: errors.New("boom"),
	}
	opts := NewRequestOptions().WithTransport(mt)

	_, err := Do(opts, http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !tracker.closed {
		t.Fatal("expected response body to be closed")
	}
}
