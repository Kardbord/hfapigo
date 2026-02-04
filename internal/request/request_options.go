package request

import (
	"context"
	"net/http"
	"net/textproto"

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
	Headers              http.Header
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
		Headers:              nil,
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

// WithHeaders returns a new RequestOptions instance with the provided headers merged in.
func (o RequestOptions) WithHeaders(h http.Header) RequestOptions {
	o.Headers = mergeHeaders(o.Headers, h)
	return o
}

// WithHeader returns a new RequestOptions instance with a single header applied.
func (o RequestOptions) WithHeader(key, value string) RequestOptions {
	o.Headers = mergeHeaders(o.Headers, http.Header{key: []string{value}})
	return o
}

// WithDefaultHeader returns a new RequestOptions instance with a header set
// only if the header is missing or empty.
func (o RequestOptions) WithDefaultHeader(key, value string) RequestOptions {
	o.Headers = ensureHeader(o.Headers, key, value)
	return o
}

// WithHeaders returns a RequestOption that sets custom headers applied to every request.
// The provided map is copied to avoid unexpected mutations by callers.
func WithHeaders(h http.Header) RequestOption {
	return func(o *RequestOptions) {
		o.Headers = mergeHeaders(o.Headers, h)
	}
}

// WithHeader returns a RequestOption that sets a single header applied to every request.
func WithHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		o.Headers = mergeHeaders(o.Headers, http.Header{key: []string{value}})
	}
}

// WithDefaultHeader returns a RequestOption that sets a header only if missing or empty.
func WithDefaultHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		o.Headers = ensureHeader(o.Headers, key, value)
	}
}

func cloneHeader(h http.Header) http.Header {
	if len(h) == 0 {
		return nil
	}
	out := make(http.Header, len(h))
	for k, v := range h {
		key := textproto.CanonicalMIMEHeaderKey(k)
		out[key] = append([]string(nil), v...)
	}
	return out
}

func mergeHeaders(base http.Header, override http.Header) http.Header {
	if len(base) == 0 && len(override) == 0 {
		return nil
	}
	out := cloneHeader(base)
	if out == nil {
		out = make(http.Header, len(override))
	}
	for k, v := range override {
		key := textproto.CanonicalMIMEHeaderKey(k)
		out[key] = append([]string(nil), v...)
	}
	return out
}
