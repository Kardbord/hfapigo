package hfapigo_test

import (
	"testing"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
)

func TestAudioClassificationRequest(t *testing.T) {
	const retries = 10

	acResps := []*hfapigo.AudioClassificationResponse{}
	var err error
	for i := 0; i < retries; i++ {
		acResps, err = hfapigo.SendAudioClassificationRequest(hfapigo.RecommendedAudioClassificationModel, TestFilesDir+"/sample.flac")
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
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
