package service

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/api"
	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// ChatService provides methods for interacting with chat completion endpoints.
type ChatService struct {
	opts request.RequestOptions
}

func NewChatService(opts request.RequestOptions) ChatService {
	return ChatService{opts: opts}
}

func (s ChatService) Complete(prompt string, opts ...api.RequestOption) (api.ChatResponse, error) {
	return request.DoJSON[api.ChatRequest, api.ChatResponse](
		s.opts.With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		api.ChatRequest{Inputs: prompt},
	)
}
