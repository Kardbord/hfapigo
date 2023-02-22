package hfapigo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const RecommendedFillMaskModel = "bert-base-uncased"

// Request structure for the Fill Mask endpoint
type FillMaskRequest struct {
	// (Required) a string to be filled from, must contain the [MASK] token (check model card for exact name of the mask)
	Inputs  []string `json:"inputs,omitempty"`
	Options Options  `json:"options,omitempty"`
}

// Response structure for the Fill Mask endpoint
type FillMaskResponse struct {
	Masks []*FillMaskResponseEntry
}

// Used in the FillMaskResponse struct
type FillMaskResponseEntry struct {
	// The actual sequence of tokens that ran against the model (may contain special tokens)
	Sequence string `json:"sequence,omitempty"`

	// The probability for this token.
	Score float64 `json:"score,omitempty"`

	// The id of the token
	TokenID int `json:"token,omitempty"`

	// The string representation of the token
	TokenStr string `json:"token_str,omitempty"`
}

func SendFillMaskRequest(model string, request *FillMaskRequest) ([]*FillMaskResponse, error) {
	if request == nil {
		return nil, errors.New("nil FillMaskRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	fmResps := []*FillMaskResponse{}
	{
		// Multi-input or multi-mask implies response is a list of lists of dicts.
		rawResps := make([][]*FillMaskResponseEntry, len(request.Inputs))
		err = json.Unmarshal(respBody, &rawResps)
		if err == nil {
			for i := range rawResps {
				fmResps = append(fmResps, &FillMaskResponse{
					Masks: rawResps[i],
				})
			}
			return fmResps, nil
		}
	}
	{
		// Single input, single mask implies response is a list of dicts.
		rawResps := make([]*FillMaskResponseEntry, len(request.Inputs))
		err2 := json.Unmarshal(respBody, &rawResps)
		if err2 != nil {
			err = fmt.Errorf("%s; %w", err, err2)
			return nil, err
		}
		fmResps = append(fmResps, &FillMaskResponse{
			Masks: rawResps,
		})
		return fmResps, nil
	}
}
