package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Kardbord/hfapigo/v2"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	rand.Seed(time.Now().UnixNano())
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

const inputImg = "./test-image.png"

func main() {
	fmt.Printf("Sending image to text request for image %s", inputImg)

	type ChanRv struct {
		resps []*hfapigo.ImageToTextResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		resps, err := hfapigo.SendImageToTextRequest(hfapigo.RecommendedImageToTextModel, inputImg)
		ch <- ChanRv{resps: resps, err: err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			for _, r := range chrv.resps {
				fmt.Printf("Caption: %s\n", r.GeneratedText)
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 400)
		}
	}
}
