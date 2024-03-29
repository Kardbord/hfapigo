package hfapigo

import (
	"encoding/json"
	"errors"
)

const RecommendedConversationalModel = "microsoft/DialoGPT-large"

// Request structure for the conversational endpoint
type ConversationalRequest struct {
	// (Required)
	Inputs ConverstationalInputs `json:"inputs,omitempty"`

	Parameters ConversationalParameters `json:"parameters,omitempty"`
	Options    Options                  `json:"options,omitempty"`
}

// Used with ConversationalRequest
type ConverstationalInputs struct {
	// (Required) The last input from the user in the conversation.
	Text string `json:"text,omitempty"`

	// A list of strings corresponding to the earlier replies from the model.
	GeneratedResponses []string `json:"generated_responses,omitempty"`

	// A list of strings corresponding to the earlier replies from the user.
	// Should be of the same length of GeneratedResponses.
	PastUserInputs []string `json:"past_user_inputs,omitempty"`
}

// Used with ConversationalRequest
type ConversationalParameters struct {
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
	RepetitionPenalty *float64 `json:"repetition_penalty,omitempty"`

	// (Default: None). Float (0-120.0). The amount of time in seconds that the query should take maximum.
	// Network can cause some overhead so it will be a soft limit.
	MaxTime *float64 `json:"maxtime,omitempty"`
}

func NewConversationalParameters() *ConversationalParameters {
	return &ConversationalParameters{}
}
func (c *ConversationalParameters) SetMinLength(minLength int) *ConversationalParameters {
	c.MinLength = &minLength
	return c
}
func (c *ConversationalParameters) SetMaxLength(maxLength int) *ConversationalParameters {
	c.MaxLength = &maxLength
	return c
}
func (c *ConversationalParameters) SetTopK(topK int) *ConversationalParameters {
	c.TopK = &topK
	return c
}
func (c *ConversationalParameters) SetTopP(topP float64) *ConversationalParameters {
	c.TopP = &topP
	return c
}
func (c *ConversationalParameters) SetTempurature(temperature float64) *ConversationalParameters {
	c.Temperature = &temperature
	return c
}
func (c *ConversationalParameters) SetRepetitionPenalty(penalty float64) *ConversationalParameters {
	c.RepetitionPenalty = &penalty
	return c
}
func (c *ConversationalParameters) SetMaxTime(maxTime float64) *ConversationalParameters {
	c.MaxTime = &maxTime
	return c
}

// Response structure for the conversational endpoint
type ConversationalResponse struct {
	// The answer of the model
	GeneratedText string `json:"generated_text,omitempty"`

	// A facility dictionary to send back for the next input (with the new user input addition).
	Conversation Conversation `json:"conversation,omitempty"`
}

// Used with ConversationalResponse
type Conversation struct {
	// The last outputs from the model in the conversation, after the model has run.
	GeneratedResponses []string `json:"generated_responses,omitempty"`

	// The last inputs from the user in the conversation, after the model has run.
	PastUserInputs []string `json:"past_user_inputs,omitempty"`
}

// Deprecated: HF's conversational endpoint seems to be under construction
// and slated to be either updated or replaced.
// TODO: Update or remove conversational support once it becomes
// clear what its replacement is.
func SendConversationalRequest(model string, request *ConversationalRequest) (*ConversationalResponse, error) {
	if request == nil {
		return nil, errors.New("nil ConversationalRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	cresp := ConversationalResponse{}
	err = json.Unmarshal(respBody, &cresp)
	if err != nil {
		return nil, err
	}

	return &cresp, nil
}
