package api

// ChatRequest represents a completion request for the chat API.
type ChatRequest struct {
	Inputs string `json:"inputs"`
}

// ChatResponse represents a completion response from the chat API.
type ChatResponse struct {
	GeneratedText string `json:"generated_text"`
}
