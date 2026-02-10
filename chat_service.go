package hfapigo

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// chatService implements chat completion calls using the configured request options.
type chatService struct {
	opts request.RequestOptions
}

// newChatService builds a chat service with a snapshot of the provided options.
func newChatService(opts request.RequestOptions) chatService {
	return chatService{opts: opts}
}

// Complete sends a chat completion request and returns a chat completion response.
func (s chatService) Complete(req ChatRequest, opts ...RequestOption) (ChatResponse, error) {
	return request.DoJSON[ChatRequest, ChatResponse](
		s.opts.With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		req,
	)
}
