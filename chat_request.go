package hfapigo

import (
	"encoding/json"
	"fmt"
)

// ChatRequest represents a completion request for the chat API.
// Output type depends on the Stream parameter.
type ChatRequest struct {
	// Model to use for the chat completion.
	// Required.
	Model *string `json:"model,omitempty"`

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

	// Response format configuration. Known types: text, json_schema, json_object.
	// Non-empty provider-specific values are also accepted.
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

	// Tool choice behavior. Known values: auto, none, required, or
	// {"function":{"name":"..."}}. Meanings:
	// - auto: model can pick between generating a message or calling tools.
	// - none: model will not call any tool and instead generates a message.
	// - required: model must call one or more tools.
	// Non-empty provider-specific values are also accepted.
	ToolChoice *ChatToolChoice `json:"tool_choice,omitempty"`

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

// MarshalJSON enforces required fields for ChatRequest.
//
//nolint:gocritic // hugeParam: using a value receiver guarantees no modifications to the caller
func (r ChatRequest) MarshalJSON() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	type alias ChatRequest

	return json.Marshal(alias(r))
}

//nolint:gocritic // hugeParam: using a value receiver guarantees no modifications to the caller
func (r ChatRequest) validate() error {
	if r.Model == nil || *r.Model == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat request model is required",
			Err:     nil,
		}
	}

	if len(r.Messages) == 0 {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat request must include at least one message",
			Err:     nil,
		}
	}

	return nil
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
	Content ChatMessageContent `json:"content"`

	// ToolCalls is used instead of Content when providing tool call messages.
	ToolCalls []ChatToolCall `json:"tool_calls,omitempty"`
}

// MarshalJSON enforces the union shape for ChatMessage.
func (m ChatMessage) MarshalJSON() ([]byte, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	type alias ChatMessage

	return json.Marshal(alias(m))
}

// validate enforces the union shape for ChatMessage.
func (m ChatMessage) validate() error {
	if m.Role == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat message: role must be set",
			Err:     nil,
		}
	}
	contentSet := m.Content.Text != nil || m.Content.Chunks != nil
	if contentSet && len(m.ToolCalls) > 0 {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat message: content and tool_calls are mutually exclusive",
			Err:     nil,
		}
	}
	if !contentSet && len(m.ToolCalls) == 0 {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat message: either content or tool_calls must be set",
			Err:     nil,
		}
	}

	return nil
}

// ChatMessageContent can be a string or []ChatMessageChunk.
// Use a string for pure text content or a chunk list for multimodal content.
type ChatMessageContent struct {
	// NOTE: ChatMessageContent uses json:"-" on its fields because the JSON payload is
	// a union (string or []ChatMessageChunk) handled by custom marshal/unmarshal.

	// Text holds plain string content.
	Text *string `json:"-"`
	// Chunks holds structured content chunks.
	Chunks []ChatMessageChunk `json:"-"`
}

// MarshalJSON enforces the union shape for ChatMessageContent.
func (c ChatMessageContent) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	if c.Text != nil {
		return json.Marshal(*c.Text)
	}
	if c.Chunks != nil {
		return json.Marshal(c.Chunks)
	}

	return []byte("null"), nil
}

// validate enforces the union shape for ChatMessageContent.
func (c ChatMessageContent) validate() error {
	if c.Text != nil && len(c.Chunks) > 0 {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat message content: both text and chunks set",
			Err:     nil,
		}
	}

	return nil
}

// ChatMessageChunk represents a content chunk (text or image URL).
// Type must be "text" with Text set, or "image_url" with ImageURL set.
type ChatMessageChunk struct {
	// Required when Type is "text".
	Text *string `json:"text,omitempty"`
	// Required when Type is "image_url".
	ImageURL *ChatImageURL `json:"image_url,omitempty"`
	// Required. Possible values: text, image_url.
	Type MessageChunkType `json:"type"`
}

