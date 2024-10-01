package hfapigo_test

import (
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

func TestSpeechRecognitionRequest(t *testing.T) {
	arResp := &hfapigo.SpeechRecognitionResponse{}
	var err error
	arResp, err = hfapigo.SendSpeechRecognitionRequest(hfapigo.RecommendedSpeechRecongnitionModelEnglish, TestFilesDir+"/sample.flac")
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
