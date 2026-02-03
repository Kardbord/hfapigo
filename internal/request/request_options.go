package request

import (
	"context"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/version"
)

// RequestOptions holds configuration settings for API requests.
type RequestOptions struct {
	Ctx                  context.Context
	BaseURL              string
	Token                string
	Model                string
	Provider             string
	UserAgent            string
	MaxResponseBodyBytes int64
	Transport            Transport
}

const (
	// DefaultBaseURL is the default HuggingFace API endpoint.
	DefaultBaseURL = "https://router.huggingface.co"
	// DefaultToken is the default authentication token (empty string).
	DefaultToken = ""
	// DefaultModel is the default model to use
	DefaultModel = ""
	// DefaultProvider is the default inference provider
	DefaultProvider = ""
	// DefaultMaxResponseBodyBytes caps the amount of response data read into memory by default.
	DefaultMaxResponseBodyBytes int64 = 1 << 20 // 1 MiB
)

// NewRequestOptions creates a new RequestOptions instance with default values.
// The returned options use a background context, default endpoints, and the default HTTP client.
func NewRequestOptions() RequestOptions {
	return RequestOptions{
		Ctx:                  context.Background(),
		BaseURL:              DefaultBaseURL,
		Token:                DefaultToken,
		Model:                DefaultModel,
		Provider:             DefaultProvider,
		UserAgent:            version.UserAgent(),
		MaxResponseBodyBytes: DefaultMaxResponseBodyBytes,
		Transport:            NewHTTPTransport(http.DefaultClient),
	}
}

// RequestOption is a function type that modifies RequestOptions.
// It follows the functional options pattern for flexible configuration.
type RequestOption func(*RequestOptions)

// apply applies a series of RequestOption functions to the RequestOptions instance.
func (o *RequestOptions) apply(opts ...RequestOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// With returns a new RequestOptions instance with the provided options applied.
// This method creates a copy of the current options and applies modifications to it.
func (o RequestOptions) With(opts ...RequestOption) RequestOptions {
	o.apply(opts...)
	return o
}
