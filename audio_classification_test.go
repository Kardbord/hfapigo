package hfapigo_test

import (
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

func TestAudioClassificationRequest(t *testing.T) {
	acResps := []*hfapigo.AudioClassificationResponse{}
	var err error
	acResps, err = hfapigo.SendAudioClassificationRequest(hfapigo.RecommendedAudioClassificationModel, TestFilesDir+"/sample.flac")
	if err != nil {
		t.Fatal(err)
	}
	if len(acResps) == 0 {
		t.Fatal("Expected non-empty response")
	}

	for _, resp := range acResps {
		if resp == nil {
			t.Fatal("nil response received")
		}
		if resp.Score == 0.0 {
			t.Fatal("Expected non-zero score")
		}
		if resp.Label == "" {
			t.Fatal("Expected non-empty label")
		}
	}
}
