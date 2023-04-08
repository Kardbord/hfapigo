package hfapigo

import "encoding/json"

const RecommendedImageClassificationModel = "google/vit-base-patch16-224"

type ImageClassificationResponse struct {
	// The label for the class (model specific)
	Label string `json:"label,omitempty"`

	// A float that represents how likely it is that the image file belongs to this class.
	Score float64 `json:"score,omitempty"`
}

func SendImageClassificationRequest(model, imageFile string) ([]*ImageClassificationResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, imageFile)
	if err != nil {
		return nil, err
	}

	resps := []*ImageClassificationResponse{}
	err = json.Unmarshal(respBody, &resps)
	if err != nil {
		return nil, err
	}

	return resps, nil
}
