package request

import (
	"context"
	"net/http"
)

type OptionProvider interface {
	Options() RequestOptions
}

type RequestOptions struct {
	Ctx       context.Context
	BaseURL   string
	Token     string
	Model     string
	Provider  string
	Transport Transport
}

const (
	DefaultBaseURL  = "https://router.huggingface.co"
	DefaultToken    = ""
	DefaultModel    = ""
	DefaultProvider = ""
)

func NewRequestOptions() RequestOptions {
	return RequestOptions{
		Ctx:       context.Background(),
		BaseURL:   DefaultBaseURL,
		Token:     DefaultToken,
		Model:     DefaultModel,
		Provider:  DefaultProvider,
		Transport: NewHTTPTransport(http.DefaultClient),
	}
}

type RequestOption func(*RequestOptions)

func (o *RequestOptions) apply(opts ...RequestOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o RequestOptions) With(opts ...RequestOption) RequestOptions {
	o.apply(opts...)
	return o
}
