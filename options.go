package hfapigo

import (
	"context"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

func WithBaseURL(u string) request.RequestOption {
	return func(o *request.RequestOptions) { o.BaseURL = u }
}

func WithToken(t string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Token = t }
}

func WithModel(m string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Model = m }
}

func WithProvider(p string) request.RequestOption {
	return func(o *request.RequestOptions) { o.Provider = p }
}

func WithHTTPClient(c *http.Client) request.RequestOption {
	return func(o *request.RequestOptions) { o.Transport = c }
}

func WithContext(ctx context.Context) request.RequestOption {
	return func(o *request.RequestOptions) { o.Ctx = ctx }
}