// MarshalJSON enforces the union shape for ChatMessageChunk.
func (c ChatMessageChunk) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	type alias ChatMessageChunk

	return json.Marshal(alias(c))
}

// validate enforces the union shape for ChatMessageChunk.
func (c ChatMessageChunk) validate() error {
	switch c.Type {
	case MessageChunkTypeText:
		if c.Text == nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "chat message chunk: text requires text field",
				Err:     nil,
			}
		}
		if c.ImageURL != nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "chat message chunk: text cannot include image_url",
				Err:     nil,
			}
		}
	case MessageChunkTypeImageURL:
		if c.ImageURL == nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "chat message chunk: image_url requires image_url field",
				Err:     nil,
			}
		}
		if c.Text != nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "chat message chunk: image_url cannot include text",
				Err:     nil,
			}
		}
	default:
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat message chunk: type must be text or image_url",
			Err:     nil,
		}
	}

	return nil
}

// ChatImageURL contains the URL for an image chunk.
type ChatImageURL struct {
	// Required.
	URL string `json:"url"`
}

// MarshalJSON enforces the required URL on ChatImageURL.
func (u ChatImageURL) MarshalJSON() ([]byte, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}
	type alias ChatImageURL

	return json.Marshal(alias(u))
}

func (u ChatImageURL) validate() error {
	if u.URL == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat image url: url must be set",
			Err:     nil,
		}
	}

	return nil
}

// MessageChunkType enumerates supported chat message chunk types.
type MessageChunkType string

const (
	// MessageChunkTypeText represents a text chunk.
	MessageChunkTypeText MessageChunkType = "text"
	// MessageChunkTypeImageURL represents an image_url chunk.
	MessageChunkTypeImageURL MessageChunkType = "image_url"
)

// ChatToolCall represents a tool call in a message.
type ChatToolCall struct {
	// Tool call ID.
	// Required.
	ID string `json:"id"`
	// Tool call type (for example: function).
	// Required.
	Type string `json:"type"`
	// Required.
	Function ChatFunctionCall `json:"function"`
}

// MarshalJSON enforces the tool call shape for ChatToolCall.
func (c ChatToolCall) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	type alias ChatToolCall

	return json.Marshal(alias(c))
}

// validate enforces the tool call shape for ChatToolCall.
func (c ChatToolCall) validate() error {
	if c.ID == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat tool call: id must be set",
			Err:     nil,
		}
	}
	if c.Type == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat tool call: type must be set",
			Err:     nil,
		}
	}

	return nil
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

// MarshalJSON enforces required fields on ChatFunctionDefinition.
func (f ChatFunctionDefinition) MarshalJSON() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	type alias ChatFunctionDefinition

	return json.Marshal(alias(f))
}

func (f ChatFunctionDefinition) validate() error {
	if f.Name == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat function definition: name must be set",
			Err:     nil,
		}
	}

	return nil
}

// ChatResponseFormat configures the response format.
type ChatResponseFormat struct {
	// Known type values: text, json_schema, json_object.
	// Non-empty provider-specific values are also accepted.
	// For json_schema, JSONSchema is required.
	Type       ResponseFormatType    `json:"type"`
	JSONSchema *ChatJSONSchemaConfig `json:"json_schema,omitempty"`
}

// ResponseFormatType enumerates known response formats.
type ResponseFormatType string

const (
	// ResponseFormatTypeText requests a text response.
	ResponseFormatTypeText ResponseFormatType = "text"
	// ResponseFormatTypeJSONSchema requests a JSON schema response.
	ResponseFormatTypeJSONSchema ResponseFormatType = "json_schema"
	// ResponseFormatTypeJSONObject requests a JSON object response.
	ResponseFormatTypeJSONObject ResponseFormatType = "json_object"
)

// MarshalJSON enforces the union shape for ChatResponseFormat.
func (r ChatResponseFormat) MarshalJSON() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	type alias ChatResponseFormat

	return json.Marshal(alias(r))
}

