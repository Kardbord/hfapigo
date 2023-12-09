package hfapigo_test

import (
	"testing"
	"time"

	"github.com/Kardbord/hfapigo"
)

func TestImageToText(t *testing.T) {
	const retries = 10

	resps := []*hfapigo.ImageToTextResponse{}
	var err error
	for i := 0; i < retries; i++ {
		resps, err = hfapigo.SendImageToTextRequest(hfapigo.RecommendedImageToTextModel, TestFilesDir+"/test-image.png")
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}

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
