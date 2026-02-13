package hfapigo

// TODO: Review these AI-generated tests

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfapigo/v4/internal/testutils"
)

func TestChatMessageContent_Marshal(t *testing.T) {
	t.Parallel()

	text := "hello"
	imgURL := "https://example.com/image.png"

	cases := []struct {
		name     string
		value    ChatMessageContent
		wantJSON string
		wantErr  bool
	}{
		{
			name:     "text",
			value:    ChatMessageContent{Text: &text},
			wantJSON: `"hello"`,
		},
		{
			name:     "empty",
			value:    ChatMessageContent{},
			wantJSON: `null`,
		},
		{
			name: "chunks",
			value: ChatMessageContent{
				Chunks: []ChatMessageChunk{
					{Type: MessageChunkTypeImageURL, ImageURL: &ChatImageURL{URL: imgURL}},
				},
			},
			wantJSON: `[{"image_url":{"url":"https://example.com/image.png"},"type":"image_url"}]`,
		},
		{
			name: "both",
			value: ChatMessageContent{
				Text:   &text,
				Chunks: []ChatMessageChunk{{Type: MessageChunkTypeText, Text: &text}},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
			if tc.wantJSON != "" && string(data) != tc.wantJSON {
				t.Fatalf("unexpected json: %s", string(data))
			}
		})
	}
}

func TestChatMessageContent_Unmarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantText  *string
		wantChunk bool
		wantErr   bool
	}{
		{
			name:      "string",
			unmarshal: `"hi"`,
			wantText:  testutils.Ptr("hi"),
		},
		{
			name:      "null",
			unmarshal: `null`,
		},
		{
			name:      "array",
			unmarshal: `[{"text":"ok","type":"text"}]`,
			wantChunk: true,
		},
		{
			name:      "invalid",
			unmarshal: `123`,
			wantErr:   true,
		},
		{
			name:      "object",
			unmarshal: `{}`,
			wantErr:   true,
		},
		{
			name:      "empty payload",
			unmarshal: ``,
			wantErr:   true,
		},
		{
			name:      "invalid chunk type",
			unmarshal: `[{"text":"ok","type":"other"}]`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatMessageContent
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
			if tc.wantText != nil {
				if got.Text == nil || *got.Text != *tc.wantText {
					t.Fatalf("unexpected text: %+v", got.Text)
				}
				testutils.RequireNil(t, got.Chunks)
			}
			if tc.wantChunk && len(got.Chunks) == 0 {
				t.Fatalf("expected chunks")
			}
			if !tc.wantChunk && len(got.Chunks) > 0 {
				t.Fatalf("unexpected chunks")
			}
			if tc.wantChunk {
				testutils.RequireNil(t, got.Text)
			}
		})
	}
}

