package hfapigo

import (
	"encoding/json"
	"errors"
)

const RecommendedQuestionAnsweringModel = "bert-large-uncased-whole-word-masking-finetuned-squad"

// Request structure for question answering model
type QuestionAnsweringRequest struct {
	// (Required)
	Inputs  QuestionAnsweringInputs `json:"inputs,omitempty"`
	Options Options                 `json:"options,omitempty"`
}

type QuestionAnsweringInputs struct {
	// (Required) The question as a string that has an answer within Context.
	Question string `json:"question,omitempty"`

	// (Required) A string that contains the answer to the question
	Context string `json:"context,omitempty"`
}

// Response structure for question answering model
type QuestionAnsweringResponse struct {
	// A string that’s the answer within the Context text.
	Answer string `json:"answer,omitempty"`

	// A float that represents how likely that the answer is correct.
	Score float64 `json:"score,omitempty"`

	// The string index of the start of the answer within Context.
	Start int `json:"start,omitempty"`

	// The string index of the stop of the answer within Context.
	End int `json:"end,omitempty"`
}

func SendQuestionAnsweringRequest(model string, request *QuestionAnsweringRequest) (*QuestionAnsweringResponse, error) {
	if request == nil {
		return nil, errors.New("nil QuestionAnsweringRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	qaResp := QuestionAnsweringResponse{}
	err = json.Unmarshal(respBody, &qaResp)
	if err != nil {
		return nil, err
	}

	return &qaResp, nil
}
