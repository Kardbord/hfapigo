package api

import "github.com/Kardbord/hfapigo/v4/internal/request"

// RequestOption represents a functional option that configures client requests.
type RequestOption = request.RequestOption

// ChatRequest represents a completion request for the chat API.
type ChatRequest struct {
	Inputs string `json:"inputs"`
}

// ChatResponse represents a completion response from the chat API.
type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}
