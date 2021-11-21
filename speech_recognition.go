package hfapigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const RecommendedSpeechRecongnitionModelEnglish = "facebook/wav2vec2-base-960h"

type SpeechRecognitionResponse struct {
	// The string that was recognized within the audio file.
	Text string `json:"text,omitempty"`
}

// SendSpeechRecognitionRequest takes a model string and a path to an audio file.
// It reads the file and sends a request to the speech recognition endpoint.
func SendSpeechRecognitionRequest(model, audioFile string) (*SpeechRecognitionResponse, error) {
	buf, err := os.ReadFile(audioFile)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, APIBaseURL+model, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	setAuthorizationHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = checkRespForError(respBody)
	if err != nil {
		return nil, err
	}

	arresp := SpeechRecognitionResponse{}
	err = json.Unmarshal(respBody, &arresp)
	if err != nil {
		return nil, err
	}

	return &arresp, nil
}
