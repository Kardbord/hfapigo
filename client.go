package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// Client represents a HuggingFace API client with configured request options.
// Client instances are immutable; options are fixed at creation time and never mutated.
// This keeps client usage safe across goroutines and avoids surprises from mutable state.
// If options include externally-owned pointers, callers must avoid mutating them after creation
// or ensure their own synchronization.
// RawService captures a snapshot of these options when created.
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
// The request is passed by value and the SDK never mutates the received
// payload. The value copy shares the request's nested data (slices, maps, and
// pointed-to values) with the caller, so the caller must treat the request and
// the data it references as read-only while a call is in flight.
//
// Concurrency:
//   - A single Client is safe for concurrent use.
//   - Reusing one request across sequential, fully-awaited calls is safe.
//   - To invoke the same template request from multiple goroutines, pass a
//     defensive copy per call, e.g. go client.Chat(req.Clone(), ...), or build a
//     fresh request per call.
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
//   - Returns a configuration error if the request is missing a model or messages.
//   - Returns a configuration error if *req.Stream is true; use ChatStream for streaming.
//
//nolint:gocritic // hugeParam: Chat takes the request by value so the SDK never mutates the caller's payload
func (c Client) Chat(req ChatRequest, opts ...Option) (ChatResponse, error) {
	return newChatService(c.opts).complete(req, opts...)
}

// ChatStream sends a chat completion request and returns a streaming response.
// Callers should Close the returned ChatStream when finished so the underlying HTTP
// connection and decoder goroutine are released promptly.
//
// The request is passed by value and the SDK never mutates the received
// payload. The value copy shares the request's nested data (slices, maps, and
// pointed-to values) with the caller, so the caller must treat the request and
// the data it references as read-only while a call is in flight.
//
// Concurrency:
//   - A single Client is safe for concurrent use.
//   - Reusing one request across sequential, fully-awaited calls is safe.
//   - To invoke the same template request from multiple goroutines, pass a
//     defensive copy per call, e.g. go client.ChatStream(req.Clone(), ...), or
//     build a fresh request per call.
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
//   - Returns a configuration error if the request is missing a model or messages.
//   - Always sends the request with streaming enabled.
//
//nolint:gocritic // hugeParam: ChatStream takes the request by value so the SDK never mutates the caller's payload
func (c Client) ChatStream(req ChatRequest, opts ...Option) (*ChatStream, error) {
	return newChatService(c.opts).completeStream(req, opts...)
}

// ClassifyText sends a text classification request and returns the text
// classification response for a single input.
//
// For multiple classification inputs, use ClassifyTextBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) ClassifyText(
	req TextClassificationRequest,
	opts ...Option,
) ([]TextClassification, error) {
	return newTextClassificationService(c.opts).classify(req, opts...)
}

// ClassifyTextBatch sends a text classification request for a batch of inputs
// and returns a list of text classification responses for each input in the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) ClassifyTextBatch(
	req TextClassificationBatchRequest,
	opts ...Option,
) ([][]TextClassification, error) {
	return newTextClassificationService(c.opts).classifyBatch(req, opts...)
}

// ZeroShotClassifyText sends a zero-shot text classification request and
// returns the zero-shot text classification response for a single input.
//
// For multiple inputs, use ZeroShotClassifyTextBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) ZeroShotClassifyText(
	req ZeroShotTextClassificationRequest,
	opts ...Option,
) ([]ZeroShotTextClassification, error) {
	return newZeroShotTextClassificationService(c.opts).classify(req, opts...)
}

// ZeroShotClassifyTextBatch sends a zero-shot text classification request for
// a batch of inputs and returns a list of zero-shot text classification
// responses for each input in the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) ZeroShotClassifyTextBatch(
	req ZeroShotTextClassificationBatchRequest,
	opts ...Option,
) ([][]ZeroShotTextClassification, error) {
	return newZeroShotTextClassificationService(c.opts).classifyBatch(req, opts...)
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
func (c Client) FillMaskBatch(
	req FillMaskBatchRequest,
	opts ...Option,
) ([][]FillMaskPrediction, error) {
	return newFillMaskService(c.opts).fillBatch(req, opts...)
}

// Summarize sends a summarization request and returns the summarization output
// for a single input.
//
// The API always returns a list for summarization; a single input yields a
// one-element list rather than a bare summary object.
//
// For multiple inputs, use SummarizeBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) Summarize(req SummarizationRequest, opts ...Option) ([]Summarization, error) {
	return newSummarizationService(c.opts).summarize(req, opts...)
}

// SummarizeBatch sends a summarization request for a batch of inputs and returns
// a flat list of summarization outputs, one for each input in the batch, in the
// same order as the inputs.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice. The response is
// a flat list (one summary per input) — not a nested list — consistent with
// how the API returns a list even for a single input.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) SummarizeBatch(
	req SummarizationBatchRequest,
	opts ...Option,
) ([]Summarization, error) {
	return newSummarizationService(c.opts).summarizeBatch(req, opts...)
}

// Translate sends a translation request and returns the translation output
// for a single input.
//
// The API always returns a list for translation; a single input yields a
// one-element list rather than a bare translation object.
//
// For multiple inputs, use TranslateBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) Translate(req TranslationRequest, opts ...Option) ([]Translation, error) {
	return newTranslationService(c.opts).translate(req, opts...)
}

// TranslateBatch sends a translation request for a batch of inputs and returns
// a flat list of translation outputs, one for each input in the batch, in the
// same order as the inputs.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice. The response is
// a flat list (one translation per input) — not a nested list — consistent with
// how the API returns a list even for a single input.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (c Client) TranslateBatch(
	req TranslationBatchRequest,
	opts ...Option,
) ([]Translation, error) {
	return newTranslationService(c.opts).translateBatch(req, opts...)
}

// Raw returns the raw HTTP request service for this client. Unlike the other
// endpoints, which are exposed directly as Client methods, the raw path remains
// namespaced under RawService: it is the advanced escape hatch for endpoints the
// SDK does not otherwise cover, and its several method variants are easier to
// discover grouped together than splashed across the Client surface.
//
// RawService is immutable and captures a snapshot of the client options when
// created; it is lightweight, so prefer calling Raw() per use rather than
// retaining the value.
func (c Client) Raw() RawService {
	return newRawService(c.opts)
}
