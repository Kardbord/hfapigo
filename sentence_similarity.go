package hfapigo

import (
	"encoding/json"
	"errors"
)

const (
	RecommendedSentenceSimilarityModel = "sentence-transformers/all-MiniLM-L6-v2"
)

// Request structure for the Sentence Similarity endpoint.
type SentenceSimilarityRequest struct {
	// (Required) Inputs for the request.
	Inputs struct {
		// (Required) The string that you wish to compare the other strings with.
		// This can be a phrase, sentence, or longer passage, depending on the
		// model being used.
		SourceSentence string `json:"source_sentence,omitempty"`

		// A list of strings which will be compared against the source_sentence.
		Sentences []string `json:"sentences,omitempty"`
	} `json:"inputs,omitempty"`

	Options Options `json:"options,omitempty"`
}

// Response structure from the Sentence Similarity endpoint.
// The return value is a list of similarity scores, given as floats.
// Each list entry corresponds to the Inputs.Sentences list entry
// of the same index.
type SentenceSimilarityResponse []float64

func SendSentenceSimilarityRequest(model string, request *SentenceSimilarityRequest) (*SentenceSimilarityResponse, error) {
	if request == nil {
		return nil, errors.New("nil SentenceSimilarityRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	resps := make(SentenceSimilarityResponse, len(request.Inputs.Sentences))
	err = json.Unmarshal(respBody, &resps)
	if err != nil {
		return nil, err
	}

	return &resps, nil
}
