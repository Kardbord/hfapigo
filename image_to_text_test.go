package hfapigo_test

import (
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

func TestImageToText(t *testing.T) {

	resps := []*hfapigo.ImageToTextResponse{}
	var err error
	resps, err = hfapigo.SendImageToTextRequest(hfapigo.RecommendedImageToTextModel, TestFilesDir+"/test-image.png")
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) == 0 {
		t.Fatal("Expected non-empty response")
	}
	for _, resp := range resps {
		if resp == nil {
			t.Fatal("nil response received")
		}
		if resp.GeneratedText == "" {
			t.Fatal("Expected non-empty caption")
		}
	}
}