func TestChatMessageChunk_Validation(t *testing.T) {
	t.Parallel()

	text := "hello"

	cases := []struct {
		name    string
		value   ChatMessageChunk
		wantErr bool
	}{
		{
			name:  "text",
			value: ChatMessageChunk{Type: MessageChunkTypeText, Text: &text},
		},
		{
			name: "image_url",
			value: ChatMessageChunk{
				Type:     MessageChunkTypeImageURL,
				ImageURL: &ChatImageURL{URL: "x"},
			},
		},
		{
			name:    "missing text",
			value:   ChatMessageChunk{Type: MessageChunkTypeText},
			wantErr: true,
		},
		{
			name:    "missing image",
			value:   ChatMessageChunk{Type: MessageChunkTypeImageURL},
			wantErr: true,
		},
		{
			name:    "invalid type",
			value:   ChatMessageChunk{Type: MessageChunkType("other")},
			wantErr: true,
		},
		{
			name: "text with image_url",
			value: ChatMessageChunk{
				Type:     MessageChunkTypeText,
				Text:     &text,
				ImageURL: &ChatImageURL{URL: "x"},
			},
			wantErr: true,
		},
		{
			name: "image_url with text",
			value: ChatMessageChunk{
				Type:     MessageChunkTypeImageURL,
				Text:     &text,
				ImageURL: &ChatImageURL{URL: "x"},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatMessageChunk_UnmarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "text with image_url",
			unmarshal: `{"text":"hi","image_url":{"url":"x"},"type":"text"}`,
			wantErr:   true,
		},
		{
			name:      "image_url with text",
			unmarshal: `{"text":"hi","image_url":{"url":"x"},"type":"image_url"}`,
			wantErr:   true,
		},
		{
			name:      "text with extra image_url",
			unmarshal: `{"text":"hi","type":"text","image_url":{"url":"x"}}`,
			wantErr:   true,
		},
		{
			name:      "image_url with extra text",
			unmarshal: `{"image_url":{"url":"x"},"type":"image_url","text":"hi"}`,
			wantErr:   true,
		},
		{
			name:      "missing type",
			unmarshal: `{"text":"hi"}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatMessageChunk
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatMessage_Validation(t *testing.T) {
	t.Parallel()

	text := "hi"

	cases := []struct {
		name    string
		value   ChatMessage
		wantErr bool
	}{
		{
			name:  "content",
			value: ChatMessage{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
		{
			name: "tool_calls",
			value: ChatMessage{Role: "assistant", ToolCalls: []ChatToolCall{{
				ID:       "id",
				Type:     "function",
				Function: ChatFunctionDefinition{Name: "do"},
			}}},
		},
		{
			name: "both",
			value: ChatMessage{
				Role:    "assistant",
				Content: ChatMessageContent{Text: &text},
				ToolCalls: []ChatToolCall{
					{ID: "id", Type: "function", Function: ChatFunctionDefinition{Name: "do"}},
				},
			},
			wantErr: true,
		},
		{
			name:    "neither",
			value:   ChatMessage{Role: "assistant"},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatMessage_UnmarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "both content and tool_calls",
			unmarshal: `{"role":"assistant","content":"hi","tool_calls":[{"id":"id","type":"function","function":{"name":"do"}}]}`,
			wantErr:   true,
		},
		{
			name:      "content null with tool_calls",
			unmarshal: `{"role":"assistant","content":null,"tool_calls":[{"id":"id","type":"function","function":{"name":"do"}}]}`,
		},
		{
			name:      "content null without tool_calls",
			unmarshal: `{"role":"assistant","content":null}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatMessage
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatMessage_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		unmarshal     string
		wantRole      string
		wantContent   *string
		wantChunkType MessageChunkType
		wantChunks    int
		wantToolCalls int
	}{
		{
			name:        "content string",
			unmarshal:   `{"role":"user","content":"hi"}`,
			wantRole:    "user",
			wantContent: testutils.Ptr("hi"),
		},
		{
			name:          "content chunks",
			unmarshal:     `{"role":"user","content":[{"text":"hi","type":"text"}]}`,
			wantRole:      "user",
			wantChunks:    1,
			wantChunkType: MessageChunkTypeText,
		},
		{
			name:          "tool_calls",
			unmarshal:     `{"role":"assistant","tool_calls":[{"id":"id","type":"function","function":{"name":"fn"}}]}`,
			wantRole:      "assistant",
			wantToolCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatMessage
			testutils.RequireNoError(t, json.Unmarshal([]byte(tc.unmarshal), &got))
			if got.Role != tc.wantRole {
				t.Fatalf("unexpected role: %v", got.Role)
			}
			if tc.wantContent != nil {
				if got.Content.Text == nil || *got.Content.Text != *tc.wantContent {
					t.Fatalf("unexpected content: %+v", got.Content.Text)
				}
				testutils.RequireNil(t, got.Content.Chunks)
			}
			if tc.wantChunks > 0 {
				if len(got.Content.Chunks) != tc.wantChunks {
					t.Fatalf("unexpected chunks: %+v", got.Content.Chunks)
				}
				if got.Content.Chunks[0].Type != tc.wantChunkType {
					t.Fatalf("unexpected chunk type: %+v", got.Content.Chunks[0].Type)
				}
				testutils.RequireNil(t, got.Content.Text)
			}
			if tc.wantToolCalls > 0 {
				if len(got.ToolCalls) != tc.wantToolCalls {
					t.Fatalf("unexpected tool calls: %+v", got.ToolCalls)
				}
				if got.ToolCalls[0].ID != "id" || got.ToolCalls[0].Type != "function" ||
					got.ToolCalls[0].Function.Name != "fn" {
					t.Fatalf("unexpected tool call: %+v", got.ToolCalls[0])
				}
				testutils.RequireNil(t, got.Content.Text)
				testutils.RequireNil(t, got.Content.Chunks)
			}
			if tc.wantToolCalls == 0 && len(got.ToolCalls) > 0 {
				t.Fatalf("unexpected tool calls: %+v", got.ToolCalls)
			}
		})
	}
}

func TestChatRequest_MarshalSuccess(t *testing.T) {
	t.Parallel()

	text := "hi"
	imgURL := "https://example.com/image.png"
	req := ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
			{
				Role: "user",
				Content: ChatMessageContent{
					Chunks: []ChatMessageChunk{
						{Type: MessageChunkTypeImageURL, ImageURL: &ChatImageURL{URL: imgURL}},
					},
				},
			},
		},
	}

	data, err := json.Marshal(req)
	testutils.RequireNoError(t, err)

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unexpected json: %v", err)
	}
	messages, ok := got["messages"].([]any)
	if !ok || len(messages) != 2 {
		t.Fatalf("unexpected messages: %#v", got["messages"])
	}
	first, ok := messages[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected first message: %#v", messages[0])
	}
	if first["role"] != "user" || first["content"] != "hi" {
		t.Fatalf("unexpected first message fields: %#v", first)
	}
	second, ok := messages[1].(map[string]any)
	if !ok {
		t.Fatalf("unexpected second message: %#v", messages[1])
	}
	if second["role"] != "user" {
		t.Fatalf("unexpected second message role: %#v", second["role"])
	}
	chunks, ok := second["content"].([]any)
	if !ok || len(chunks) != 1 {
		t.Fatalf("unexpected content chunks: %#v", second["content"])
	}
	chunk, ok := chunks[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected chunk: %#v", chunks[0])
	}
	if chunk["type"] != string(MessageChunkTypeImageURL) {
		t.Fatalf("unexpected chunk type: %#v", chunk["type"])
	}
	imageURL, ok := chunk["image_url"].(map[string]any)
	if !ok || imageURL["url"] != imgURL {
		t.Fatalf("unexpected image_url: %#v", chunk["image_url"])
	}
}

