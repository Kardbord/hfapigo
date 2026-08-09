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
		hfgo.WithModel("google-t5/t5-small"),
	)

	input := "The weather is lovely and sunny today."

	cleanUp := true
	truncation := hfgo.TranslationTruncationOnlyFirst

	fmt.Println("Translating input:")
	PrintJSON(input)
	fmt.Println("...")

	// Make the translation request with custom parameters
	translations, err := client.Translate(
		hfgo.TranslationRequest{
			Input: input,
			Parameters: &hfgo.TranslationParameters{
				CleanUpTokenizationSpaces: &cleanUp,
				Truncation:                &truncation,
			},
		},
	)
	if err != nil {
		log.Fatalf("error running translation: %v\n", err)
	}

	fmt.Println("Results:")
	PrintJSON(translations)
}

func PrintJSON[T any](v T) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("error printing JSON: %v\n", err)
	}

	fmt.Println(string(b))
}
