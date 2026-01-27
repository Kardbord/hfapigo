package chat

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

type Service struct {
	client request.OptionProvider
}

func New(client request.OptionProvider) Service {
	return Service{client: client}
}

type ChatRequest struct {
	Inputs string `json:"inputs"`
}

type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}

func (s Service) Complete(prompt string, opts ...request.RequestOption) (ChatResponse, error) {
	return request.DoJSON[ChatRequest, ChatResponse](
		s.client.Options().With(opts...),
		http.MethodPost,
		"/v1/chat/completions",
		ChatRequest{Inputs: prompt},
	)
}
