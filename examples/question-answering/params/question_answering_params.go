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
		hfgo.WithModel("deepset/roberta-base-squad2"),
	)

	question := "What is the capital of France?"
	context := "France is a country in Europe. Its capital is Paris. It is known for the Eiffel Tower."

	fmt.Printf("Question: %s\n", question)
	fmt.Printf("Context:  %s\n", context)
	fmt.Println("...")

	topK := 3
	maxAnswerLen := 10
	answers, err := client.AnswerQuestion(
		hfgo.QuestionAnsweringRequest{
			Input: hfgo.QuestionAnsweringInput{
				Question: question,
				Context:  context,
			},
			Parameters: &hfgo.QuestionAnsweringParameters{
				TopK:         &topK,
				MaxAnswerLen: &maxAnswerLen,
			},
		},
	)
	if err != nil {
		log.Fatalf("error running question answering: %v\n", err)
	}

	fmt.Println("Results:")
	PrintJSON(answers)
}

func PrintJSON[T any](v T) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("error printing JSON: %v\n", err)
	}

	fmt.Println(string(b))
}
