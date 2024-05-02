package hfapigo

import (
	"encoding/json"
	"errors"
)

const RecommendedTokenClassificationModel = "dslim/bert-base-NER-uncased"

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

	// Same as the simple strategy except entities cannot end up with different tags. Entities will use the tag of the first token when there is ambiguity.
	AggregationStrategyFirst AggregationStrategy = "first"

	// Same as the simple strategy except entities cannot end up with different tags. Scores are averaged across tokens and then the maximum label is applied.
	AggregationStrategyAverage AggregationStrategy = "average"

	// Same as the simple strategy except entities cannot end up with different tags. Entity will be the token with the maximum score.
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
	Entities []*TokenClassificationResponseEntity
}

type TokenClassificationResponseEntity struct {
	// The type for the entity being recognized (model specific).
	Label string `json:"entity_group,omitempty"`

	// How likely the entity was recognized.
	Score float64 `json:"score,omitempty"`

	// The string that was captured
	Entity string `json:"word,omitempty"`

	// The offset stringwise where the answer is located. Useful to disambiguate if Entity occurs multiple times.
	Start int `json:"start,omitempty"`

	// The offset stringwise where the answer is located. Useful to disambiguate if Entity occurs multiple times.
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

	strategy := request.Parameters.AggregationStrategy
	if strategy != nil && *strategy == AggregationStrategyNone {
		return handleNoAggregationResponse(respBody, request)
	}

	return handleAggregatedResponse(respBody, request)
}

func handleNoAggregationResponse(respBody []byte, request *TokenClassificationRequest) ([]*TokenClassificationResponse, error) {
	type EntityAlias struct {
		Label  string  `json:"entity,omitempty"` // different JSON key for non-aggregated responses
		Score  float64 `json:"score,omitempty"`
		Entity string  `json:"word,omitempty"`
		Start  int     `json:"start,omitempty"`
		End    int     `json:"end,omitempty"`
	}

	inputResps := make([][]*EntityAlias, len(request.Inputs))
	err := json.Unmarshal(respBody, &inputResps)
	if err != nil {
		return nil, err
	}

	tcresps := make([]*TokenClassificationResponse, len(request.Inputs))
	for i, iresp := range inputResps {
		tcresps[i] = &TokenClassificationResponse{}
		tcresps[i].Entities = make([]*TokenClassificationResponseEntity, len(iresp))
		for j, eg := range iresp {
			if eg == nil {
				return nil, errors.New("nil response encountered, this should never happen -- there is a bug in the hfapigo library")
			}
			var tcrg TokenClassificationResponseEntity = TokenClassificationResponseEntity(*eg)
			tcresps[i].Entities[j] = &tcrg
		}
	}

	return tcresps, nil
}

func handleAggregatedResponse(respBody []byte, request *TokenClassificationRequest) ([]*TokenClassificationResponse, error) {
	tcentities := make([][]*TokenClassificationResponseEntity, len(request.Inputs))
	err := json.Unmarshal(respBody, &tcentities)
	if err != nil {
		return nil, err
	}

	tcresps := make([]*TokenClassificationResponse, len(request.Inputs))
	for i := range tcentities {
		tcresps[i] = &TokenClassificationResponse{
			Entities: tcentities[i],
		}
	}

	return tcresps, nil
}
