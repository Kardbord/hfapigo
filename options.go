package hfapigo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// WithBaseURL returns a RequestOption that sets the base URL for API requests.
// The base URL is the root endpoint for all HuggingFace API calls.
func WithBaseURL(u string) request.RequestOption {
	return func(o *request.RequestOptions) { o.BaseURL = u }
}

// WithToken returns a RequestOption that sets the authentication token for API requests.
// The token is used for Bearer authentication with the HuggingFace API.
func WithToken(t string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Token = t }
}

// WithModel returns a RequestOption that sets the model to use for API requests.
// The model specifies which HuggingFace model should process the request.
func WithModel(m string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Model = m }
}

// WithProvider returns a RequestOption that sets the provider for API requests.
// The provider specifies which inference provider should handle the request.
func WithProvider(p string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Provider = p }
}

// WithHTTPClient returns a RequestOption that sets a custom HTTP client for API requests.
// This allows customization of transport settings, timeouts, and other HTTP client configurations.
func WithHTTPClient(c *http.Client) request.RequestOption {
	return func(o *request.RequestOptions) { o.Transport = request.NewHTTPTransport(c) }
}

// WithContext returns a RequestOption that sets the context for API requests.
// The context can be used for cancellation, timeouts, and passing request-scoped values.
func WithContext(ctx context.Context) request.RequestOption {
	return func(o *request.RequestOptions) { o.Ctx = ctx }
}

func WithUserAgentSuffix(s string) request.RequestOption {
	return func(o *request.RequestOptions) { o.UserAgent = fmt.Sprintf("%s %s", UserAgentPrefix(), s) }
}
