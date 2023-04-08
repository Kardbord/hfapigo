package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/hfapigo"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

const inputImg = "./test-image.png"

func main() {
	fmt.Printf("Requesting classification of image: %s\n", inputImg)
	resps, err := hfapigo.SendImageClassificationRequest(hfapigo.RecommendedImageClassificationModel, inputImg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for _, r := range resps {
		fmt.Printf("%s: %f\n", r.Label, r.Score)
	}
}
