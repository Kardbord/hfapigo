package hfapigo

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// ChatResponse represents a completion response from the chat API.
type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}

type chatRequest struct {
	Inputs string `json:"inputs"`
}

// ChatService provides methods for interacting with chat completion endpoints.
type ChatService struct {
	opts request.RequestOptions
}

func newChatService(opts request.RequestOptions) ChatService {
	return ChatService{opts: opts}
}

func (s ChatService) Complete(prompt string, opts ...request.RequestOption) (ChatResponse, error) {
	return request.DoJSON[chatRequest, ChatResponse](
		s.opts.With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		chatRequest{Inputs: prompt},
	)
}
