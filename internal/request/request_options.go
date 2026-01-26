package request

import (
	"context"
	"net/http"
)

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

func NewWithDefault() RequestOptions {
	return RequestOptions{
		Ctx:       context.Background(),
		BaseURL:   DefaultBaseURL,
		Token:     DefaultToken,
		Model:     DefaultModel,
		Provider:  DefaultProvider,
		Transport: http.DefaultClient,
	}
}

type RequestOption func(*RequestOptions)

func (o *RequestOptions) Apply(opts ...RequestOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *RequestOptions) NewOverride(opts ...RequestOption) *RequestOptions {
	newOpts := *o
	newOpts.Apply(opts...)
	return &newOpts
}

func NewOverride(o RequestOptions, opts ...RequestOption) *RequestOptions {
	o.Apply(opts...)
	return &o
}

func NewFromDefault(opts ...RequestOption) *RequestOptions {
	return NewOverride(NewWithDefault(), opts...)
}
