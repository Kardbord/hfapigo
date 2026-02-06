package hfapigo

import (
	"context"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/api"
	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// WithBaseURL returns a RequestOption that sets the base URL for API requests.
// The base URL is the root endpoint for all HuggingFace API calls and must not
// include query parameters or fragments.
func WithBaseURL(u string) api.RequestOption {
	return request.WithBaseURL(u)
}

// WithToken returns a RequestOption that sets the authentication token for API requests.
// The token is used for Bearer authentication with the HuggingFace API.
func WithToken(t string) api.RequestOption {
	return request.WithToken(t)
}

// WithModel returns a RequestOption that sets the model to use for API requests.
// The model specifies which HuggingFace model should process the request.
func WithModel(m string) api.RequestOption {
	return request.WithModel(m)
}

// WithProvider returns a RequestOption that sets the provider for API requests.
// The provider specifies which inference provider should handle the request.
func WithProvider(p string) api.RequestOption {
	return request.WithProvider(p)
}

// WithHTTPClient returns a RequestOption that sets a custom HTTP client for API requests.
// This allows customization of transport settings, timeouts, and other HTTP client configurations.
func WithHTTPClient(c *http.Client) api.RequestOption {
	return request.WithHTTPClient(c)
}

// WithContext returns a RequestOption that sets the context for API requests.
// The context can be used for cancellation, timeouts, and passing request-scoped values.
// If a nil context is provided, the SDK will fall back to context.Background().
func WithContext(ctx context.Context) api.RequestOption {
	return request.WithContext(ctx)
}

// WithUserAgentSuffix returns a RequestOption that appends a suffix to the SDK user agent string.
func WithUserAgentSuffix(s string) api.RequestOption {
	return request.WithUserAgentSuffix(s)
}

// WithMaxResponseBodyBytes returns a RequestOption that sets the maximum number of bytes
// read from any response body. Values <= 0 fall back to the default.
func WithMaxResponseBodyBytes(n int64) api.RequestOption {
	return request.WithMaxResponseBodyBytes(n)
}

// WithHeaders returns a RequestOption that sets custom headers applied to every request,
// overriding any existing values for matching keys.
// Per-request headers can still override these values when provided.
func WithHeaders(h http.Header) api.RequestOption {
	return request.WithHeaders(h)
}

// WithHeader returns a RequestOption that sets a single header applied to every request.
func WithHeader(key, value string) api.RequestOption {
	return request.WithHeader(key, value)
}

// WithDefaultHeader returns a RequestOption that sets a header only if missing or empty.
func WithDefaultHeader(key, value string) api.RequestOption {
	return request.WithDefaultHeader(key, value)
}
