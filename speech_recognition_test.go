package hfapigo_test

import (
	"testing"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func TestSpeechRecognitionRequest(t *testing.T) {
	const retries = 10

	arResp := &hfapigo.SpeechRecognitionResponse{}
	var err error
	for i := 0; i < retries; i++ {
		arResp, err = hfapigo.SendSpeechRecognitionRequest(hfapigo.RecommendedSpeechRecongnitionModelEnglish, TestFilesDir+"/sample.flac")
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	if arResp == nil {
		t.Fatal("Expected non-nil response")
	}
	if arResp.Text == "" {
		t.Fatal("Expected non-empty Text response")
	}
}
