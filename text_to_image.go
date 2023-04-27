package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"

	_ "image/jpeg"
	_ "image/png"
)

const RecommendedTextToImageModel = "runwayml/stable-diffusion-v1-5"

// Request structure for text-to-image model
type TextToImageRequest struct {
	Inputs  string  `json:"inputs,omitempty"`
	Options Options `json:"options,omitempty"`
}

// Send a TextToImageRequest. If successful, returns the generated image object, format name, and nil.
// If unsuccessful, returns nil, "", and an error.
func SendTextToImageRequest(model string, request *TextToImageRequest) (image.Image, string, error) {
	if request == nil {
		return nil, "", errors.New("nil TextToImageRequest")
	}

	jsonBuf, err := json.Marshal(request)
	if err != nil {
		return nil, "", err
	}

	respBody, err := MakeHFAPIRequest(jsonBuf, model)
	if err != nil {
		return nil, "", err
	}

	return image.Decode(bytes.NewReader(respBody))
}
