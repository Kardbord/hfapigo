package chat

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

type Service struct {
	ctx *request.Context
}

func New(ctx *request.Context) *Service {
	return &Service{ctx: ctx}
}

type ChatRequest struct {
	Inputs string `json:"inputs"`
}

type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}

func (s *Service) Complete(prompt string) (ChatResponse, error) {
	return request.DoJSON[ChatRequest, ChatResponse](
		s.ctx,
		http.MethodPost,
		"/v1/chat/completions",
		ChatRequest{Inputs: prompt},
	)
}
