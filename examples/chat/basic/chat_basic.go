package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Kardbord/hfgo/v4"
)

// This example demonstrates how to use the ChatService for basic (non-streaming)
// chat completions. It sends a single message to the model and prints the response.
func main() {
	token := os.Getenv("HUGGING_FACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGING_FACE_TOKEN environment variable is not set")
	}

	// Create a new client with your API token and desired model
	client := hfgo.NewClient(
		hfgo.WithToken(token),
		hfgo.WithModel("deepseek-ai/DeepSeek-R1"),
	)

	// Create a chat request with a simple message
	request := &hfgo.ChatRequest{
		Messages: []hfgo.ChatMessage{
			{
				Role: "user",
				Content: hfgo.ChatMessageContent{
					Text: Ptr("Hello! What is the capital of France?"),
				},
			},
		},
		MaxTokens: Ptr(1024),
	}

	// Send the request and get the response
	response, err := client.Chat().Complete(request)
	if err != nil {
		log.Fatalf("Failed to complete chat request: %v", err)
	}

	// Print the response
	fmt.Println("Chat Completion Response:")
	fmt.Printf("  ID: %s\n", response.ID)
	fmt.Printf("  Model: %s\n", response.Model)
	fmt.Printf("  Choices: %d\n", len(response.Choices))

	for i, choice := range response.Choices {
		fmt.Printf("\n  Choice %d:\n", i)
		fmt.Printf("    Finish Reason: %s\n", choice.FinishReason)
		fmt.Printf("    Role: %s\n", choice.Message.Role)
		if choice.Message.Content != nil {
			fmt.Printf("    Content: %s\n", *choice.Message.Content)
		}
	}

	fmt.Printf("\nUsage:\n")
	fmt.Printf("  Prompt Tokens: %d\n", response.Usage.PromptTokens)
	fmt.Printf("  Completion Tokens: %d\n", response.Usage.CompletionTokens)
	fmt.Printf("  Total Tokens: %d\n", response.Usage.TotalTokens)
}

// Helper function to create pointers from values
func Ptr[T any](v T) *T {
	return &v
}
