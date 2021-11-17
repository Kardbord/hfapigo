package hfapigo

import (
	"encoding/json"
	"errors"
)

const RecommendedTableQuestionAnsweringModel = "google/tapas-base-finetuned-wtq"

// Request structure for table question answering model
type TableQuestionAnsweringRequest struct {
	Inputs  TableQuestionAnsweringInputs `json:"inputs,omitempty"`
	Options Options                      `json:"options,omitempty"`
}

type TableQuestionAnsweringInputs struct {
	// (Required) The query in plain text that you want to ask the table
	Query string `json:"query,omitempty"`

	// (Required) A table of data represented as a dict of list where entries
	// are headers and the lists are all the values, all lists must
	// have the same size.
	Table map[string][]string `json:"table,omitempty"`
}

// Response structure for table question answering model
type TableQuestionAnsweringResponse struct {
	// The plaintext answer
	Answer string `json:"answer,omitempty"`

	// A list of coordinates of the cells references in the answer
	Coordinates [][]int `json:"coordinates,omitempty"`

	// A list of coordinates of the cells contents
	Cells []string `json:"cells,omitempty"`

	// The aggregator used to get the answer
	Aggregator string `json:"aggregator,omitempty"`
}

func SendTableQuestionAnsweringRequest(model string, request *TableQuestionAnsweringRequest) (*TableQuestionAnsweringResponse, error) {
	if request == nil {
		return nil, errors.New("nil tableQuestionAnsweringRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	tqaResp := TableQuestionAnsweringResponse{}
	err = json.Unmarshal(respBody, &tqaResp)
	if err != nil {
		return nil, err
	}

	return &tqaResp, nil
}
