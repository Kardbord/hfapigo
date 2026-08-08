package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// SummarizationService implements summarization calls
// using the configured request options.
type SummarizationService struct {
	opts request.Options
}

// newSummarizationService builds a SummarizationService with a snapshot
// of the provided options.
func newSummarizationService(opts request.Options) SummarizationService {
	return SummarizationService{opts: opts}
}

// Summarize sends a summarization request and returns the
// summarization output for a single input.
//
// The API always returns a list for summarization; a single input yields a
// one-element list rather than a bare summary object.
//
// For multiple inputs, use SummarizeBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s SummarizationService) Summarize(
	req SummarizationRequest,
	opts ...Option,
) ([]Summarization, error) {
	return doModelInference[SummarizationRequest, []Summarization](
		s.opts.With(opts...),
		"summarization",
		req,
	)
}

// SummarizeBatch sends a summarization request for a batch of inputs
// and returns a flat list of summarization outputs, one for each input in
// the batch, in the same order as the inputs.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice. The response is
// a flat list (one summary per input) — not a nested list — consistent with
// how the API returns a list even for a single input.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s SummarizationService) SummarizeBatch(
	req SummarizationBatchRequest,
	opts ...Option,
) ([]Summarization, error) {
	return doModelInference[SummarizationBatchRequest, []Summarization](
		s.opts.With(opts...),
		"summarization",
		req,
	)
}
