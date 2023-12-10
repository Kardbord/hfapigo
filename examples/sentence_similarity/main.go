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
	sourceSentence := "That is a happy person"
	inputSentences := []string{"That is a happy dog", "That is a very happy person", "Today is a sunny day"}

	fmt.Printf("Source Sentence: %s\n", sourceSentence)
	fmt.Printf("Input Sentences: %s\n", strings.Join(inputSentences, `", "`))
	fmt.Print("\nSending request")

	type ChanRv struct {
		resp *hfapigo.SentenceSimilarityResponse
		err  error
	}
	ch := make(chan ChanRv)

	go func() {
		resp, err := hfapigo.SendSentenceSimilarityRequest(hfapigo.RecommendedSentenceSimilarityModel, &hfapigo.SentenceSimilarityRequest{
			Inputs: hfapigo.SentenceSimilarityInputs{
				SourceSentence: sourceSentence,
				Sentences:      inputSentences,
			},
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{resp, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}

			fmt.Println()
			for i, s := range inputSentences {
				fmt.Printf("%s -> %f similarity\n", s, (*chrv.resp)[i])
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
