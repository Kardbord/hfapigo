package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kardbord/hfapigo/v3"
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
		"You know, I find you quite fascinating.",
		"I don't really care for your disposition.",
	}

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.TextClassificationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		tcresps, err := hfapigo.SendTextClassificationRequest(hfapigo.RecommendedTextClassificationModel, &hfapigo.TextClassificationRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{resps: tcresps, err: err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}

			for i, input := range inputs {
				fmt.Printf("\nInput %d: %s\n", i, input)
				for _, label := range chrv.resps[i].Labels {
					fmt.Printf("Label: %s Score: %f\n", label.Name, label.Score)
				}
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
