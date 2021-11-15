package hfapigo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const RecommmendedSummarizationModel = "facebook/bart-large-cnn"

// Request structure for the summarization endpoint
type SummarizationRequest struct {
	// Strings to be summarized
	Inputs     []string                `json:"inputs,omitempty"`
	Parameters SummarizationParameters `json:"parameters,omitempty"`
	Options    Options                 `json:"options,omitempty"`
}

// Used with SummarizationRequest
type SummarizationParameters struct {
	// (Default: None). Integer to define the minimum length in tokens of the output summary.
	MinLength int `json:"min_length,omitempty"`

	// (Default: None). Integer to define the maximum length in tokens of the output summary.
	MaxLength int `json:"max_length,omitempty"`

	// (Default: None). Integer to define the top tokens considered within the sample operation to create
	// new text.
	TopK int `json:"top_k,omitempty"`

	// (Default: None). Float to define the tokens that are within the sample` operation of text generation.
	// Add tokens in the sample for more probable to least probable until the sum of the probabilities is
	// greater than top_p.
	TopP float64 `json:"top_p,omitempty"`

	// (Default: 1.0). Float (0.0-100.0). The temperature of the sampling operation. 1 means regular sampling,
	// 0 mens top_k=1, 100.0 is getting closer to uniform probability.
	Temperature *float64 `json:"temperature,omitempty"`

	// (Default: None). Float (0.0-100.0). The more a token is used within generation the more it is penalized
	// to not be picked in successive generation passes.
	RepetitionPenalty float64 `json:"repetitionpenalty,omitempty"`

	// (Default: None). Float (0-120.0). The amount of time in seconds that the query should take maximum.
	// Network can cause some overhead so it will be a soft limit.
	MaxTime float64 `json:"maxtime,omitempty"`

	// This option is used in the API example, but not documented.
	// Including it anyway until they either remove it or document it.
	DoSample bool `json:"do_sample"`
}

func NewSummarizationParameters() *SummarizationParameters {
	return &SummarizationParameters{}
}

func (sp *SummarizationParameters) SetTempurature(temperature float64) *SummarizationParameters {
	sp.Temperature = &temperature
	return sp
}

// Response structure for the summarization endpoint
type SummarizationResponse struct {
	// The summarized input string
	SummaryText string `json:"summary_text,omitempty"`
}

func SendSummarizationRequest(model string, request *SummarizationRequest) ([]*SummarizationResponse, error) {
	endpoint := APIBaseURL + model
	if request == nil {
		return nil, errors.New("nil SummarizationRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := BuildHFAPIRequest(jsonBuf, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sresps := make([]*SummarizationResponse, len(request.Inputs))
	err = json.Unmarshal(respBody, &sresps)
	if err != nil {
		return nil, errors.New(string(respBody))
	}

	return sresps, nil
}
