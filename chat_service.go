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

// Complete sends a basic chat completion with a single user prompt.
func (s chatService) Complete(prompt string, opts ...RequestOption) (ChatResponse, error) {
	content := ChatMessageContent{Text: &prompt}
	return request.DoJSON[ChatRequest, ChatResponse](
		s.opts.With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		ChatRequest{
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: content,
				},
			},
		},
	)
}
