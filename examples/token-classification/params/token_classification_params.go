package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Kardbord/hfgo/v4"
)

func main() {
	token := os.Getenv("HF_TOKEN")
	if token == "" {
		log.Fatal("HF_TOKEN environment variable is not set")
	}

	client := hfgo.NewClient(
		hfgo.WithToken(token),
		hfgo.WithModel("dslim/bert-base-NER"),
	)

	input := "My name is Sarah and I live in London."

	fmt.Println("Classifying tokens with aggregation strategy:")
	PrintJSON(input)
	fmt.Println("...")

	entities, err := client.ClassifyTokens(
		hfgo.TokenClassificationRequest{
			Input: input,
			Parameters: &hfgo.TokenClassificationParameters{
				AggregationStrategy: Ptr(hfgo.TokenClassificationAggregationSimple),
				IgnoreLabels:        []string{"O"},
			},
		},
	)
	if err != nil {
		log.Fatalf("error running token classification: %v\n", err)
	}

	fmt.Println("Results:")
	PrintJSON(entities)
}

func Ptr[T any](v T) *T {
	return &v
}

func PrintJSON[T any](v T) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("error printing JSON: %v\n", err)
	}

	fmt.Println(string(b))
}
