package hfapigo

import (
	"encoding/json"
)

const RecommendedObjectDetectionModel = "facebook/detr-resnet-50"

type ObjectDetectionResponse struct {
	// The label for the class (model specific) of a detected object.
	Label string `json:"label,omitempty"`

	// A float that represents how likely it is that the detected object belongs to the given class.
	Score float64 `json:"score,omitempty"`

	// Bounding box of the detected object
	Box ObjectBox
}

type ObjectBox struct {
	XMin int `json:"xmin,omitempty"`
	YMin int `json:"ymin,omitempty"`
	XMax int `json:"xmax,omitempty"`
	YMax int `json:"ymax,omitempty"`
}

func SendObjectDetectionRequest(model, imageFile string) ([]*ObjectDetectionResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, imageFile)

	resps := []*ObjectDetectionResponse{}
	err = json.Unmarshal(respBody, &resps)
	if err != nil {
		return nil, err
	}

	return resps, nil
}