func TestChatToolChoice_Marshal(t *testing.T) {
	t.Parallel()

	mode := ToolChoiceMode("provider-mode")

	cases := []struct {
		name     string
		value    ChatToolChoice
		wantJSON string
		wantErr  bool
	}{
		{
			name:     "mode",
			value:    ChatToolChoice{Mode: &mode},
			wantJSON: `"provider-mode"`,
		},
		{
			name:     "null",
			value:    ChatToolChoice{},
			wantJSON: `null`,
		},
		{
			name:     "function",
			value:    ChatToolChoice{Function: &ChatFunctionName{Name: "do"}},
			wantJSON: `{"function":{"name":"do"}}`,
		},
		{
			name:    "both",
			value:   ChatToolChoice{Mode: &mode, Function: &ChatFunctionName{Name: "do"}},
			wantErr: true,
		},
		{
			name:    "empty mode",
			value:   ChatToolChoice{Mode: testutils.Ptr(ToolChoiceMode(""))},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
			if tc.wantJSON != "" && string(data) != tc.wantJSON {
				t.Fatalf("unexpected json: %s", string(data))
			}
		})
	}
}

func TestChatToolChoice_Unmarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantMode  *ToolChoiceMode
		wantFunc  *ChatFunctionName
		wantErr   bool
	}{
		{
			name:      "function",
			unmarshal: `{"function":{"name":"do"}}`,
			wantFunc:  &ChatFunctionName{Name: "do"},
		},
		{
			name:      "null",
			unmarshal: `null`,
		},
		{
			name:      "mode string",
			unmarshal: `"auto"`,
			wantMode:  testutils.Ptr(ToolChoiceMode("auto")),
		},
		{
			name:      "empty mode",
			unmarshal: `""`,
			wantErr:   true,
		},
		{
			name:      "empty object",
			unmarshal: `{}`,
			wantErr:   true,
		},
		{
			name:      "function null",
			unmarshal: `{"function":null}`,
			wantErr:   true,
		},
		{
			name:      "array payload",
			unmarshal: `[]`,
			wantErr:   true,
		},
		{
			name:      "empty payload",
			unmarshal: ``,
			wantErr:   true,
		},
		{
			name:      "mode in object",
			unmarshal: `{"mode":"auto"}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatToolChoice
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
			if tc.wantMode != nil {
				if got.Mode == nil || *got.Mode != *tc.wantMode {
					t.Fatalf("unexpected mode: %+v", got.Mode)
				}
			}
			if tc.wantFunc != nil {
				if got.Function == nil || got.Function.Name != tc.wantFunc.Name {
					t.Fatalf("unexpected function: %+v", got.Function)
				}
			}
		})
	}
}

func TestChatResponseFormat(t *testing.T) {
	t.Parallel()

	providerType := ResponseFormatType("provider-format")

	cases := []struct {
		name    string
		value   ChatResponseFormat
		wantErr bool
	}{
		{
			name: "json_schema",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeJSONSchema,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
		},
		{
			name:  "provider type",
			value: ChatResponseFormat{Type: providerType},
		},
		{
			name:  "text",
			value: ChatResponseFormat{Type: ResponseFormatTypeText},
		},
		{
			name:  "json_object",
			value: ChatResponseFormat{Type: ResponseFormatTypeJSONObject},
		},
		{
			name:    "json_schema missing",
			value:   ChatResponseFormat{Type: ResponseFormatTypeJSONSchema},
			wantErr: true,
		},
		{
			name: "other with json_schema",
			value: ChatResponseFormat{
				Type:       providerType,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr: true,
		},
		{
			name: "text with json_schema",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeText,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr: true,
		},
		{
			name: "json_object with json_schema",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeJSONObject,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr: true,
		},
		{
			name:    "empty type",
			value:   ChatResponseFormat{},
			wantErr: true,
		},
		{
			name: "json_schema empty name",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeJSONSchema,
				JSONSchema: &ChatJSONSchemaConfig{Name: ""},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatResponseFormat_UnmarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "json_schema null",
			unmarshal: `{"type":"json_schema","json_schema":null}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatResponseFormat
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatResponseFormat_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantType  ResponseFormatType
	}{
		{
			name:      "text",
			unmarshal: `{"type":"text"}`,
			wantType:  ResponseFormatTypeText,
		},
		{
			name:      "provider",
			unmarshal: `{"type":"provider-format"}`,
			wantType:  ResponseFormatType("provider-format"),
		},
		{
			name:      "json_object",
			unmarshal: `{"type":"json_object"}`,
			wantType:  ResponseFormatTypeJSONObject,
		},
		{
			name:      "json_schema",
			unmarshal: `{"type":"json_schema","json_schema":{"name":"n"}}`,
			wantType:  ResponseFormatTypeJSONSchema,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatResponseFormat
			testutils.RequireNoError(t, json.Unmarshal([]byte(tc.unmarshal), &got))
			if got.Type != tc.wantType {
				t.Fatalf("unexpected type: %v", got.Type)
			}
			if got.Type == ResponseFormatTypeJSONSchema &&
				(got.JSONSchema == nil || got.JSONSchema.Name != "n") {
				t.Fatalf("unexpected json_schema: %+v", got.JSONSchema)
			}
		})
	}
}

