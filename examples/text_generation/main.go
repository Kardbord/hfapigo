package main

import (
	"fmt"
	"os"
	"strings"
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
	inputs := []string{
		"The answer to life, the universe, and everything is",
		"Somebody once told me that the world is gonna roll me",
	}
	const numReturnSequences = 3

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))

	type ChanRv struct {
		resps []*hfapigo.TextGenerationResponse
		err   error
	}
	ch := make(chan ChanRv)

	fmt.Print("Sending request")
	go func() {
		resps, err := hfapigo.SendTextGenerationRequest(hfapigo.RecommendedTextGenerationModel, &hfapigo.TextGenerationRequest{
			Inputs:     inputs,
			Parameters: *hfapigo.NewTextGenerationParameters().SetNumReturnSequences(numReturnSequences),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
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
			for i := range inputs {
				fmt.Printf("\nInput %d results:\n", i)
				for _, gt := range chrv.resps[i].GeneratedTexts {
					gt = strings.Replace(gt, "\n", " ", -1)
					gt = strings.Replace(gt, "\r", " ", -1)
					fmt.Println(gt)
				}
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
