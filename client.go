package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// Client represents a HuggingFace API client with configured request options.
// Client instances are immutable; options are fixed at creation time and never mutated.
// This keeps client usage safe across goroutines and avoids surprises from mutable state.
// If options include externally-owned pointers, callers must avoid mutating them after creation
// or ensure their own synchronization.
// Services capture a snapshot of these options when created.
type Client struct {
	opts request.RequestOptions
}

// NewClient creates a new Client instance with the provided request options.
// If no options are provided, default options will be used.
// Clients are immutable; to change options, create a new Client to keep calls deterministic.
func NewClient(opts ...RequestOption) Client {
	return Client{
		opts: request.NewRequestOptions().With(opts...),
	}
}

// Chat returns a ChatService instance configured with this client's options.
// The chat service provides methods for interacting with chat completion endpoints.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call Chat() per use instead of retaining the value.
func (c Client) Chat() ChatService {
	return newChatService(c.opts)
}

// Raw returns a RawService instance configured with this client's options.
// The raw service provides methods for sending raw HTTP requests to any desired endpoint.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call Raw() per use instead of retaining the value.
func (c Client) Raw() RawService {
	return newRawService(c.opts)
}
