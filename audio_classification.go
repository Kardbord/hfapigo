package hfapigo

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

const RecommendedAudioClassificationModel = "superb/hubert-large-superb-er"

type AudioClassificationRequestParameters struct {
	// Possible values: sigmoid, softmax, none.
	FunctionToApply string `json:"function_to_apply,omitempty"`

	// When specified, limits the output to the top K most probable classes.
	TopK uint32 `json:"top_k,omitempty"`
}

type AudioClassificationRequest struct {
	// Path to an audio file to send with the request
	InputFile string `json:"-"`

	// The input audio data as a base64-encoded string.
	// Automatically populated from InputFile if not provided.
	Input *string `json:"inputs,omitempty"`

	// Optional input parameters
	Parameters *AudioClassificationRequestParameters `json:"parameters,omitempty"`
}

// Response structure for audio classification endpoint
type AudioClassificationResponse struct {
	Score float64 `json:"score,omitempty"`
	Label string  `json:"label,omitempty"`
}

func SendAudioClassificationRequest(model string, request *AudioClassificationRequest) ([]*AudioClassificationResponse, error) {
	if request.Input == nil || *request.Input == "" {
		data, err := os.ReadFile(request.InputFile)
		if err != nil {
			return nil, err
		}
		input := base64.StdEncoding.EncodeToString(data)
		request.Input = &input
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	acresp := []*AudioClassificationResponse{}
	err = json.Unmarshal(respBody, &acresp)
	if err != nil {
		return nil, err
	}

	return acresp, nil
}
