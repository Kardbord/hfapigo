package hfapigo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	RecommendedRussianToEnglishModel = "Helsinki-NLP/opus-mt-ru-en"
)

// Request structure for the Translation endpoint
type TranslationRequest struct {
	// (Required) a string to be translated in the original languages
	Input []string `json:"inputs,omitempty"`

	Options Options `json:"options,omitempty"`
}

// Response structure from the Translation endpoint
type TranslationResponse struct {
	// The translated Input string
	TranslationText string `json:"translation_text,omitempty"`
}

func SendTranslationRequest(model string, request *TranslationRequest) ([]*TranslationResponse, error) {
	endpoint := APIBaseURL + model
	if request == nil {
		return nil, errors.New("nil TranslationRequest")
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

	tresps := make([]*TranslationResponse, 1)
	err = json.Unmarshal(respBody, &tresps)
	if err != nil {
		return nil, errors.New(string(respBody))
	}

	return tresps, nil
}
