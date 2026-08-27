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
		hfgo.WithModel("google/tapas-base-finetuned-wtq"),
	)

	table := map[string][]string{
		"Name": {"Alice", "Bob", "Carol"},
		"Age":  {"25", "30", "35"},
		"City": {"New York", "London", "Paris"},
	}
	question := "How old is Bob?"

	fmt.Printf("Question: %s\n", question)
	fmt.Printf("Table:\n")
	PrintJSON(table)
	fmt.Println("...")

	answers, err := client.AnswerTableQuestion(
		hfgo.TableQuestionAnsweringRequest{
			Input: hfgo.TableQuestionAnsweringInput{
				Question: question,
				Table:    table,
			},
		},
	)
	if err != nil {
		log.Fatalf("error running table question answering: %v\n", err)
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
