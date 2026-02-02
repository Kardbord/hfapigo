package hfapigo

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// ChatService provides methods for interacting with chat completion endpoints.
type ChatService interface {
	Complete(prompt string, opts ...request.RequestOption) (ChatResponse, error)
}

// ChatResponse represents a completion response from the chat API.
type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}

type chatRequest struct {
	Inputs string `json:"inputs"`
}

type chatService struct {
	opts request.RequestOptions
}

func newChatService(opts request.RequestOptions) ChatService {
	return chatService{opts: opts}
}

func (s chatService) Complete(prompt string, opts ...request.RequestOption) (ChatResponse, error) {
	return request.DoJSON[chatRequest, ChatResponse](
		s.opts.With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		chatRequest{Inputs: prompt},
	)
}
