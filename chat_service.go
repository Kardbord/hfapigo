package hfapigo

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// ChatService provides methods for interacting with chat completion endpoints.
type chatService struct {
	opts request.RequestOptions
}

func newChatService(opts request.RequestOptions) chatService {
	return chatService{opts: opts}
}

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
