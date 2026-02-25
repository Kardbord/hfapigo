package hfapigo

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// EndpointChatCompletion specifies the chat completion endpoint.
const EndpointChatCompletion = "/v1/chat/completions"

// ChatService implements chat completion calls using the configured request options.
type ChatService struct {
	opts request.RequestOptions
}

// newChatService builds a chat service with a snapshot of the provided options.
func newChatService(opts request.RequestOptions) ChatService {
	return ChatService{opts: opts}
}

// Complete sends a chat completion request and returns a chat completion response.
func (s ChatService) Complete(req *ChatRequest, opts ...RequestOption) (ChatResponse, error) {
	if req == nil {
		return ChatResponse{}, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat request is nil",
			Err:     nil,
		}
	}

	payload := *req
	optsOverride := s.opts.With(opts...)
	if payload.Model == nil || *payload.Model == "" {
		if optsOverride.Model != "" {
			model := optsOverride.Model
			payload.Model = &model
		}
	}

	return request.DoJSON[ChatRequest, ChatResponse](
		optsOverride,
		http.MethodPost,
		EndpointChatCompletion,
		payload,
	)
}
