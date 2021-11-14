package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	MaxCandidateLabels       = 10
	RecommendedZeroShotModel = "facebook/bart-large-mnli"
)

// One of the following fields is required:
//   Input
//   Inputs
type ZeroShotRequest struct {
	// (Required) Input or Inputs are required request fields
	Inputs []string `json:"inputs"`

	// (Required)
	Parameters ZeroShotParameters `json:"parameters"`

	Options Options `json:"options,omitempty"`
}

// Used with ZeroShotRequest
type ZeroShotParameters struct {
	// (Required) A list of strings that are potential classes for inputs. Max 10 candidate_labels,
	// for more, simply run multiple requests, results are going to be misleading if using
	// too many candidate_labels anyway. If you want to keep the exact same, you can
	// simply run multi_label=True and do the scaling on your end.
	CandidateLabels []string `json:"candidate_labels"`

	// (Default: false) Boolean that is set to True if classes can overlap
	MultiLabel *bool `json:"multi_label,omitempty"`
}

func (zsp *ZeroShotParameters) SetMultiLabel(multiLabel bool) *ZeroShotParameters {
	zsp.MultiLabel = &multiLabel
	return zsp
}

type ZeroShotResponse struct {
	// The string sent as an input
	Sequence string `json:"sequence"`

	// The list of labels sent in the request, sorted in descending order
	// by probability that the input corresponds to the to the label.
	Labels []string `json:"labels"`

	// a list of floats that correspond the the probability of label, in the same order as labels.
	Scores []float64 `json:"scores"`
}

func SendZeroShotRequest(request *ZeroShotRequest, endpoint string) ([]*ZeroShotResponse, error) {
	if request == nil {
		return nil, errors.New("nil ZeroShotRequestMultiInput")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonBuf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	SetAuthorizationHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
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