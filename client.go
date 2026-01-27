package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/request"
	"github.com/Kardbord/hfapigo/v4/services/chat"
)

var DefaultClient = NewClient()

type Client struct {
	opts request.RequestOptions
}

func NewClient(opts ...request.RequestOption) Client {
	return Client{
		opts: request.NewRequestOptions().With(opts...),
	}
}

func (c Client) Chat() chat.Service {
	return chat.New(c)
}

func (c Client) Options() request.RequestOptions {
	return c.opts
}
