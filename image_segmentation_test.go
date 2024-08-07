package hfapigo_test

import (
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

func TestImageSegmentationRequest(t *testing.T) {
	resps, err := hfapigo.SendImageSegmentationRequest(hfapigo.RecommendedImageSegmentationModel, TestFilesDir+"/test-image.png")
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) == 0 {
		t.Fatal("expected non-empty response")
	}
	for _, resp := range resps {
		if resp == nil {
			t.Fatal("nil response received")
		}
		if resp.Score == 0.0 {
			t.Fatal("expected non-zero score")
		}
		if resp.Label == "" {
			t.Fatal("expected non-empty label")
		}
		if resp.Mask == "" {
			t.Fatal("expected non-empty mask")
		}
	}
}
