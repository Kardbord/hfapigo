package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/request"
	"github.com/Kardbord/hfapigo/v4/services/chat"
)

type Client struct {
	ctx request.RequestOptions
}

func NewClient(opts ...ClientOption) (*Client, error) {
	options := defaultClientOptions()

	for _, opt := range opts {
		opt(&options)
	}

	return &Client{
		ctx: request.RequestOptions{
			BaseURL:   options.baseURL,
			Token:     options.token,
			Model:     options.model,
			Provider:  options.provider,
			Transport: request.NewHTTPTransport(options.httpc),
		},
	}, nil
}

func (c *Client) Chat() *chat.Service {
	return chat.New(&c.ctx)
}
