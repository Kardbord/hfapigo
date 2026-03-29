package hfgo

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kardbord/hfgo/v4/internal/chatstream"
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// EndpointChatCompletion specifies the chat completion endpoint.
const EndpointChatCompletion = "/v1/chat/completions"

// ChatService implements chat completion calls using the configured request options.
type ChatService struct {
	opts request.Options
}

// newChatService builds a chat service with a snapshot of the provided options.
func newChatService(opts request.Options) ChatService {
	return ChatService{opts: opts}
}

// Complete sends a chat completion request and returns a chat completion response.
//
// The caller must not mutate req while the request is being processed.
// For safe concurrent usage, create a new ChatRequest for each concurrent call.
//
// Model Precedence:
// The Model field is resolved with the following precedence (highest to lowest):
//  1. ChatRequest.Model field (if non-nil and non-empty)
//  2. Per-request options Model override
//  3. Client-level Model option
//
// Provider Precedence:
// The Provider field is applied as a fallback only if the resolved Model does not
// already contain a provider (indicated by ":" in the model string). If the Model
// is in the format "model:provider", the Provider option is ignored.
//
// For example:
//   - Model="mistral-7b", Provider="huggingface" → "mistral-7b:huggingface"
//   - Model="mistral-7b:huggingface", Provider="mistral" → "mistral-7b:huggingface" (Provider ignored)
//   - Model="mistral-7b:huggingface", Provider="" → "mistral-7b:huggingface"
func (s ChatService) Complete(req *ChatRequest, opts ...Option) (ChatResponse, error) {
	if req == nil {
		return ChatResponse{}, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat request is nil",
			Err:     nil,
		}
	}

	payload := *req
	optsOverride := s.opts.With(opts...)

	resolveModel(&payload, optsOverride)

	if payload.Stream != nil && *payload.Stream {
		return ChatResponse{}, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat completion streaming is not supported by ChatService.Complete; use a streaming method instead",
			Err:     nil,
		}
	}

	return request.DoJSON[ChatRequest, ChatResponse](
		optsOverride,
		http.MethodPost,
		EndpointChatCompletion,
		payload,
	)
}

// CompleteStream sends a chat completion request and returns a streaming response.
// Callers should Close the returned ChatStream when finished so the underlying HTTP
// connection and decoder goroutine are released promptly.
//
// The caller must not mutate req while the request is being processed.
// For safe concurrent usage, create a new ChatRequest for each concurrent call.
//
// Model Precedence:
// The Model field is resolved with the following precedence (highest to lowest):
//  1. ChatRequest.Model field (if non-nil and non-empty)
//  2. Per-request options Model override
//  3. Client-level Model option
//
// Provider Precedence:
// The Provider field is applied as a fallback only if the resolved Model does not
// already contain a provider (indicated by ":" in the model string). If the Model
// is in the format "model:provider", the Provider option is ignored.
//
// For example:
//   - Model="mistral-7b", Provider="huggingface" → "mistral-7b:huggingface"
//   - Model="mistral-7b:huggingface", Provider="mistral" → "mistral-7b:huggingface" (Provider ignored)
//   - Model="mistral-7b:huggingface", Provider="" → "mistral-7b:huggingface"
func (s ChatService) CompleteStream(req *ChatRequest, opts ...Option) (*ChatStream, error) {
	if req == nil {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat request is nil",
			Err:     nil,
		}
	}

	payload := *req
	optsOverride := s.opts.With(opts...)

	resolveModel(&payload, optsOverride)

	stream := true
	payload.Stream = &stream

	streamResp, err := request.DoJSONStream[ChatRequest, ChatStreamResponse](
		optsOverride,
		http.MethodPost,
		EndpointChatCompletion,
		payload,
	)
	if err != nil {
		return nil, err
	}

	return &ChatStream{
		stream:       streamResp,
		toolCallAccr: chatstream.ToolCallAccumulator{},
	}, nil
}

// resolveModel resolves the model with precedence and applies the provider fallback.
// Model precedence: request > options > client.
// Provider is applied as a fallback only if the resolved model doesn't already contain a provider.
func resolveModel(payload *ChatRequest, optsOverride request.Options) {
	// Resolve model with precedence: request > options > client
	if payload.Model == nil || *payload.Model == "" {
		if optsOverride.Model != "" {
			model := optsOverride.Model
			payload.Model = &model
		}
	}

	// Apply provider fallback to the final model
	payload.Model = applyProvider(payload.Model, optsOverride.Provider)
}

// applyProvider applies the provider to the model if the model
// doesn't already contain a provider (indicated by ":").
func applyProvider(model *string, provider string) *string {
	if model == nil || *model == "" || provider == "" {
		return model
	}

	if !strings.Contains(*model, ":") {
		newModel := fmt.Sprintf("%s:%s", *model, provider)

		return &newModel
	}

	return model
}

// ChatStream wraps a streaming chat completion response.
type ChatStream struct {
	stream       *request.JSONStream[ChatStreamResponse]
	toolCallAccr chatstream.ToolCallAccumulator
}

// Recv blocks until the next streaming chunk arrives or the context is done.
func (c *ChatStream) Recv(ctx context.Context) (ChatStreamResponse, error) {
	if c.stream == nil {
		return ChatStreamResponse{}, &SDKError{
			Kind:    SDKErrorKindInternal,
			Message: "chat stream is nil",
			Err:     nil,
		}
	}

	chunk, err := c.stream.Recv(ctx)
	if err != nil {
		return chunk, err
	}

	c.mergeToolCallMetadata(&chunk)

	return chunk, nil
}

// Close releases the underlying stream resources.
func (c *ChatStream) Close() error {
	if c == nil || c.stream == nil {
		return nil
	}

	return c.stream.Close()
}

// mergeToolCallMetadata ensures streaming tool call deltas include the cached
// id/type/function-name values observed earlier in the stream.
func (c *ChatStream) mergeToolCallMetadata(resp *ChatStreamResponse) {
	if c == nil || resp == nil {
		return
	}
	for i := range resp.Choices {
		choice := &resp.Choices[i]
		if len(choice.Delta.ToolCalls) == 0 {
			continue
		}
		for j := range choice.Delta.ToolCalls {
			call := &choice.Delta.ToolCalls[j]
			toolID, callType, functionName := c.toolCallAccr.Merge(
				choice.Index,
				call.Index,
				call.ID,
				call.Type,
				call.Function.Name,
			)
			call.ID = toolID
			call.Type = callType
			call.Function.Name = functionName
		}
	}
}
