package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Kardbord/hfapigo/v2"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func main() {
	input := "The answer to life, the universe, and everything is"

	fmt.Printf("Input: \"%s\"\n", input)

	type ChanRv struct {
		resps []*hfapigo.TextGenerationResponse
		err   error
	}
	ch := make(chan ChanRv)

	fmt.Print("Sending request")
	go func() {
		resps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Input:   input,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		ch <- ChanRv{resps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}
			fmt.Printf("Response: %s\n", chrv.resps[0].GeneratedText)
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
