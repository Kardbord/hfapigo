package hfapigo_test

import (
	"testing"
	"time"

	"github.com/Kardbord/hfapigo/v3"
)

func TestObjectDetectionRequest(t *testing.T) {
	const retries = 10

	resps := []*hfapigo.ObjectDetectionResponse{}
	var err error
	for i := 0; i < retries; i++ {
		resps, err = hfapigo.SendObjectDetectionRequest(hfapigo.RecommendedObjectDetectionModel, TestFilesDir+"/test-image.png")
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
		if resp.Score == 0.0 {
			t.Fatal("Expected non-zero score")
		}
		if resp.Label == "" {
			t.Fatal("Expected non-empty label")
		}
		if equal(resp.Box.XMin, resp.Box.XMax, resp.Box.YMin, resp.Box.YMax) {
			t.Fatal("expected non-equal coordinates")
		}
	}
}

func equal(nums ...int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[0] {
			return false
		}
	}
	return true
}
