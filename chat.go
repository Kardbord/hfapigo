package hfapigo

import (
	"encoding/json"
	"errors"
	"fmt"
)

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

// UnmarshalJSON enforces the union shape for ChatMessage.
func (m *ChatMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("chat message: nil receiver")
	}
	type alias ChatMessage
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatMessage(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*m = out

	return nil
}

// validate enforces the union shape for ChatMessage.
func (m ChatMessage) validate() error {
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

// UnmarshalJSON enforces the union shape for ChatMessageContent.
func (c *ChatMessageContent) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("chat message content: nil receiver")
	}
	if string(data) == "null" {
		return nil
	}
	if len(data) == 0 {
		return errors.New("chat message content: empty payload")
	}
	switch data[0] {
	case '"':
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		c.Text = &text
		c.Chunks = nil

		return c.validate()
	case '[':
		var chunks []ChatMessageChunk
		if err := json.Unmarshal(data, &chunks); err != nil {
			return err
		}
		c.Text = nil
		c.Chunks = chunks

		return c.validate()
	default:
		return errors.New("chat message content: expected string or array")
	}
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

// UnmarshalJSON enforces the union shape for ChatMessageChunk.
func (c *ChatMessageChunk) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("chat message chunk: nil receiver")
	}
	type alias ChatMessageChunk
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatMessageChunk(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*c = out

	return nil
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
	Function ChatFunctionDefinition `json:"function"`
}

// MarshalJSON enforces the tool call shape for ChatToolCall.
func (c ChatToolCall) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	type alias ChatToolCall

	return json.Marshal(alias(c))
}

// UnmarshalJSON enforces the tool call shape for ChatToolCall.
func (c *ChatToolCall) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("chat tool call: nil receiver")
	}
	type alias ChatToolCall
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatToolCall(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*c = out

	return nil
}

// validate enforces the tool call shape for ChatToolCall.
func (c ChatToolCall) validate() error {
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

// UnmarshalJSON enforces the union shape for ChatResponseFormat.
func (r *ChatResponseFormat) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("chat response format: nil receiver")
	}
	type alias ChatResponseFormat
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatResponseFormat(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*r = out

	return nil
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

// UnmarshalJSON enforces the tool shape for ChatTool.
func (t *ChatTool) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("chat tool: nil receiver")
	}
	type alias ChatTool
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatTool(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*t = out

	return nil
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

// UnmarshalJSON enforces the union shape for ChatToolChoice.
func (t *ChatToolChoice) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("tool choice: nil receiver")
	}
	if string(data) == "null" {
		return nil
	}
	if len(data) == 0 {
		return errors.New("tool choice: empty payload")
	}
	mode, function, err := parseToolChoice(data)
	if err != nil {
		return err
	}
	t.Mode = mode
	t.Function = function

	return t.validate()
}

func parseToolChoice(data []byte) (*ToolChoiceMode, *ChatFunctionName, error) {
	// TODO: Is this the best way to distinguish between string and struct return types?
	switch data[0] {
	case '"':
		var mode ToolChoiceMode
		if err := json.Unmarshal(data, &mode); err != nil {
			return nil, nil, err
		}
		if mode == "" {
			return nil, nil, errors.New("tool choice: mode must be set")
		}

		return &mode, nil, nil
	case '{':
		var payload toolChoiceFunctionPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, nil, err
		}
		if payload.Function == nil {
			return nil, nil, errors.New("tool choice: function object required")
		}

		return nil, payload.Function, nil
	default:
		return nil, nil, errors.New("tool choice: expected string or object")
	}
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

// MarshalJSON enforces the union shape for ChatCompletionMessage.
func (m ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	type alias ChatCompletionMessage

	return json.Marshal(alias(m))
}

// UnmarshalJSON enforces the union shape for ChatCompletionMessage.
func (m *ChatCompletionMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("completion message: nil receiver")
	}
	type alias ChatCompletionMessage
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatCompletionMessage(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*m = out

	return nil
}

// validate enforces the union shape for ChatCompletionMessage.
func (m ChatCompletionMessage) validate() error {
	if m.Content != nil {
		if len(m.ToolCalls) > 0 {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "completion message: content and tool_calls are mutually exclusive",
				Err:     nil,
			}
		}
	} else if len(m.ToolCalls) == 0 {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "completion message: either content or tool_calls must be set",
			Err:     nil,
		}
	}
	if m.Content == nil && m.ToolCallID != nil {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "completion message: tool_call_id requires content",
			Err:     nil,
		}
	}

	return nil
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

// MarshalJSON enforces the tool call shape for ChatToolCallOutput.
func (c ChatToolCallOutput) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	type alias ChatToolCallOutput

	return json.Marshal(alias(c))
}

// UnmarshalJSON enforces the tool call shape for ChatToolCallOutput.
func (c *ChatToolCallOutput) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("chat tool call output: nil receiver")
	}
	type alias ChatToolCallOutput
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatToolCallOutput(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*c = out

	return nil
}

// validate enforces the tool call shape for ChatToolCallOutput.
func (c ChatToolCallOutput) validate() error {
	if c.Type == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat tool call output: type must be set",
			Err:     nil,
		}
	}

	return nil
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

// MarshalJSON enforces the union shape for ChatStreamDelta.
func (d ChatStreamDelta) MarshalJSON() ([]byte, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}
	type alias ChatStreamDelta

	return json.Marshal(alias(d))
}

// UnmarshalJSON enforces the union shape for ChatStreamDelta.
func (d *ChatStreamDelta) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("stream delta: nil receiver")
	}
	type alias ChatStreamDelta
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatStreamDelta(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*d = out

	return nil
}

// validate enforces the union shape for ChatStreamDelta.
func (d ChatStreamDelta) validate() error {
	if len(d.ToolCalls) > 0 {
		if d.Content != nil || d.ToolCallID != nil {
			return &SDKError{
				Kind:    SDKErrorKindConfiguration,
				Message: "stream delta: tool_calls cannot include content or tool_call_id",
				Err:     nil,
			}
		}
	}

	return nil
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

// MarshalJSON enforces the tool call shape for ChatStreamToolCall.
func (c ChatStreamToolCall) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	type alias ChatStreamToolCall

	return json.Marshal(alias(c))
}

// UnmarshalJSON enforces the tool call shape for ChatStreamToolCall.
func (c *ChatStreamToolCall) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("chat stream tool call: nil receiver")
	}
	type alias ChatStreamToolCall
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatStreamToolCall(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*c = out

	return nil
}

// validate enforces the tool call shape for ChatStreamToolCall.
func (c ChatStreamToolCall) validate() error {
	if c.Type == "" {
		return &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "chat stream tool call: type must be set",
			Err:     nil,
		}
	}

	return nil
}

// ChatStreamFunction represents a streamed function call.
type ChatStreamFunction struct {
	// Required.
	Name string `json:"name"`
	// Required.
	Arguments string `json:"arguments"`
}
