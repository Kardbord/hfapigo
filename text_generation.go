package hfapigo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const RecommendedTextGenerationModel = "gpt2-large"

type TextGenerationRequest struct {
	// (Required) a string to be generated from
	Inputs     []string                 `json:"inputs,omitempty"`
	Parameters TextGenerationParameters `json:"parameters,omitempty"`
	Options    Options                  `json:"options,omitempty"`
}

type TextGenerationParameters struct {
	// (Default: None). Integer to define the top tokens considered within the sample operation to create new text.
	TopK *int `json:"top_k,omitempty"`

	// (Default: None). Float to define the tokens that are within the sample` operation of text generation. Add
	// tokens in the sample for more probable to least probable until the sum of the probabilities is greater
	// than top_p.
	TopP *float64 `json:"top_p,omitempty"`

	// (Default: 1.0). Float (0.0-100.0). The temperature of the sampling operation. 1 means regular sampling,
	// 0 means top_k=1, 100.0 is getting closer to uniform probability.
	Temperature *float64 `json:"temperature,omitempty"`

	// (Default: None). Float (0.0-100.0). The more a token is used within generation the more it is penalized
	// to not be picked in successive generation passes.
	RepetitionPenalty *float64 `json:"repetition_penalty,omitempty"`

	// (Default: None). Int (0-250). The amount of new tokens to be generated, this does not include the input
	// length it is a estimate of the size of generated text you want. Each new tokens slows down the request,
	// so look for balance between response times and length of text generated.
	MaxNewTokens *int `json:"max_new_tokens,omitempty"`

	// (Default: None). Float (0-120.0). The amount of time in seconds that the query should take maximum.
	// Network can cause some overhead so it will be a soft limit. Use that in combination with max_new_tokens
	// for best results.
	MaxTime *float64 `json:"max_time,omitempty"`

	// (Default: True). Bool. If set to False, the return results will not contain the original query making it
	// easier for prompting.
	ReturnFullText *bool `json:"return_full_text,omitempty"`

	// (Default: 1). Integer. The number of proposition you want to be returned.
	NumReturnSequences *int `json:"num_return_sequences,omitempty"`
}

func NewTextGenerationParameters() *TextGenerationParameters {
	return &TextGenerationParameters{}
}
func (params *TextGenerationParameters) SetTopK(topK int) *TextGenerationParameters {
	params.TopK = &topK
	return params
}
func (params *TextGenerationParameters) SetTopP(topP float64) *TextGenerationParameters {
	params.TopP = &topP
	return params
}
func (params *TextGenerationParameters) SetTempurature(temp float64) *TextGenerationParameters {
	params.Temperature = &temp
	return params
}
func (params *TextGenerationParameters) SetRepetitionPenaly(penalty float64) *TextGenerationParameters {
	params.RepetitionPenalty = &penalty
	return params
}
func (params *TextGenerationParameters) SetMaxNewTokens(maxNewTokens int) *TextGenerationParameters {
	params.MaxNewTokens = &maxNewTokens
	return params
}
func (params *TextGenerationParameters) SetMaxTime(maxTime float64) *TextGenerationParameters {
	params.MaxTime = &maxTime
	return params
}
func (params *TextGenerationParameters) SetReturnFullText(returnFullText bool) *TextGenerationParameters {
	params.ReturnFullText = &returnFullText
	return params
}
func (params *TextGenerationParameters) SetNumReturnSequences(numReturnSequences int) *TextGenerationParameters {
	params.NumReturnSequences = &numReturnSequences
	return params
}

type TextGenerationResponse struct {
	// A list of generated texts. The length of this list is the value of
	// NumReturnSequences in the request.
	GeneratedTexts []string
}

type textGenerationResponseSequence struct {
	GeneratedText string `json:"generated_text,omitempty"`
}

func (tgs textGenerationResponseSequence) String() string {
	return tgs.GeneratedText
}

func SendTextGenerationRequest(model string, request *TextGenerationRequest) ([]*TextGenerationResponse, error) {
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

	tgrespsRaw := make([][]*textGenerationResponseSequence, len(request.Inputs))
	err = json.Unmarshal(respBody, &tgrespsRaw)
	if err != nil {
		return nil, err
	}
	if len(tgrespsRaw) != len(request.Inputs) {
		return nil, fmt.Errorf("expected %d responses, got %d; response=%s", len(request.Inputs), len(tgrespsRaw), string(respBody))
	}

	tgresps := make([]*TextGenerationResponse, len(request.Inputs))
	for i := range tgrespsRaw {
		tgresps[i] = &TextGenerationResponse{}
		for _, t := range tgrespsRaw[i] {
			tgresps[i].GeneratedTexts = append(tgresps[i].GeneratedTexts, t.GeneratedText)
		}
	}

	return tgresps, nil
}
