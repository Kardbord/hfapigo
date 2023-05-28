package hfapigo

import "encoding/json"

const RecommendedImageToTextModel = "nlpconnect/vit-gpt2-image-captioning"

type ImageToTextResponse struct {
	// The generated caption
	GeneratedText string `json:"generated_text"`
}

func SendImageToTextRequest(model, imageFile string) ([]*ImageToTextResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, imageFile)
	if err != nil {
		return nil, err
	}

	resps := []*ImageToTextResponse{}
	err = json.Unmarshal(respBody, &resps)
	if err != nil {
		return nil, err
	}

	return resps, nil
}
