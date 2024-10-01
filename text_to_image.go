package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"

	_ "image/jpeg"
	_ "image/png"
)

const RecommendedTextToImageModel = "stable-diffusion-v1-5/stable-diffusion-v1-5"

// Request structure for text-to-image model
type TextToImageRequest struct {
	// The prompt or prompts to guide the image generation.
	Inputs     string                       `json:"inputs,omitempty"`
	Options    Options                      `json:"options,omitempty"`
	Parameters TextToImageRequestParameters `json:"parameters,omitempty"`
}

type TextToImageRequestParameters struct {
	// The prompt or prompts not to guide the image generation.
	// Ignored when not using guidance (i.e., ignored if guidance_scale is less than 1).
	NegativePrompt string `json:"negative_prompt,omitempty"`
	// The height in pixels of the generated image.
	Height int64 `json:"height,omitempty"`
	// The width in pixels of the generated image.
	Width int64 `json:"width,omitempty"`
	// The number of denoising steps. More denoising steps usually lead to a higher quality
	// image at the expense of slower inference. Defaults to 50.
	NumInferenceSteps int64 `json:"num_inference_steps,omitempty"`
	// Higher guidance scale encourages to generate images that are closely linked to the text
	// input, usually at the expense of lower image quality. Defaults to 7.5.
	GuidanceScale float64 `json:"guidance_scale,omitempty"`
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
