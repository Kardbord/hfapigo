package hfapigo

import (
	"encoding/json"
	"errors"
	"regexp"
)

const RecommendedTokenClassificationModel = "dbmdz/bert-large-cased-finetuned-conll03-english"

// Request structure for the token classification endpoint
type TokenClassificationRequest struct {
	// (Required) strings to be classified
	Inputs     []string                      `json:"inputs,omitempty"`
	Parameters TokenClassificationParameters `json:"parameters,omitempty"`
	Options    Options                       `json:"options,omitempty"`
}

type AggregationStrategy string

const (
	// Every token gets classified without further aggregation.
	AggregationStrategyNone AggregationStrategy = "none"

	// Entities are grouped according to the default schema (B-, I- tags get merged when the tag is similar).
	AggregationStrategySimple AggregationStrategy = "simple"

	// Same as the simple strategy except words cannot end up with different tags. Words will use the tag of the first token when there is ambiguity.
	AggregationStrategyFirst AggregationStrategy = "first"

	// Same as the simple strategy except words cannot end up with different tags. Scores are averaged across tokens and then the maximum label is applied.
	AggregationStrategyAverage AggregationStrategy = "average"

	// Same as the simple strategy except words cannot end up with different tags. Word entity will be the token with the maximum score.
	AggregationStrategyMax AggregationStrategy = "max"
)

type TokenClassificationParameters struct {
	// (Default: simple)
	AggregationStrategy *AggregationStrategy `json:"aggregation_strategy,omitempty"`
}

func NewTokenClassificationParameters() *TokenClassificationParameters {
	return &TokenClassificationParameters{}
}

func (params *TokenClassificationParameters) SetAggregationStrategy(aggregationStrategy AggregationStrategy) *TokenClassificationParameters {
	params.AggregationStrategy = &aggregationStrategy
	return params
}

// Response structure for the token classification endpoint
type TokenClassificationResponse struct {
	EntityGroups []*TokenClassificationResponseEntityGroup
}

type TokenClassificationResponseEntityGroup struct {
	// The type for the entity being recognized (model specific).
	Name string `json:"entity_group,omitempty"`

	// How likely the entity was recognized.
	Score float64 `json:"score,omitempty"`

	// The string that was captured
	Word string `json:"word,omitempty"`

	// The offset stringwise where the answer is located. Useful to disambiguate if Word occurs multiple times.
	Start int `json:"start,omitempty"`

	// The offset stringwise where the answer is located. Useful to disambiguate if word occurs multiple times.
	End int `json:"end,omitempty"`
}

func SendTokenClassificationRequest(model string, request *TokenClassificationRequest) ([]*TokenClassificationResponse, error) {
	if request == nil {
		return nil, errors.New("nil TokenClassificationRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, err
	}

	// This is a hack that has to exist because the HF API does not use a consistent key for the EntityGroup
	// in the JSON response. Sometimes, even within the context of a single request, it is "entity_group"
	// as specified in the docs, and sometimes it is just "entity".
	entityRegexp, err := regexp.Compile(`"entity":`)
	if err != nil {
		return nil, err
	}
	respBody = entityRegexp.ReplaceAll(respBody, []byte(`"entity_group":`))

	tcgroups := make([][]*TokenClassificationResponseEntityGroup, len(request.Inputs))
	err = json.Unmarshal(respBody, &tcgroups)
	if err != nil {
		return nil, err
	}

	tcresps := make([]*TokenClassificationResponse, len(request.Inputs))
	for i := range tcgroups {
		tcresps[i] = &TokenClassificationResponse{
			EntityGroups: tcgroups[i],
		}
	}

	return tcresps, nil
}
