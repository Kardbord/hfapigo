package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// Client represents a HuggingFace API client with configured request options.
// Client instances are immutable; options are fixed at creation time and never mutated.
// This keeps client usage safe across goroutines and avoids surprises from mutable state.
// If options include externally-owned pointers, callers must avoid mutating them after creation
// or ensure their own synchronization.
// Services capture a snapshot of these options when created.
type Client struct {
	opts request.Options
}

// NewClient creates a new Client instance with the provided request options.
// If no options are provided, default options will be used.
// Clients are immutable; to change options, create a new Client to keep calls deterministic.
func NewClient(opts ...Option) Client {
	return Client{
		opts: request.NewOptions().With(opts...),
	}
}

// Chat sends a chat completion request and returns a chat completion response.
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
//
// Behavior:
//   - Returns a configuration error if req is nil.
//   - Returns a configuration error if *req.Stream is true; use ChatStream for streaming.
func (c Client) Chat(req *ChatRequest, opts ...Option) (ChatResponse, error) {
	return newChatService(c.opts).complete(req, opts...)
}

// ChatStream sends a chat completion request and returns a streaming response.
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
//
// Behavior:
//   - Returns a configuration error if req is nil.
//   - Always sends the request with streaming enabled.
func (c Client) ChatStream(req *ChatRequest, opts ...Option) (*ChatStream, error) {
	return newChatService(c.opts).completeStream(req, opts...)
}

// ClassifyText returns a TextClassificationService instance configured with this client's options.
// The classification service provides methods for interacting with text classification endpoints.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call ClassifyText() per use instead of retaining the value.
func (c Client) ClassifyText() TextClassificationService {
	return newTextClassificationService(c.opts)
}

// ZeroShotClassifyText returns a ZeroShotClassificationService instance configured with this client's options.
// The classification service provides methods for interacting with zero-shot text classification endpoints.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call ZeroShotClassifyText() per use instead of retaining the value.
func (c Client) ZeroShotClassifyText() ZeroShotTextClassificationService {
	return newZeroShotTextClassificationService(c.opts)
}

// FillMask sends a fill mask request and returns the mask filling predictions
// for a single input.
//
// For multiple inputs, use FillMaskBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) FillMask(req FillMaskRequest, opts ...Option) ([]FillMaskPrediction, error) {
	return newFillMaskService(c.opts).fill(req, opts...)
}

// FillMaskBatch sends a fill mask request for a batch of inputs and returns a
// list of mask filling predictions for each input in the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) FillMaskBatch(req FillMaskBatchRequest, opts ...Option) ([][]FillMaskPrediction, error) {
	return newFillMaskService(c.opts).fillBatch(req, opts...)
}

// Summarization returns a SummarizationService instance configured with this client's options.
// The summarization service provides methods for interacting with summarization endpoints.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call Summarization() per use instead of retaining the value.
func (c Client) Summarization() SummarizationService {
	return newSummarizationService(c.opts)
}

// Raw returns a RawService instance configured with this client's options.
// The raw service provides methods for sending raw HTTP requests to any desired endpoint.
// Service configurations are captured at creation time and do not change if the client options change later.
// Clients are immutable to keep concurrency simple and request behavior predictable.
// Services are lightweight; prefer to call Raw() per use instead of retaining the value.
func (c Client) Raw() RawService {
	return newRawService(c.opts)
}
