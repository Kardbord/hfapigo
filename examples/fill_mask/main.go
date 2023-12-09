package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kardbord/hfapigo"
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
		"The answer to life, the universe, and everything is [MASK].",
		"So long, and thanks for all the [MASK].",
	}

	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.FillMaskResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		fmresps, err := hfapigo.SendFillMaskRequest(hfapigo.RecommendedFillMaskModel, &hfapigo.FillMaskRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})
		ch <- ChanRv{fmresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return
			}

			for i := range inputs {
				fmt.Printf("\nInput %d results:\n", i)
				jsonBuf, err := json.MarshalIndent(chrv.resps[i].Masks, "", "  ")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(string(jsonBuf))
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
