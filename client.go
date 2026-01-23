package hfapigo

import (
	"net/http"
	"time"
)

type client struct {
	options clientOptions
}

func NewClient(opts ...clientOption) (*client, error) {
	client := &client{
		options: clientOptions{
			baseURL:   DefaultBaseURL,
			token:     DefaultToken,
			transport: DefaultTransport(),
			model:     DefaultModel,
			provider:  DefaultProvider,
			timeout:   DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(&client.options)
	}

	return client, nil
}

type clientOptions struct {
	baseURL   string
	token     string
	transport Transport
	model     string
	provider  string
	timeout   time.Duration
	// TODO: Add streaming support
}

const (
	DefaultBaseURL  = "https://router.huggingface.co"
	DefaultToken    = ""
	DefaultModel    = ""
	DefaultProvider = ""
	DefaultTimeout  = time.Second * 5
)

func DefaultTransport() Transport {
	return &HTTPTransport{http.DefaultClient}
}

type clientOption func(*clientOptions)

func WithBaseURL(u string) clientOption {
	return func(o *clientOptions) { o.baseURL = u }
}

func WithToken(t string) clientOption {
	return func(o *clientOptions) { o.token = t }
}

func WithTransport(t Transport) clientOption {
	return func(o *clientOptions) { o.transport = t }
}

func WithModel(m string) clientOption {
	return func(o *clientOptions) { o.model = m }
}

func WithProvider(p string) clientOption {
	return func(o *clientOptions) { o.provider = p }
}

func WithTimeout(t time.Duration) clientOption {
	return func(o *clientOptions) { o.timeout = t }
}
