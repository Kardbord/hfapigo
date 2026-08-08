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

	// Create a new client with your API token and desired model
	client := hfgo.NewClient(
		hfgo.WithToken(token),
		hfgo.WithModel("google-bert/bert-base-uncased"),
	)

	inputs := []string{
		"I [MASK] my dog everyday.",
		"She is a [MASK] programmer.",
		"The meeting was [MASK] due to the storm.",
	}

	fmt.Println("Filling the mask in inputs:")
	PrintJSON(inputs)
	fmt.Println("...")

	// Make the batched fill mask request
	predictions, err := client.FillMaskBatch(
		hfgo.FillMaskBatchRequest{
			Inputs: inputs,
		},
	)
	if err != nil {
		log.Fatalf("error running batched fill mask: %v\n", err)
	}

	fmt.Println("Results:")
	PrintJSON(predictions)
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
