package hfapigo

import (
	"encoding/json"
	"errors"
)

const (
	MaxCandidateLabels       = 10
	RecommendedZeroShotModel = "facebook/bart-large-mnli"
)

// Request structure for the Zero-shot classification endpoint.
//
// One of the following fields is required:
//   Input
//   Inputs
type ZeroShotRequest struct {
	// (Required) Input or Inputs are required request fields
	Inputs []string `json:"inputs,omitempty"`

	// (Required)
	Parameters ZeroShotParameters `json:"parameters,omitempty"`

	Options Options `json:"options,omitempty"`
}

// Used with ZeroShotRequest
type ZeroShotParameters struct {
	// (Required) A list of strings that are potential classes for inputs. Max 10 candidate_labels,
	// for more, simply run multiple requests, results are going to be misleading if using
	// too many candidate_labels anyway. If you want to keep the exact same, you can
	// simply run multi_label=True and do the scaling on your end.
	CandidateLabels []string `json:"candidate_labels,omitempty"`

	// (Default: false) Boolean that is set to True if classes can overlap
	MultiLabel *bool `json:"multi_label,omitempty"`
}

func (zsp *ZeroShotParameters) SetMultiLabel(multiLabel bool) *ZeroShotParameters {
	zsp.MultiLabel = &multiLabel
	return zsp
}

// Response structure from the Zero-shot classification endpoint.
type ZeroShotResponse struct {
	// The string sent as an input
	Sequence string `json:"sequence,omitempty"`

	// The list of labels sent in the request, sorted in descending order
	// by probability that the input corresponds to the to the label.
	Labels []string `json:"labels,omitempty"`

	// a list of floats that correspond the the probability of label, in the same order as labels.
	Scores []float64 `json:"scores,omitempty"`
}

func SendZeroShotRequest(model string, request *ZeroShotRequest) ([]*ZeroShotResponse, error) {
	if request == nil {
		return nil, errors.New("nil ZeroShotRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	zresps := make([]*ZeroShotResponse, len(request.Inputs))
	err = json.Unmarshal(respBody, &zresps)
	if err != nil {
		return nil, err
	}

	return zresps, nil
}