func TestChatCompletionMessage(t *testing.T) {
	t.Parallel()

	text := "ok"

	cases := []struct {
		name    string
		value   ChatCompletionMessage
		wantErr bool
	}{
		{
			name:  "content",
			value: ChatCompletionMessage{Role: "assistant", Content: &text},
		},
		{
			name: "tool_calls",
			value: ChatCompletionMessage{Role: "assistant", ToolCalls: []ChatToolCallOutput{{
				ID:       "id",
				Type:     "function",
				Function: ChatFunctionCall{Name: "fn", Arguments: "{}"},
			}}},
		},
		{
			name: "both",
			value: ChatCompletionMessage{
				Role:    "assistant",
				Content: &text,
				ToolCalls: []ChatToolCallOutput{
					{
						ID:       "id",
						Type:     "function",
						Function: ChatFunctionCall{Name: "fn", Arguments: "{}"},
					},
				},
			},
			wantErr: true,
		},
		{
			name:    "neither",
			value:   ChatCompletionMessage{Role: "assistant"},
			wantErr: true,
		},
		{
			name:    "tool_call_id without content",
			value:   ChatCompletionMessage{Role: "assistant", ToolCallID: testutils.Ptr("id")},
			wantErr: true,
		},
		{
			name: "content with tool_call_id",
			value: ChatCompletionMessage{
				Role:       "assistant",
				Content:    &text,
				ToolCallID: testutils.Ptr("id"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatCompletionMessage_UnmarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "both content and tool_calls",
			unmarshal: `{"role":"assistant","content":"hi","tool_calls":[{"id":"id","type":"function","function":{"name":"fn","arguments":"{}"}}]}`,
			wantErr:   true,
		},
		{
			name:      "neither content nor tool_calls",
			unmarshal: `{"role":"assistant"}`,
			wantErr:   true,
		},
		{
			name:      "tool_call_id without content",
			unmarshal: `{"role":"assistant","tool_call_id":"id"}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatCompletionMessage
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatStreamDelta(t *testing.T) {
	t.Parallel()

	text := "ok"

	cases := []struct {
		name    string
		value   ChatStreamDelta
		wantErr bool
	}{
		{
			name:  "content",
			value: ChatStreamDelta{Content: &text},
		},
		{
			name:  "role only",
			value: ChatStreamDelta{Role: testutils.Ptr("assistant")},
		},
		{
			name: "tool_calls",
			value: ChatStreamDelta{ToolCalls: []ChatStreamToolCall{{
				ID:       "id",
				Type:     "function",
				Index:    0,
				Function: ChatStreamFunction{Name: "fn", Arguments: "{}"},
			}}},
		},
		{
			name: "tool_calls with content",
			value: ChatStreamDelta{
				Content: &text,
				ToolCalls: []ChatStreamToolCall{
					{
						ID:       "id",
						Type:     "function",
						Index:    0,
						Function: ChatStreamFunction{Name: "fn", Arguments: "{}"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "tool_calls with tool_call_id",
			value: ChatStreamDelta{
				ToolCallID: testutils.Ptr("id"),
				ToolCalls: []ChatStreamToolCall{
					{
						ID:       "id",
						Type:     "function",
						Index:    0,
						Function: ChatStreamFunction{Name: "fn", Arguments: "{}"},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatStreamDelta_UnmarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "tool_calls with content",
			unmarshal: `{"content":"hi","tool_calls":[{"id":"id","type":"function","index":0,"function":{"name":"fn","arguments":"{}"}}]}`,
			wantErr:   true,
		},
		{
			name:      "tool_calls with tool_call_id",
			unmarshal: `{"tool_call_id":"id","tool_calls":[{"id":"id","type":"function","index":0,"function":{"name":"fn","arguments":"{}"}}]}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatStreamDelta
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatStreamDelta_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		unmarshal     string
		wantRole      *string
		wantContent   *string
		wantToolCalls int
	}{
		{
			name:      "role only",
			unmarshal: `{"role":"assistant"}`,
			wantRole:  testutils.Ptr("assistant"),
		},
		{
			name:        "content only",
			unmarshal:   `{"content":"hi"}`,
			wantContent: testutils.Ptr("hi"),
		},
		{
			name:          "tool_calls only",
			unmarshal:     `{"tool_calls":[{"id":"id","type":"function","index":0,"function":{"name":"fn","arguments":"{}"}}]}`,
			wantToolCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatStreamDelta
			testutils.RequireNoError(t, json.Unmarshal([]byte(tc.unmarshal), &got))
			if tc.wantRole != nil {
				if got.Role == nil || *got.Role != *tc.wantRole {
					t.Fatalf("unexpected role: %+v", got.Role)
				}
				testutils.RequireNil(t, got.Content)
				testutils.RequireNil(t, got.ToolCallID)
				testutils.RequireNil(t, got.ToolCalls)
			}
			if tc.wantContent != nil {
				if got.Content == nil || *got.Content != *tc.wantContent {
					t.Fatalf("unexpected content: %+v", got.Content)
				}
				testutils.RequireNil(t, got.Role)
				testutils.RequireNil(t, got.ToolCallID)
				testutils.RequireNil(t, got.ToolCalls)
			}
			if tc.wantToolCalls > 0 {
				if len(got.ToolCalls) != tc.wantToolCalls {
					t.Fatalf("unexpected tool calls: %+v", got.ToolCalls)
				}
				call := got.ToolCalls[0]
				if call.ID != "id" || call.Type != "function" || call.Index != 0 ||
					call.Function.Name != "fn" ||
					call.Function.Arguments != "{}" {
					t.Fatalf("unexpected tool call: %+v", call)
				}
				testutils.RequireNil(t, got.Content)
				testutils.RequireNil(t, got.ToolCallID)
			}
		})
	}
}

func TestChatToolCall_TypeMustBeSet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "type empty",
			data:    []byte(`{"id":"id","type":"","function":{"name":"fn"}}`),
			wantErr: true,
		},
		{
			name:    "type missing",
			data:    []byte(`{"id":"id","function":{"name":"fn"}}`),
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out ChatToolCall
			err := json.Unmarshal(tc.data, &out)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatToolCall_MarshalTypeMissing(t *testing.T) {
	t.Parallel()

	value := ChatToolCall{
		ID:       "id",
		Function: ChatFunctionDefinition{Name: "fn"},
	}
	_, err := json.Marshal(value)
	testutils.RequireError(t, err)
}

func TestChatToolCall_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":"id","type":"function","function":{"name":"fn"}}`)
	var got ChatToolCall
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.ID != "id" || got.Type != "function" || got.Function.Name != "fn" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestChatTool_TypeMustBeSet(t *testing.T) {
	t.Parallel()

	data := []byte(`{"function":{"name":"fn"}}`)
	var out ChatTool
	err := json.Unmarshal(data, &out)
	testutils.RequireError(t, err)
}

func TestChatTool_MarshalTypeMissing(t *testing.T) {
	t.Parallel()

	value := ChatTool{
		Function: ChatFunctionDefinition{Name: "fn"},
	}
	_, err := json.Marshal(value)
	testutils.RequireError(t, err)
}

func TestChatTool_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"type":"function","function":{"name":"fn"}}`)
	var got ChatTool
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Type != "function" || got.Function.Name != "fn" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestChatToolCallOutput_TypeMustBeSet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "type empty",
			data:    []byte(`{"id":"id","type":"","function":{"name":"fn","arguments":"{}"}}`),
			wantErr: true,
		},
		{
			name:    "type missing",
			data:    []byte(`{"id":"id","function":{"name":"fn","arguments":"{}"}}`),
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out ChatToolCallOutput
			err := json.Unmarshal(tc.data, &out)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatToolCallOutput_MarshalTypeMissing(t *testing.T) {
	t.Parallel()

	value := ChatToolCallOutput{
		ID:       "id",
		Function: ChatFunctionCall{Name: "fn", Arguments: "{}"},
	}
	_, err := json.Marshal(value)
	testutils.RequireError(t, err)
}

func TestChatToolCallOutput_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":"id","type":"function","function":{"name":"fn","arguments":"{}"}}`)
	var got ChatToolCallOutput
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.ID != "id" || got.Type != "function" || got.Function.Name != "fn" ||
		got.Function.Arguments != "{}" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestChatStreamToolCall_TypeMustBeSet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name: "type empty",
			data: []byte(
				`{"id":"id","type":"","index":0,"function":{"name":"fn","arguments":"{}"}}`,
			),
			wantErr: true,
		},
		{
			name:    "type missing",
			data:    []byte(`{"id":"id","index":0,"function":{"name":"fn","arguments":"{}"}}`),
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out ChatStreamToolCall
			err := json.Unmarshal(tc.data, &out)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatStreamToolCall_MarshalTypeMissing(t *testing.T) {
	t.Parallel()

	value := ChatStreamToolCall{
		ID:       "id",
		Index:    0,
		Function: ChatStreamFunction{Name: "fn", Arguments: "{}"},
	}
	_, err := json.Marshal(value)
	testutils.RequireError(t, err)
}

func TestChatStreamToolCall_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(
		`{"id":"id","type":"function","index":0,"function":{"name":"fn","arguments":"{}"}}`,
	)
	var got ChatStreamToolCall
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.ID != "id" || got.Type != "function" || got.Index != 0 || got.Function.Name != "fn" {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestChatRequest_UnmarshalTypeValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "max_tokens string",
			unmarshal: `{"messages":[{"role":"user","content":"hi"}],"max_tokens":"no"}`,
			wantErr:   true,
		},
		{
			name:      "logprobs string",
			unmarshal: `{"messages":[{"role":"user","content":"hi"}],"logprobs":"false"}`,
			wantErr:   true,
		},
		{
			name:      "top_logprobs string",
			unmarshal: `{"messages":[{"role":"user","content":"hi"}],"top_logprobs":"5"}`,
			wantErr:   true,
		},
		{
			name:      "temperature object",
			unmarshal: `{"messages":[{"role":"user","content":"hi"}],"temperature":{}}`,
			wantErr:   true,
		},
		{
			name:      "presence_penalty bool",
			unmarshal: `{"messages":[{"role":"user","content":"hi"}],"presence_penalty":false}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatRequest
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatResponse_UnmarshalTypeValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		unmarshal string
		wantErr   bool
	}{
		{
			name:      "usage total_tokens string",
			unmarshal: `{"id":"id","created":1,"model":"m","system_fingerprint":"s","choices":[],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":"3"}}`,
			wantErr:   true,
		},
		{
			name:      "choices not array",
			unmarshal: `{"id":"id","created":1,"model":"m","system_fingerprint":"s","choices":{},"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}}`,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got ChatResponse
			err := json.Unmarshal([]byte(tc.unmarshal), &got)
			if testutils.AssertError(t, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestChatResponse_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(
		`{"id":"id","created":1,"model":"m","system_fingerprint":"s","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"hi"}}],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}}`,
	)
	var got ChatResponse
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.ID != "id" || got.Model != "m" || len(got.Choices) != 1 {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Choices[0].Message.Content == nil || *got.Choices[0].Message.Content != "hi" {
		t.Fatalf("unexpected message content: %+v", got.Choices[0].Message.Content)
	}
}

func TestChatStreamResponse_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(
		`{"id":"id","created":1,"model":"m","system_fingerprint":"s","choices":[{"delta":{"role":"assistant"},"index":0}]}`,
	)
	var got ChatStreamResponse
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.ID != "id" || len(got.Choices) != 1 {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Choices[0].Delta.Role == nil || *got.Choices[0].Delta.Role != "assistant" {
		t.Fatalf("unexpected delta role: %+v", got.Choices[0].Delta.Role)
	}
}

func TestChatLogProbs_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(
		`{"content":[{"token":"t","logprob":0.1,"top_logprobs":[{"token":"t","logprob":0.1}]}]}`,
	)
	var got ChatLogProbs
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if len(got.Content) != 1 || got.Content[0].Token != "t" {
		t.Fatalf("unexpected logprobs: %+v", got)
	}
}

func TestChatUsage_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}`)
	var got ChatUsage
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.PromptTokens != 1 || got.CompletionTokens != 2 || got.TotalTokens != 3 {
		t.Fatalf("unexpected usage: %+v", got)
	}
}

func TestChatFunctionDefinition_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"name":"fn","description":"desc","parameters":{"type":"object"}}`)
	var got ChatFunctionDefinition
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Name != "fn" || got.Description == nil || *got.Description != "desc" {
		t.Fatalf("unexpected function definition: %+v", got)
	}
}

