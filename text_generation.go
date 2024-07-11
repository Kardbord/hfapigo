package hfapigo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	RecommendedTextGenerationModel = "microsoft/phi-2"
	TextGenerationGrammarTypeJSON  = "json"
	TextGenerationGrammarTypeRegex = "regex"
)

type TextGenerationRequest struct {
	// (Required) a string to be generated from
	Input      string                   `json:"inputs,omitempty"`
	Parameters TextGenerationParameters `json:"parameters,omitempty"`
	Options    Options                  `json:"options,omitempty"`
}

type TextGenerationParameters struct {
	BestOf              *int     `json:"best_of,omitempty"`
	DecoderInputDetails *bool    `json:"decoder_input_details,omitempty"`
	Details             *bool    `json:"details,omitempty"`
	DoSample            *bool    `json:"do_sample,omitempty"`
	FrequencyPenalty    *float64 `json:"frequency_penalty,omitempty"`
	Grammar             *string  `json:"grammar,omitempty"`
	MaxNewTokens        *int     `json:"max_new_tokens,omitempty"`
	RepetitionPenalty   *float64 `json:"repetition_penalty,omitempty"`
	ReturnFullText      *bool    `json:"return_full_text,omitempty"`
	Seed                *int64   `json:"seed,omitempty"`
	Stop                []string `json:"stop,omitempty"`
	Temperature         *float64 `json:"temperature,omitempty"`
	TopK                *int     `json:"top_k,omitempty"`
	TopNTokens          *int     `json:"top_n_tokens,omitempty"`
	TopP                *float64 `json:"top_p,omitempty"`
	Truncate            *int     `json:"truncate,omitempty"`
	TypicalP            *float64 `json:"typical_p,omitempty"`
	Watermark           *bool    `json:"watermark,omitempty"`
}

func NewTextGenerationParameters() *TextGenerationParameters {
	return &TextGenerationParameters{}
}
func (params *TextGenerationParameters) SetBestOf(bestOf int) *TextGenerationParameters {
	params.BestOf = &bestOf
	return params
}
func (params *TextGenerationParameters) SetDecoderInputDetails(decoderInputDetails bool) *TextGenerationParameters {
	params.DecoderInputDetails = &decoderInputDetails
	return params
}
func (params *TextGenerationParameters) SetDetails(details bool) *TextGenerationParameters {
	params.Details = &details
	return params
}
func (params *TextGenerationParameters) SetDoSample(doSample bool) *TextGenerationParameters {
	params.DoSample = &doSample
	return params
}
func (params *TextGenerationParameters) SetFrequencyPenalty(frequencyPenalty float64) *TextGenerationParameters {
	params.FrequencyPenalty = &frequencyPenalty
	return params
}
func (params *TextGenerationParameters) SetGrammar(grammar string) *TextGenerationParameters {
	params.Grammar = &grammar
	return params
}
func (params *TextGenerationParameters) SetMaxNewTokens(maxNewTokens int) *TextGenerationParameters {
	params.MaxNewTokens = &maxNewTokens
	return params
}
func (params *TextGenerationParameters) SetRepetitionPenalty(repetitionPenalty float64) *TextGenerationParameters {
	params.RepetitionPenalty = &repetitionPenalty
	return params
}
func (params *TextGenerationParameters) SetReturnFullText(returnFullText bool) *TextGenerationParameters {
	params.ReturnFullText = &returnFullText
	return params
}
func (params *TextGenerationParameters) SetSeed(seed int64) *TextGenerationParameters {
	params.Seed = &seed
	return params
}
func (params *TextGenerationParameters) SetStop(stop []string) *TextGenerationParameters {
	params.Stop = stop
	return params
}
func (params *TextGenerationParameters) SetTemperature(temperature float64) *TextGenerationParameters {
	params.Temperature = &temperature
	return params
}
func (params *TextGenerationParameters) SetTopK(topK int) *TextGenerationParameters {
	params.TopK = &topK
	return params
}
func (params *TextGenerationParameters) SetTopNTokens(topNTokens int) *TextGenerationParameters {
	params.TopNTokens = &topNTokens
	return params
}
func (params *TextGenerationParameters) SetTopP(topP float64) *TextGenerationParameters {
	params.TopP = &topP
	return params
}
func (params *TextGenerationParameters) SetTruncate(truncate int) *TextGenerationParameters {
	params.Truncate = &truncate
	return params
}
func (params *TextGenerationParameters) SetTypicalP(typicalP float64) *TextGenerationParameters {
	params.TypicalP = &typicalP
	return params
}
func (params *TextGenerationParameters) SetWatermark(watermark bool) *TextGenerationParameters {
	params.Watermark = &watermark
	return params
}
func (params *TextGenerationParameters) SetRepetitionPenaly(penalty float64) *TextGenerationParameters {
	params.RepetitionPenalty = &penalty
	return params
}

type TextGenerationResponse struct {
	GeneratedText string                        `json:"generated_text,omitempty"`
	Details       TextGenerationResponseDetails `json:"details,omitempty"`
}

type TextGenerationResponseDetails struct {
	BestOfSequences []*TextGenerationBestOfSequence `json:"best_of_sequences,omitempty"`
	FinishReason    string                          `json:"finish_reason,omitempty"`
	GeneratedTokens int                             `json:"generated_tokens,omitempty"`
	Prefill         []*TextGenerationPrefillToken   `json:"prefill,omitempty"`
	Seed            int64                           `json:"seed,omitempty"`
	Tokens          []*TextGenerationToken          `json:"tokens,omitempty"`
	TopTokens       []*TextGenerationToken          `json:"top_tokens,omitempty"`
}

type TextGenerationBestOfSequence struct {
	FinishReason    string                        `json:"finish_reason,omitempty"`
	GeneratedText   string                        `json:"generated_text,omitempty"`
	GeneratedTokens int                           `json:"generated_tokens,omitempty"`
	Prefill         []*TextGenerationPrefillToken `json:"prefill,omitempty"`
	Seed            int64                         `json:"seed,omitempty"`
	Tokens          []*TextGenerationToken        `json:"tokens,omitempty"`
	TopTokens       [][]*TextGenerationToken      `json:"top_tokens,omitempty"`
}

type TextGenerationPrefillToken struct {
	ID      int     `json:"id,omitempty"`
	LogProb float64 `json:"logprob,omitempty"`
	Text    string  `json:"text,omitempty"`
}

type TextGenerationToken struct {
	TextGenerationPrefillToken
	Special bool `json:"special,omitempty"`
}

func SendTextGenerationRequest(model string, request *TextGenerationRequest) ([]*TextGenerationResponse, error) {
	if request == nil {
		return nil, errors.New("nil TextGenerationRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	tgresps := make([]*TextGenerationResponse, 1)
	err = json.Unmarshal(respBody, &tgresps)
	if err != nil {
		return nil, err
	}
	if len(tgresps) < 1 {
		return nil, fmt.Errorf("expected at least 1 response, got none; response=%s", string(respBody))
	}

	return tgresps, nil
}
