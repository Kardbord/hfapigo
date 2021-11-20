package hfapigo

import (
	"encoding/json"
	"errors"
)

const (
	RecommendedTextClassificationModel = "distilbert-base-uncased-finetuned-sst-2-english"
)

// Request structure for the Text classification endpoint
type TextClassificationRequest struct {
	// (Required) strings to be classified
	Inputs  []string `json:"inputs,omitempty"`
	Options Options  `json:"options,omitempty"`
}

// Response structure for the Text classification endpoint
type TextClassificationResponse struct {
	// HFAPI returns a list of labels and their associated scores for
	// each input.
	Labels []*TextClassificationResponseLabel
}

// Used in TextClassificationResponse
type TextClassificationResponseLabel struct {
	// The label for the class (model specific)
	Name string `json:"label,omitempty"`

	// A float that represents how likely is that the text belongs in this class.
	Score float64 `json:"score,omitempty"`
}

func SendTextClassificationRequest(model string, request *TextClassificationRequest) ([]*TextClassificationResponse, error) {
	if request == nil {
		return nil, errors.New("nil TextClassificationRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	tclabels := make([][]*TextClassificationResponseLabel, len(request.Inputs))
	err = json.Unmarshal(respBody, &tclabels)
	if err != nil {
		return nil, err
	}

	tcresps := make([]*TextClassificationResponse, len(tclabels))
	for i := range tclabels {
		tcresps[i] = &TextClassificationResponse{
			Labels: tclabels[i],
		}
	}

	return tcresps, nil
}
