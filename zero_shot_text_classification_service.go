package hfgo

import (
	"fmt"
	"net/http"

	"github.com/Kardbord/hfgo/v4/internal/request"
)

// ZeroShotTextClassificationService implements zero-shot
// text classification calls using the configured request
// options.
type ZeroShotTextClassificationService struct {
	opts request.Options
}

// newZeroShotTextClassificationService builds a ZeroShotTextClassificationService with
// a snapshot of the provided options.
func newZeroShotTextClassificationService(opts request.Options) ZeroShotTextClassificationService {
	return ZeroShotTextClassificationService{opts: opts}
}

// Classify sends a zero-shot text classification request and returns a zero-shot
// text classification response for a single input.
//
// For multiple inputs, use ClassifyBatch.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s ZeroShotTextClassificationService) Classify(
	req ZeroShotTextClassificationRequest,
	opts ...Option,
) ([]ZeroShotTextClassification, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind: SDKErrorKindConfiguration,
			//nolint:goconst // repeated string is incidental
			Message: "the model option must be set for text classification to succeed",
			Err:     nil,
		}
	}

	if req.Parameters == nil || len(req.Parameters.CandidateLabels) == 0 {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "candidate labels must be provided for zero-shot text classification",
			Err:     nil,
		}
	}

	resp, err := request.DoJSON[ZeroShotTextClassificationRequest, []ZeroShotTextClassification](
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

// ClassifyBatch sends a zero-shot text classification request for a batch of inputs
// and returns a list of zero-shot text classification responses for each input in the
// batch.
//
// NOTE: Batched inference is supported by the upstream API, but is not
// officially documented; behavior may change without notice.
//
// Callers should check the length of the response list before indexing.
//
// The Provider option is ignored for now, as hf-inference is currently the only supported provider.
func (s ZeroShotTextClassificationService) ClassifyBatch(
	req ZeroShotTextClassificationBatchRequest,
	opts ...Option,
) ([][]ZeroShotTextClassification, error) {
	optsOverride := s.opts.With(opts...)

	if optsOverride.Model == "" {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for text classification to succeed",
			Err:     nil,
		}
	}

	if req.Parameters == nil || len(req.Parameters.CandidateLabels) == 0 {
		return nil, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "candidate labels must be provided for zero-shot text classification",
			Err:     nil,
		}
	}

	resp, err := request.DoJSON[ZeroShotTextClassificationBatchRequest, []zeroShotTextClassificationBatched](
		optsOverride,
		http.MethodPost,
		"hf-inference/models/"+optsOverride.Model,
		req,
	)
	if err != nil {
		return nil, err
	}

	return normalizeZeroShotTextClassificationResponse(resp, req.Inputs)
}

// normalizeZeroShotTextClassificationResponse handles a quirk in the HuggingFace API where
// the response format is a list of labels and a list of scores.
//
//	[
//		{
//			"sequence:" "sequence 1",
//			"labels": ["label1", "label2"],
//			"scores": [score1, score2]
//		},
//		{
//			"sequence:" "sequence 2",
//			"labels": ["label1", "label2"],
//			"scores": [score1, score2]
//		}
//	]
//
// Because batched request handling is available but undocumented, we're taking the liberty for
// now of normalizing the response into the same format as text classification for consistency.
//
//	[ // This list is ordered the same as the input list
//		[
//			{"label": "label1", "score": score1},
//			{"label": "label2", "score": score2}
//		],
//		[
//			{"label": "label1", "score": score1},
//			{"label": "label2", "score": score2}
//		],
//	]
//
// NOTE: Should batched zero-shot classification become officially supported and documented,
// we'll want to simply return the same format as the upstream API.
func normalizeZeroShotTextClassificationResponse(
	resp []zeroShotTextClassificationBatched,
	inputs []string,
) ([][]ZeroShotTextClassification, error) {
	// Validate response length matches input length
	if len(resp) != len(inputs) {
		return nil, &SDKError{
			Kind: SDKErrorKindSerialization,
			Message: fmt.Sprintf(
				"response item count (%d) does not match input count (%d); API response format may have changed",
				len(resp),
				len(inputs),
			),
			Err: nil,
		}
	}

	result := make([][]ZeroShotTextClassification, len(resp))

	//nolint:varnamelen // "i" is commonly used as an outer loop index variable
	for i, item := range resp {
		// Validate that the response sequence matches the input
		if item.Sequence != inputs[i] {
			return nil, &SDKError{
				Kind: SDKErrorKindSerialization,
				Message: fmt.Sprintf(
					`response item %d sequence does not match input; expected "%q" but got "%q"; API response format may have changed or order is not preserved`,
					i,
					inputs[i],
					item.Sequence,
				),
				Err: nil,
			}
		}

		// Validate that labels and scores have matching lengths
		if len(item.Labels) != len(item.Scores) {
			return nil, &SDKError{
				Kind: SDKErrorKindSerialization,
				Message: fmt.Sprintf(
					"response item %d has mismatched label and score counts (labels: %d, scores: %d); API response format may have changed",
					i,
					len(item.Labels),
					len(item.Scores),
				),
				Err: nil,
			}
		}

		classifications := make([]ZeroShotTextClassification, len(item.Labels))
		for j := range item.Labels {
			classifications[j] = ZeroShotTextClassification{
				Label: item.Labels[j],
				Score: item.Scores[j],
			}
		}
		result[i] = classifications
	}

	return result, nil
}
