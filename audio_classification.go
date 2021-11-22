package hfapigo

import "encoding/json"

const RecommendedAudioClassificationModel = "superb/hubert-large-superb-er"

// Response structure for audio classification endpoint
type AudioClassificationResponse struct {
	Score float64 `json:"score,omitempty"`
	Label string  `json:"label,omitempty"`
}

func SendAudioClassificationRequest(model, audioFile string) ([]*AudioClassificationResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, audioFile)

	acresp := []*AudioClassificationResponse{}
	err = json.Unmarshal(respBody, &acresp)
	if err != nil {
		return nil, err
	}

	return acresp, nil
}
