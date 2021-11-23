package hfapigo

import (
	"encoding/json"
	"errors"
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
	MinLength *int `json:"min_length,omitempty"`

	// (Default: None). Integer to define the maximum length in tokens of the output summary.
	MaxLength *int `json:"max_length,omitempty"`

	// (Default: None). Integer to define the top tokens considered within the sample operation to create
	// new text.
	TopK *int `json:"top_k,omitempty"`

	// (Default: None). Float to define the tokens that are within the sample` operation of text generation.
	// Add tokens in the sample for more probable to least probable until the sum of the probabilities is
	// greater than top_p.
	TopP *float64 `json:"top_p,omitempty"`

	// (Default: 1.0). Float (0.0-100.0). The temperature of the sampling operation. 1 means regular sampling,
	// 0 mens top_k=1, 100.0 is getting closer to uniform probability.
	Temperature *float64 `json:"temperature,omitempty"`

	// (Default: None). Float (0.0-100.0). The more a token is used within generation the more it is penalized
	// to not be picked in successive generation passes.
	RepetitionPenalty *float64 `json:"repetitionpenalty,omitempty"`

	// (Default: None). Float (0-120.0). The amount of time in seconds that the query should take maximum.
	// Network can cause some overhead so it will be a soft limit.
	MaxTime *float64 `json:"maxtime,omitempty"`
}

func NewSummarizationParameters() *SummarizationParameters {
	return &SummarizationParameters{}
}
func (sp *SummarizationParameters) SetMinLength(minLength int) *SummarizationParameters {
	sp.MinLength = &minLength
	return sp
}
func (sp *SummarizationParameters) SetMaxLength(maxLength int) *SummarizationParameters {
	sp.MaxLength = &maxLength
	return sp
}
func (sp *SummarizationParameters) SetTopK(topK int) *SummarizationParameters {
	sp.TopK = &topK
	return sp
}
func (sp *SummarizationParameters) SetTopP(topP float64) *SummarizationParameters {
	sp.TopP = &topP
	return sp
}
func (sp *SummarizationParameters) SetTempurature(temperature float64) *SummarizationParameters {
	sp.Temperature = &temperature
	return sp
}
func (sp *SummarizationParameters) SetRepetitionPenalty(penalty float64) *SummarizationParameters {
	sp.RepetitionPenalty = &penalty
	return sp
}
func (sp *SummarizationParameters) SetMaxTime(maxTime float64) *SummarizationParameters {
	sp.MaxTime = &maxTime
	return sp
}

// Response structure for the summarization endpoint
type SummarizationResponse struct {
	// The summarized input string
	SummaryText string `json:"summary_text,omitempty"`
}

func SendSummarizationRequest(model string, request *SummarizationRequest) ([]*SummarizationResponse, error) {
	if request == nil {
		return nil, errors.New("nil SummarizationRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	sresps := make([]*SummarizationResponse, len(request.Inputs))
	err = json.Unmarshal(respBody, &sresps)
	if err != nil {
		return nil, err
	}

	return sresps, nil
}