// validate enforces the union shape for ChatResponseFormat.
func (r ChatResponseFormat) validate() error {
	switch r.Type {
	case ResponseFormatTypeJSONSchema:
		if r.JSONSchema == nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "chat response format: json_schema requires json_schema field",
				Err:     nil,
			}
		}
	case "":
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat response format: type must be set",
			Err:     nil,
		}
	case ResponseFormatTypeText, ResponseFormatTypeJSONObject:
		if r.JSONSchema != nil {
			return &SDKError{
				Kind: SDKErrorKindConfiguration,
				Message: fmt.Sprintf(
					"chat response format: %s cannot include json_schema field",
					r.Type,
				),
				Err: nil,
			}
		}
	default:
		// default is reserved for provider-specific response format types.
		if r.JSONSchema != nil {
			return &SDKError{
				Kind: SDKErrorKindConfiguration,
				Message: fmt.Sprintf(
					"chat response format: %s cannot include json_schema field",
					r.Type,
				),
				Err: nil,
			}
		}
	}

	return nil
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

// MarshalJSON enforces the tool shape for ChatTool.
func (t ChatTool) MarshalJSON() ([]byte, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}
	type alias ChatTool

	return json.Marshal(alias(t))
}

// validate enforces the tool shape for ChatTool.
func (t ChatTool) validate() error {
	if t.Type == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat tool: type must be set",
			Err:     nil,
		}
	}

	return nil
}

// ChatToolChoice represents the tool choice union type.
// Either Mode is set to auto/none/required, or Function is set.
type ChatToolChoice struct {
	// NOTE: ChatToolChoice uses json:"-" on its fields because the JSON payload is
	// a union (string or {"function":{...}}) handled by custom marshal/unmarshal.

	// Mode is a tool choice mode. Known values: auto, none, required.
	// Non-empty provider-specific values are also accepted.
	Mode *ToolChoiceMode `json:"-"`
	// Function selects a specific tool function by name.
	Function *ChatFunctionName `json:"-"`
}

// ToolChoiceMode enumerates known tool choice modes.
type ToolChoiceMode string

const (
	// ToolChoiceModeAuto lets the provider decide the tool choice.
	ToolChoiceModeAuto ToolChoiceMode = "auto"
	// ToolChoiceModeNone disables tool usage.
	ToolChoiceModeNone ToolChoiceMode = "none"
	// ToolChoiceModeRequired requires tool usage.
	ToolChoiceModeRequired ToolChoiceMode = "required"
)

// toolChoiceFunctionPayload is the JSON object shape for function tool choices.
type toolChoiceFunctionPayload struct {
	Function *ChatFunctionName `json:"function"`
}

// MarshalJSON enforces the union shape for ChatToolChoice.
func (t ChatToolChoice) MarshalJSON() ([]byte, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}
	if t.Mode != nil {
		return json.Marshal(*t.Mode)
	}
	if t.Function != nil {
		return json.Marshal(toolChoiceFunctionPayload{
			Function: t.Function,
		})
	}

	return []byte("null"), nil
}

// validate enforces the union shape for ChatToolChoice.
func (t ChatToolChoice) validate() error {
	if t.Mode != nil && t.Function != nil {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "tool choice: mode and function are mutually exclusive",
			Err:     nil,
		}
	}
	if t.Mode != nil && *t.Mode == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "tool choice: mode must be set",
			Err:     nil,
		}
	}

	return nil
}

// ChatFunctionName identifies a tool function by name.
type ChatFunctionName struct {
	// Required.
	Name string `json:"name"`
}

// MarshalJSON enforces required fields on ChatFunctionName.
func (f ChatFunctionName) MarshalJSON() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	type alias ChatFunctionName

	return json.Marshal(alias(f))
}

func (f ChatFunctionName) validate() error {
	if f.Name == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat function name: name must be set",
			Err:     nil,
		}
	}

	return nil
}
