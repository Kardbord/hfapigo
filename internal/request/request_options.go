package request

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"

	"github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/version"
)

// RequestOptions holds configuration settings for API requests.
// Built-in option helpers return a new value and defensively clone headers,
// while context and the HTTP client are shared as-is. Custom options should avoid
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
	HTTPClient           *http.Client
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

// DefaultHTTPClient returns a new HTTP client configured with a cloned default transport.
func DefaultHTTPClient() *http.Client {
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Client{}
	}
	return &http.Client{Transport: defaultTransport.Clone()}
}

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
		HTTPClient:           DefaultHTTPClient(),
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

// Validate returns a configuration error if the options are invalid.
func (o RequestOptions) Validate() error {
	if o.HTTPClient == nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: "http client is nil",
		}
	}
	parsedBase, err := url.Parse(o.BaseURL)
	if err != nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("invalid base URL %q", o.BaseURL),
			Err:     err,
		}
	}
	if parsedBase.Scheme == "" || parsedBase.Host == "" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("base URL must include scheme and host, got %q", o.BaseURL),
		}
	}
	if parsedBase.RawQuery != "" || parsedBase.Fragment != "" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("base URL must not include query or fragment, got %q", o.BaseURL),
		}
	}
	return nil
}

// clone returns a shallow copy with reference types safely duplicated where possible.
// Headers are deep-copied; context and HTTP client are shared by design.
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
// The base URL must not include query parameters or fragments.
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
	if s == "" {
		return o
	}
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

// WithDefaultHTTPClient returns a new RequestOptions instance that uses the default HTTP client.
func (o RequestOptions) WithDefaultHTTPClient() RequestOptions {
	o = o.clone()
	o.HTTPClient = DefaultHTTPClient()
	return o
}

// WithHTTPClientFactory returns a new RequestOptions instance with an http.Client created by the factory.
// The factory should return a fresh client value; avoid sharing mutable internals like Transport unless synchronized.
// If the factory is nil, the HTTP client is set to nil.
func (o RequestOptions) WithHTTPClientFactory(factory func() http.Client) RequestOptions {
	o = o.clone()
	if factory == nil {
		o.HTTPClient = nil
		return o
	}
	client := factory()
	o.HTTPClient = &client
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
// The base URL must not include query parameters or fragments.
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
		if s == "" {
			return
		}
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

// WithDefaultHTTPClient returns a RequestOption that sets the default HTTP client.
func WithDefaultHTTPClient() RequestOption {
	return func(o *RequestOptions) {
		o.HTTPClient = DefaultHTTPClient()
	}
}

// WithHTTPClientFactory returns a RequestOption that sets a client created by the factory.
// The factory should return a fresh client value; avoid sharing mutable internals like Transport unless synchronized.
// If the factory is nil, the HTTP client is set to nil.
func WithHTTPClientFactory(factory func() http.Client) RequestOption {
	return func(o *RequestOptions) {
		if factory == nil {
			o.HTTPClient = nil
			return
		}
		client := factory()
		o.HTTPClient = &client
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

// cloneHeader returns a deep copy of the provided headers with canonicalized keys.
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
