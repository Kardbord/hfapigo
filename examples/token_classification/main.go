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
		"My name is Sarah Jessica Parker but you can call me Jessica",
		"My name is Clara and I live in Berkeley, California.",
	}

	classifyNoAggregation(inputs)
	fmt.Println("===========================================================")
	classifyWithAggregation(inputs)
}

func classifyNoAggregation(inputs []string) {
	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request (AggregationStrategy=%s)", hfapigo.AggregationStrategyNone)

	type ChanRv struct {
		resps []*hfapigo.TokenClassificationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
			Inputs:     inputs,
			Parameters: *hfapigo.NewTokenClassificationParameters().SetAggregationStrategy(hfapigo.AggregationStrategyNone),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
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
				for _, entity := range chrv.resps[i].Entities {
					fmt.Printf("Entity: \"%s\" Label: \"%s\" Score: %f\n", entity.Entity, entity.Label, entity.Score)
				}
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func classifyWithAggregation(inputs []string) {
	fmt.Printf("Inputs: [\"%s\"]\n", strings.Join(inputs, `", "`))
	fmt.Printf("\nSending request (AggregationStrategy=%s)", hfapigo.AggregationStrategySimple)

	type ChanRv struct {
		resps []*hfapigo.TokenClassificationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		tcresps, err := hfapigo.SendTokenClassificationRequest(hfapigo.RecommendedTokenClassificationModel, &hfapigo.TokenClassificationRequest{
			Inputs:     inputs,
			Parameters: *hfapigo.NewTokenClassificationParameters().SetAggregationStrategy(hfapigo.AggregationStrategySimple),
			Options:    *hfapigo.NewOptions().SetWaitForModel(true),
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
				for _, entity := range chrv.resps[i].Entities {
					fmt.Printf("Entity: \"%s\" Label: \"%s\" Score: %f\n", entity.Entity, entity.Label, entity.Score)
				}
			}
			return
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}
