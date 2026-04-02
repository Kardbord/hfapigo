package hfgo

import (
	"net/http"

	"github.com/Kardbord/hfgo/v4/internal/request"
)

// TextClassificationService implements text classification
// calls using the configured request options.
type TextClassificationService struct {
	opts request.Options
}

// newTextClassificationService builds a TextClassificationService with a snapshot
// of the provided options.
func newTextClassificationService(opts request.Options) TextClassificationService {
	return TextClassificationService{opts: opts}
}

// Classify sends a text classification request and returns a text
// classification response for a single input.
//
// For multiple classification inputs, use ClassifyBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s TextClassificationService) Classify(
	req TextClassificationRequest,
	opts ...Option,
) ([]TextClassification, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for text classification to succeed",
			Err:     nil,
		}
	}

	// NOTE: The API documentation indicates that this should return an JSON array of
	// TextClassification objects, but in reality it returns an array of arrays, where
	// the outer array contains only a single entry (the inner array), and the inner array
	// contains a list of TextClassification objects.
	resp, err := request.DoJSON[TextClassificationRequest, [][]TextClassification](
		optsOverride,
		http.MethodPost,
		"hf-inference/models/"+optsOverride.Model,
		req,
	)
	if err != nil {
		return nil, err
	}

	// Unexpected, but technically legal API response.
	if len(resp) < 1 {
		return nil, nil
	}

	return resp[0], nil
}

// ClassifyBatch sends a text classification request for a batch of inputs
// and returns a list of text classification responses for each input in
// the batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s TextClassificationService) ClassifyBatch(
	req TextClassificationBatchRequest,
	opts ...Option,
) ([][]TextClassification, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for text classification to succeed",
			Err:     nil,
		}
	}

	resp, err := request.DoJSON[TextClassificationBatchRequest, [][]TextClassification](
		optsOverride,
		http.MethodPost,
		"hf-inference/models/"+optsOverride.Model,
		req,
	)
	if err != nil {
		return nil, err
	}

	return normalizeTextClassificationResponse(resp, len(req.Inputs)), nil
}

// normalizeTextClassificationResponse handles a quirk in the HuggingFace API where
// the response format differs based on whether the TopK parameter is explicitly set:
//
//   - When TopK is explicitly set (e.g., to 1, 2, or any value): Returns
//     [[classifications for input1], [classifications for input2], ...] (per-input format)
//   - When TopK is unset (nil): Returns [[all classifications together]] (flat format)
//
// This function detects the flat format case and reshapes it to the expected per-input
// format for API consistency. The detection heuristic is:
// - Single outer array (len(resp) == 1)
// - Number of inner classifications equals number of inputs (len(resp[0]) == numInputs)
// - More than one input was sent (numInputs > 1)
//
// When all conditions are met, we reshape [[class1, class2, class3]] into
// [[class1], [class2], [class3]] to maintain consistent per-input structure.
func normalizeTextClassificationResponse(
	resp [][]TextClassification,
	numInputs int,
) [][]TextClassification {
	if numInputs > 1 && len(resp) == 1 && len(resp[0]) == numInputs {
		// Reshape flat format to per-input format
		reshaped := make([][]TextClassification, numInputs)
		for i := range numInputs {
			reshaped[i] = []TextClassification{resp[0][i]}
		}

		return reshaped
	}

	return resp
}
