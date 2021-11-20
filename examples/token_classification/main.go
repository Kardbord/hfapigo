package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TannerKvarfordt/hfapigo"
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
		"My name is Sarah Jessica Parker but you can call me Jessica",
		"My name is Clara and I live in Berkeley, California.",
	}

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.TokenClassificationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
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
				for _, group := range chrv.resps[i].EntityGroups {
					fmt.Printf("Word: \"%s\" Label: %s Score: %f\n", group.Word, group.Name, group.Score)
				}
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