func TestChatFunctionCall_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"name":"fn","arguments":"{}"}`)
	var got ChatFunctionCall
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Name != "fn" || got.Arguments != "{}" {
		t.Fatalf("unexpected function call: %+v", got)
	}
}

func TestChatFunctionName_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"name":"fn"}`)
	var got ChatFunctionName
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Name != "fn" {
		t.Fatalf("unexpected function name: %+v", got)
	}
}

func TestChatStreamFunction_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"name":"fn","arguments":"{}"}`)
	var got ChatStreamFunction
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Name != "fn" || got.Arguments != "{}" {
		t.Fatalf("unexpected stream function: %+v", got)
	}
}

func TestChatChoice_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(
		`{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"hi"}}`,
	)
	var got ChatChoice
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.FinishReason != "stop" || got.Message.Content == nil || *got.Message.Content != "hi" {
		t.Fatalf("unexpected choice: %+v", got)
	}
}

func TestChatStreamChoice_UnmarshalSuccess(t *testing.T) {
	t.Parallel()

	data := []byte(`{"delta":{"role":"assistant"},"index":0}`)
	var got ChatStreamChoice
	testutils.RequireNoError(t, json.Unmarshal(data, &got))
	if got.Index != 0 || got.Delta.Role == nil || *got.Delta.Role != "assistant" {
		t.Fatalf("unexpected stream choice: %+v", got)
	}
}
