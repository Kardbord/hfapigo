package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// FillMaskService implements fill mask calls using the configured request options.
type FillMaskService struct {
	opts request.Options
}

// newFillMaskService builds a FillMaskService with a snapshot of the provided options.
func newFillMaskService(opts request.Options) FillMaskService {
	return FillMaskService{opts: opts}
}

// Fill sends a fill mask request and returns the mask filling
// predictions for a single input.
//
// For multiple inputs, use FillBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s FillMaskService) Fill(
	req FillMaskRequest,
	opts ...Option,
) ([]FillMaskPrediction, error) {
	return doModelInference[FillMaskRequest, []FillMaskPrediction](
		s.opts.With(opts...),
		"fill mask",
		req,
	)
}

// FillBatch sends a fill mask request for a batch of inputs
// and returns a list of mask filling predictions for each input in
// the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s FillMaskService) FillBatch(
	req FillMaskBatchRequest,
	opts ...Option,
) ([][]FillMaskPrediction, error) {
	return doModelInference[FillMaskBatchRequest, [][]FillMaskPrediction](
		s.opts.With(opts...),
		"fill mask",
		req,
	)
}
