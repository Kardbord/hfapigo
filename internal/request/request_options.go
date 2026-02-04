package request

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"

	"github.com/Kardbord/hfapigo/v4/internal/version"
)

// RequestOptions holds configuration settings for API requests.
// Built-in option helpers return a new value and defensively clone headers,
// while context and transport are shared as-is. Custom options should avoid
// reusing mutable header maps if they want the same defensive-copy behavior.
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
// Custom options that set Headers should avoid reusing mutable header maps
// if they want to preserve the defensive-copy behavior of built-in helpers.
type RequestOption func(*RequestOptions)

// apply applies a series of RequestOption functions to the RequestOptions instance.
func (o *RequestOptions) apply(opts ...RequestOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// clone returns a shallow copy with reference types safely duplicated where possible.
// Headers are deep-copied; context and transport are shared by design.
func (o RequestOptions) clone() RequestOptions {
	o.Headers = cloneHeader(o.Headers)
	return o
}

// With returns a new RequestOptions instance with the provided options applied.
// This method creates a copy of the current options and applies modifications to it.
func (o RequestOptions) With(opts ...RequestOption) RequestOptions {
	o = o.clone()
	o.apply(opts...)
	return o
}

// WithBaseURL returns a new RequestOptions instance with the base URL updated.
func (o RequestOptions) WithBaseURL(u string) RequestOptions {
	o = o.clone()
	o.BaseURL = u
	return o
}

// WithToken returns a new RequestOptions instance with the authentication token updated.
func (o RequestOptions) WithToken(t string) RequestOptions {
	o = o.clone()
	o.Token = t
	return o
}

// WithModel returns a new RequestOptions instance with the model updated.
func (o RequestOptions) WithModel(m string) RequestOptions {
	o = o.clone()
	o.Model = m
	return o
}

// WithProvider returns a new RequestOptions instance with the provider updated.
func (o RequestOptions) WithProvider(p string) RequestOptions {
	o = o.clone()
	o.Provider = p
	return o
}

// WithUserAgent returns a new RequestOptions instance with the User-Agent value updated.
func (o RequestOptions) WithUserAgent(ua string) RequestOptions {
	o = o.clone()
	o.UserAgent = ua
	return o
}

// WithUserAgentSuffix returns a new RequestOptions instance with a suffix appended to the SDK User-Agent.
func (o RequestOptions) WithUserAgentSuffix(s string) RequestOptions {
	o = o.clone()
	base := o.UserAgent
	if base == "" {
		base = version.UserAgent()
	}
	o.UserAgent = fmt.Sprintf("%s %s", base, s)
	return o
}

// WithMaxResponseBodyBytes returns a new RequestOptions instance with the response size cap updated.
func (o RequestOptions) WithMaxResponseBodyBytes(n int64) RequestOptions {
	o = o.clone()
	o.MaxResponseBodyBytes = n
	return o
}

// WithContext returns a new RequestOptions instance with the context updated.
func (o RequestOptions) WithContext(ctx context.Context) RequestOptions {
	o = o.clone()
	o.Ctx = ctx
	return o
}

// WithHTTPClient returns a new RequestOptions instance with the transport updated from the HTTP client.
func (o RequestOptions) WithHTTPClient(c *http.Client) RequestOptions {
	o = o.clone()
	o.Transport = NewHTTPTransport(c)
	return o
}

// WithTransport returns a new RequestOptions instance with the transport updated.
func (o RequestOptions) WithTransport(t Transport) RequestOptions {
	o = o.clone()
	o.Transport = t
	return o
}

// WithHeaders returns a new RequestOptions instance with the provided headers applied,
// overriding any existing values for matching keys.
func (o RequestOptions) WithHeaders(h http.Header) RequestOptions {
	o.Headers = overrideHeaders(o.Headers, h)
	return o
}

// WithHeader returns a new RequestOptions instance with a single header applied.
func (o RequestOptions) WithHeader(key, value string) RequestOptions {
	o.Headers = overrideHeaders(o.Headers, http.Header{key: []string{value}})
	return o
}

// WithDefaultHeader returns a new RequestOptions instance with a header set
// only if the header is missing or empty.
func (o RequestOptions) WithDefaultHeader(key, value string) RequestOptions {
	o.Headers = ensureHeader(o.Headers, key, value)
	return o
}

// WithBaseURL returns a RequestOption that sets the base URL for API requests.
func WithBaseURL(u string) RequestOption {
	return func(o *RequestOptions) {
		o.BaseURL = u
	}
}

// WithToken returns a RequestOption that sets the authentication token for API requests.
func WithToken(t string) RequestOption {
	return func(o *RequestOptions) {
		o.Token = t
	}
}

// WithModel returns a RequestOption that sets the model to use for API requests.
func WithModel(m string) RequestOption {
	return func(o *RequestOptions) {
		o.Model = m
	}
}

// WithProvider returns a RequestOption that sets the provider for API requests.
func WithProvider(p string) RequestOption {
	return func(o *RequestOptions) {
		o.Provider = p
	}
}

// WithUserAgent returns a RequestOption that sets the User-Agent header value.
func WithUserAgent(ua string) RequestOption {
	return func(o *RequestOptions) {
		o.UserAgent = ua
	}
}

// WithUserAgentSuffix returns a RequestOption that appends a suffix to the SDK user agent string.
func WithUserAgentSuffix(s string) RequestOption {
	return func(o *RequestOptions) {
		base := o.UserAgent
		if base == "" {
			base = version.UserAgent()
		}
		o.UserAgent = fmt.Sprintf("%s %s", base, s)
	}
}

// WithMaxResponseBodyBytes returns a RequestOption that sets the maximum response size to read.
func WithMaxResponseBodyBytes(n int64) RequestOption {
	return func(o *RequestOptions) {
		o.MaxResponseBodyBytes = n
	}
}

// WithContext returns a RequestOption that sets the context for API requests.
func WithContext(ctx context.Context) RequestOption {
	return func(o *RequestOptions) {
		o.Ctx = ctx
	}
}

// WithHTTPClient returns a RequestOption that sets a custom HTTP client for API requests.
func WithHTTPClient(c *http.Client) RequestOption {
	return func(o *RequestOptions) {
		o.Transport = NewHTTPTransport(c)
	}
}

// WithTransport returns a RequestOption that sets a custom transport for API requests.
func WithTransport(t Transport) RequestOption {
	return func(o *RequestOptions) {
		o.Transport = t
	}
}

// WithHeaders returns a RequestOption that sets custom headers applied to every request,
// overriding any existing values for matching keys.
// The provided map is copied to avoid unexpected mutations by callers.
func WithHeaders(h http.Header) RequestOption {
	return func(o *RequestOptions) {
		o.Headers = overrideHeaders(o.Headers, h)
	}
}

// WithHeader returns a RequestOption that sets a single header applied to every request.
func WithHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		o.Headers = overrideHeaders(o.Headers, http.Header{key: []string{value}})
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

// overrideHeaders copies base headers and replaces values with override entries.
// Header keys are canonicalized to avoid duplicate variants.
func overrideHeaders(base http.Header, override http.Header) http.Header {
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
