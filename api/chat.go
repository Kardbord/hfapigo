package api

import "encoding/json"

// ChatRequest represents a completion request for the chat API.
// Output type depends on the Stream parameter.
type ChatRequest struct {
	// Number between -2.0 and 2.0. Positive values penalize new tokens based
	// on their existing frequency in the text so far, decreasing the model's
	// likelihood to repeat the same line verbatim.
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`

	// Whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned
	// in the content of message.
	LogProbs *bool `json:"logprobs,omitempty"`

	// The maximum number of tokens that can be generated in the chat completion.
	// Default: 1024. Minimum: 0.
	MaxTokens *int `json:"max_tokens,omitempty"`

	// A list of messages comprising the conversation so far.
	// Required.
	Messages []ChatMessage `json:"messages"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based
	// on whether they appear in the text so far, increasing the model's
	// likelihood to talk about new topics.
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`

	// Response format configuration. Possible types: text, json_schema, json_object.
	ResponseFormat *ChatResponseFormat `json:"response_format,omitempty"`

	// Seed for deterministic sampling.
	// Minimum: 0.
	Seed *int64 `json:"seed,omitempty"`

	// Up to 4 sequences where the API will stop generating further tokens.
	Stop []string `json:"stop,omitempty"`

	// If true, generated tokens are returned as a stream using SSE.
	// For more information about streaming, see:
	// https://huggingface.co/docs/text-generation-inference/conceptual/streaming
	Stream *bool `json:"stream,omitempty"`

	// Stream options for SSE responses.
	StreamOptions *ChatStreamOptions `json:"stream_options,omitempty"`

	// Sampling temperature to use, between 0 and 2.
	// We generally recommend altering this or TopP but not both.
	Temperature *float64 `json:"temperature,omitempty"`

	// Tool choice behavior. Possible values: auto, none, required, or
	// {"function":{"name":"..."}}. Meanings:
	// - auto: model can pick between generating a message or calling tools.
	// - none: model will not call any tool and instead generates a message.
	// - required: model must call one or more tools.
	ToolChoice json.RawMessage `json:"tool_choice,omitempty"`

	// A prompt to be appended before the tools.
	ToolPrompt *string `json:"tool_prompt,omitempty"`

	// A list of tools the model may call.
	// Currently, only functions are supported as a tool. Use this to provide
	// a list of functions the model may generate JSON inputs for.
	Tools []ChatTool `json:"tools,omitempty"`

	// Number of most likely tokens to return at each token position.
	// LogProbs must be true if this is set.
	// An integer between 0 and 5.
	TopLogProbs *int `json:"top_logprobs,omitempty"`

	// Nucleus sampling probability mass.
	// For example, 0.1 means only the tokens comprising the top 10% probability
	// mass are considered.
	TopP *float64 `json:"top_p,omitempty"`
}

// ChatMessage represents a single chat message.
type ChatMessage struct {
	// Role of the message author (for example: system, user, assistant, tool).
	// Required.
	Role string `json:"role"`
	// Optional name for the participant.
	Name *string `json:"name,omitempty"`

	// Content may be a string or a list of content chunks.
	// Either Content or ToolCalls should be supplied.
	Content ChatMessageContent `json:"content,omitempty"`

	// ToolCalls is used instead of Content when providing tool call messages.
	ToolCalls []ChatToolCall `json:"tool_calls,omitempty"`
}

// ChatMessageContent can be a string or []ChatMessageChunk.
// Use a string for pure text content or a chunk list for multimodal content.
type ChatMessageContent any

// ChatMessageChunk represents a content chunk (text or image URL).
// Type must be "text" with Text set, or "image_url" with ImageURL set.
type ChatMessageChunk struct {
	// Required when Type is "text".
	Text *string `json:"text,omitempty"`
	// Required when Type is "image_url".
	ImageURL *ChatImageURL `json:"image_url,omitempty"`
	// Required. Possible values: text, image_url.
	Type string `json:"type"`
}

// ChatImageURL contains the URL for an image chunk.
type ChatImageURL struct {
	// Required.
	URL string `json:"url"`
}

// ChatToolCall represents a tool call in a message.
type ChatToolCall struct {
	// Tool call ID.
	// Required.
	ID string `json:"id"`
	// Tool call type (for example: function).
	// Required.
	Type string `json:"type"`
	// Required.
	Function ChatFunctionDefinition `json:"function"`
}

// ChatFunctionDefinition describes a callable function.
type ChatFunctionDefinition struct {
	// Required.
	Name string `json:"name"`
	// A description of what the function does.
	Description *string `json:"description,omitempty"`
	// JSON schema describing function parameters.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// ChatResponseFormat configures the response format.
type ChatResponseFormat struct {
	// type values: text, json_schema, json_object.
	// For json_schema, JSONSchema is required.
	Type       string                `json:"type"`
	JSONSchema *ChatJSONSchemaConfig `json:"json_schema,omitempty"`
}

// ChatJSONSchemaConfig defines JSON schema response formatting.
type ChatJSONSchemaConfig struct {
	// The name of the response format.
	// Required.
	Name string `json:"name"`
	// A description of what the response format is for.
	Description *string `json:"description,omitempty"`
	// The schema for the response format as a JSON Schema object.
	// Learn how to build JSON schemas at https://json-schema.org/.
	Schema json.RawMessage `json:"schema,omitempty"`
	// Whether to enable strict schema adherence.
	Strict *bool `json:"strict,omitempty"`
}

// ChatStreamOptions configures streaming behavior.
type ChatStreamOptions struct {
	// If set, an additional chunk is streamed before the [DONE] message
	// showing overall usage. The usage field on this chunk shows the token
	// usage statistics for the entire request, and the choices field will
	// always be an empty array. All other chunks include a usage field with
	// a null value.
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// ChatTool represents a tool definition provided to the model.
type ChatTool struct {
	// Tool type (currently only function is supported).
	Type     string                 `json:"type"`
	Function ChatFunctionDefinition `json:"function"`
}

// ChatResponse represents a non-streaming completion response from the chat API.
// This is returned when Stream is false (the default).
type ChatResponse struct {
	// Required.
	ID string `json:"id"`
	// Unix timestamp in seconds.
	// Required.
	Created int64 `json:"created"`
	// Required.
	Model string `json:"model"`
	// Required.
	SystemFingerprint string `json:"system_fingerprint"`
	// Required.
	Choices []ChatChoice `json:"choices"`
	// Required.
	Usage ChatUsage `json:"usage"`
}

// ChatChoice is a single non-streaming completion choice.
type ChatChoice struct {
	// Required.
	FinishReason string `json:"finish_reason"`
	// Required.
	Index    int           `json:"index"`
	LogProbs *ChatLogProbs `json:"logprobs,omitempty"`
	// Required.
	Message ChatCompletionMessage `json:"message"`
}

// ChatLogProbs contains per-token log probabilities.
type ChatLogProbs struct {
	// Required.
	Content []ChatLogProb `json:"content"`
}

// ChatLogProb contains logprob information for a token.
type ChatLogProb struct {
	// Required.
	Token string `json:"token"`
	// Required.
	LogProb float64 `json:"logprob"`
	// Required.
	TopLogProbs []ChatTopLogProb `json:"top_logprobs"`
}

// ChatTopLogProb is a top log probability entry for a token position.
type ChatTopLogProb struct {
	// Required.
	Token string `json:"token"`
	// Required.
	LogProb float64 `json:"logprob"`
}

// ChatCompletionMessage is a message returned by the model.
// It is either a text message (Content) or a tool call message (ToolCalls).
type ChatCompletionMessage struct {
	Role string `json:"role"`
	// Content is present for text responses.
	Content *string `json:"content,omitempty"`
	// ToolCallID is set when returning tool-specific content.
	ToolCallID *string `json:"tool_call_id,omitempty"`
	// ToolCalls is present for tool call responses.
	ToolCalls []ChatToolCallOutput `json:"tool_calls,omitempty"`
}

// ChatToolCallOutput represents a tool call in a response message.
type ChatToolCallOutput struct {
	// Required.
	ID string `json:"id"`
	// Required.
	Type string `json:"type"`
	// Required.
	Function ChatFunctionCall `json:"function"`
}

// ChatFunctionCall represents a tool function call with arguments.
type ChatFunctionCall struct {
	// Required.
	Name string `json:"name"`
	// Required.
	Arguments   string  `json:"arguments"`
	Description *string `json:"description,omitempty"`
}

// ChatUsage contains token usage statistics.
type ChatUsage struct {
	// Required.
	CompletionTokens int `json:"completion_tokens"`
	// Required.
	PromptTokens int `json:"prompt_tokens"`
	// Required.
	TotalTokens int `json:"total_tokens"`
}

// ChatStreamResponse represents a streaming response chunk.
// This is returned when Stream is true.
type ChatStreamResponse struct {
	ID string `json:"id"`
	// Unix timestamp in seconds.
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	SystemFingerprint string             `json:"system_fingerprint"`
	Choices           []ChatStreamChoice `json:"choices"`
	Usage             *ChatUsage         `json:"usage,omitempty"`
}

// ChatStreamChoice is a single streaming completion choice.
type ChatStreamChoice struct {
	// Required.
	Delta        ChatStreamDelta `json:"delta"`
	FinishReason *string         `json:"finish_reason,omitempty"`
	// Required.
	Index    int           `json:"index"`
	LogProbs *ChatLogProbs `json:"logprobs,omitempty"`
}

// ChatStreamDelta holds incremental updates for a stream.
// Deltas may include content/role/tool_call_id or role/tool_calls.
type ChatStreamDelta struct {
	// Content is present for text deltas.
	Content *string `json:"content,omitempty"`
	// Role may be included with the first delta.
	Role *string `json:"role,omitempty"`
	// ToolCallID may be included for tool-specific content.
	ToolCallID *string `json:"tool_call_id,omitempty"`
	// ToolCalls is present for tool call deltas.
	ToolCalls []ChatStreamToolCall `json:"tool_calls,omitempty"`
}

// ChatStreamToolCall represents a tool call within a streaming delta.
type ChatStreamToolCall struct {
	// Required.
	ID string `json:"id"`
	// Required.
	Type string `json:"type"`
	// Required.
	Index    int                `json:"index"`
	Function ChatStreamFunction `json:"function"`
}

// ChatStreamFunction represents a streamed function call.
type ChatStreamFunction struct {
	// Required.
	Name string `json:"name"`
	// Required.
	Arguments string `json:"arguments"`
}
