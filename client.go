package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/request"
	"github.com/Kardbord/hfapigo/v4/services/chat"
)

type Client struct {
	opts request.RequestOptions
}

func NewClient(opts ...request.RequestOption) (*Client, error) {
	return &Client{
		opts: *request.NewFromDefault(opts...),
	}, nil
}

func (c *Client) Chat() *chat.Service {
	return chat.New(&c.opts)
}
