package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/request"
	"github.com/Kardbord/hfapigo/v4/services/chat"
)

// DefaultClient is a pre-configured client instance that can be used for quick API calls
// without needing to create a new client.
var DefaultClient = NewClient()

// Client represents a HuggingFace API client with configured request options.
type Client struct {
	opts request.RequestOptions
}

// NewClient creates a new Client instance with the provided request options.
// If no options are provided, default options will be used.
func NewClient(opts ...request.RequestOption) Client {
	return Client{
		opts: request.NewRequestOptions().With(opts...),
	}
}

// Chat returns a chat.Service instance configured with this client's options.
// The chat service provides methods for interacting with chat completion endpoints.
func (c Client) Chat() chat.Service {
	return chat.New(c)
}

// Options returns a copy of the client's request options.
func (c Client) Options() request.RequestOptions {
	return c.opts
}
