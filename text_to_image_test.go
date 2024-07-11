package hfapigo_test

import (
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

func TestTextToImage(t *testing.T) {
	{ // Test valid request
		img, fmt, err := hfapigo.SendTextToImageRequest(hfapigo.RecommendedTextToImageModel, &hfapigo.TextToImageRequest{
			Inputs:  "A dog and a cat sleeping adorably.",
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		if fmt == "" {
			t.Fatal("empty encoding returned")
		}
		if img == nil {
			t.Fatal("nil image returned")
		}
	}

	{ // Test invalid request
		img, fmt, err := hfapigo.SendTextToImageRequest("not-a-model", &hfapigo.TextToImageRequest{
			Inputs:  "A dog and a cat sleeping adorably.",
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		if err == nil {
			t.Fatal("expected an error")
		}
		if fmt != "" {
			t.Fatal("expected an empty encoding string")
		}
		if img != nil {
			t.Fatal("expected a nil image")
		}
	}
}
