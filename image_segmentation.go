package hfapigo

import "encoding/json"

const RecommendedImageSegmentationModel = "facebook/detr-resnet-50-panoptic"

type ImageSegmentationResponse struct {
	// The label for the class (model specific) of a segment.
	Label string `json:"label,omitempty"`

	// A float that represents how likely it is that the segment belongs to the given class.
	Score float64 `json:"score,omitempty"`

	// A str (base64 str of a single channel black-and-white img) representing the mask of a segment.
	Mask string `json:"mask,omitempty"`
}

func SendImageSegmentationRequest(model, imageFile string) ([]*ImageSegmentationResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, imageFile)
	if err != nil {
		return nil, err
	}

	resps := []*ImageSegmentationResponse{}
	err = json.Unmarshal(respBody, &resps)
	if err != nil {
		return nil, err
	}

	return resps, nil
}