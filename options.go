package hfapigo

import (
	"net/http"
)

type ClientOption func(*clientOptions)

type clientOptions struct {
	baseURL  string
	token    string
	model    string
	provider string
	httpc    *http.Client
}

func defaultClientOptions() clientOptions {
	return clientOptions{
		baseURL:  DefaultBaseURL,
		token:    DefaultToken,
		model:    DefaultModel,
		provider: DefaultProvider,
		httpc:    http.DefaultClient,
	}
}

const (
	DefaultBaseURL  = "https://router.huggingface.co"
	DefaultToken    = ""
	DefaultModel    = ""
	DefaultProvider = ""
)

func WithBaseURL(u string) ClientOption {
	return func(o *clientOptions) { o.baseURL = u }
}

func WithToken(t string) ClientOption {
	return func(o *clientOptions) { o.token = t }
}

func WithModel(m string) ClientOption {
	return func(o *clientOptions) { o.model = m }
}

func WithProvider(p string) ClientOption {
	return func(o *clientOptions) { o.provider = p }
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(o *clientOptions) { o.httpc = c }
}
