package hfapigo

import (
	"encoding/json"
)

const RecommendedSpeechRecongnitionModelEnglish = "facebook/wav2vec2-base-960h"

type SpeechRecognitionResponse struct {
	// The string that was recognized within the audio file.
	Text string `json:"text,omitempty"`
}

// SendSpeechRecognitionRequest takes a model string and a path to an audio file.
// It reads the file and sends a request to the speech recognition endpoint.
func SendSpeechRecognitionRequest(model, audioFile string) (*SpeechRecognitionResponse, error) {
	respBody, err := MakeHFAPIRequestWithMedia(model, audioFile)

	arresp := SpeechRecognitionResponse{}
	err = json.Unmarshal(respBody, &arresp)
	if err != nil {
		return nil, err
	}

	return &arresp, nil
}
