package hfgo

import (
	"net/http"

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

// FillMask sends a fill mask request and returns the mask filling
// predictions for a single input.
//
// For multiple inputs, use FillMaskBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s FillMaskService) FillMask(
	req FillMaskRequest,
	opts ...Option,
) ([]FillMaskPrediction, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for the fill mask request to succeed",
			Err:     nil,
		}
	}

	resp, err := request.DoJSON[FillMaskRequest, []FillMaskPrediction](
		optsOverride,
		http.MethodPost,
		"hf-inference/models/"+optsOverride.Model,
		req,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FillMaskBatch sends a fill mask request for a batch of inputs
// and returns a list of mask filling predictions for each input in
// the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s FillMaskService) FillMaskBatch(
	req FillMaskBatchRequest,
	opts ...Option,
) ([][]FillMaskPrediction, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for the fill mask request to succeed",
			Err:     nil,
		}
	}

	resp, err := request.DoJSON[FillMaskBatchRequest, [][]FillMaskPrediction](
		optsOverride,
		http.MethodPost,
		"hf-inference/models/"+optsOverride.Model,
		req,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
