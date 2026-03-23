package hfgo

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestChatMessageContent_Marshal(t *testing.T) {
	t.Parallel()

	text := "hello"
	imgURL := "https://example.com/image.png"

	cases := []struct {
		name        string
		value       ChatMessageContent
		wantJSON    string
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
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
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
			if tc.wantJSON != "" && string(data) != tc.wantJSON {
				t.Fatalf("unexpected json: %s", string(data))
			}
		})
	}
}

func TestChatMessageChunk_Validation(t *testing.T) {
	t.Parallel()

	text := "hello"

	cases := []struct {
		name        string
		value       ChatMessageChunk
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
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
			name:        "missing text",
			value:       ChatMessageChunk{Type: MessageChunkTypeText},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "missing image",
			value:       ChatMessageChunk{Type: MessageChunkTypeImageURL},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "missing type",
			value:       ChatMessageChunk{},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "invalid type",
			value:       ChatMessageChunk{Type: MessageChunkType("other")},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "text with image_url",
			value: ChatMessageChunk{
				Type:     MessageChunkTypeText,
				Text:     &text,
				ImageURL: &ChatImageURL{URL: "x"},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "image_url with text",
			value: ChatMessageChunk{
				Type:     MessageChunkTypeImageURL,
				Text:     &text,
				ImageURL: &ChatImageURL{URL: "x"},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
		})
	}
}

func TestChatImageURL_Validation(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(ChatImageURL{})
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
}

func TestChatMessage_Validation(t *testing.T) {
	t.Parallel()

	text := "hi"

	cases := []struct {
		name        string
		value       ChatMessage
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
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
				Function: ChatFunctionCall{Name: "do", Arguments: "{}"},
			}}},
		},
		{
			name: "both",
			value: ChatMessage{
				Role:    "assistant",
				Content: ChatMessageContent{Text: &text},
				ToolCalls: []ChatToolCall{
					{
						ID:       "id",
						Type:     "function",
						Function: ChatFunctionCall{Name: "do", Arguments: "{}"},
					},
				},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "neither",
			value:       ChatMessage{Role: "assistant"},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "missing role",
			value:       ChatMessage{Content: ChatMessageContent{Text: &text}},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
		})
	}
}

func TestChatRequest_MarshalSuccess(t *testing.T) {
	t.Parallel()

	text := "hi"
	imgURL := "https://example.com/image.png"
	model := "model"
	req := ChatRequest{
		Model: &model,
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
	require.NoError(t, err)

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unexpected json: %v", err)
	}
	if got["model"] != model {
		t.Fatalf("unexpected model: %#v", got["model"])
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

func TestChatRequest_MarshalValidation(t *testing.T) {
	t.Parallel()

	text := "hi"
	model := "model"

	cases := []*struct {
		name        string
		value       ChatRequest
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
	}{
		{
			name: "missing model",
			value: ChatRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: ChatMessageContent{Text: &text}},
				},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "empty model",
			value: ChatRequest{
				Model: testutils.Ptr(""),
				Messages: []ChatMessage{
					{Role: "user", Content: ChatMessageContent{Text: &text}},
				},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "missing messages",
			value: ChatRequest{
				Model: &model,
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "valid request",
			value: ChatRequest{
				Model: &model,
				Messages: []ChatMessage{
					{Role: "user", Content: ChatMessageContent{Text: &text}},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestChatToolChoice_Marshal(t *testing.T) {
	t.Parallel()

	mode := ToolChoiceMode("provider-mode")

	cases := []struct {
		name        string
		value       ChatToolChoice
		wantJSON    string
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
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
			name:        "both",
			value:       ChatToolChoice{Mode: &mode, Function: &ChatFunctionName{Name: "do"}},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "empty mode",
			value:       ChatToolChoice{Mode: testutils.Ptr(ToolChoiceMode(""))},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
			if tc.wantJSON != "" && string(data) != tc.wantJSON {
				t.Fatalf("unexpected json: %s", string(data))
			}
		})
	}
}

func TestChatResponseFormat(t *testing.T) {
	t.Parallel()

	providerType := ResponseFormatType("provider-format")

	cases := []struct {
		name        string
		value       ChatResponseFormat
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
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
			name:        "json_schema missing",
			value:       ChatResponseFormat{Type: ResponseFormatTypeJSONSchema},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "other with json_schema",
			value: ChatResponseFormat{
				Type:       providerType,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "text with json_schema",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeText,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "json_object with json_schema",
			value: ChatResponseFormat{
				Type:       ResponseFormatTypeJSONObject,
				JSONSchema: &ChatJSONSchemaConfig{Name: "n"},
			},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name:        "empty type",
			value:       ChatResponseFormat{},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
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
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
		})
	}
}

func TestChatFunctionDefinition_Validation(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(ChatFunctionDefinition{})
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
}

func TestChatFunctionName_MarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		value       ChatFunctionName
		wantJSON    string
		wantErr     bool
		wantErrKind hferrors.SDKErrorKind
	}{
		{
			name:     "success",
			value:    ChatFunctionName{Name: "fn"},
			wantJSON: `{"name":"fn"}`,
		},
		{
			name:        "missing name",
			value:       ChatFunctionName{},
			wantErr:     true,
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantJSON, string(data))
		})
	}
}

func TestChatTool_MarshalTypeMissing(t *testing.T) {
	t.Parallel()

	value := ChatTool{
		Function: ChatFunctionDefinition{Name: "fn"},
	}
	_, err := json.Marshal(value)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
}

func TestChatToolCall_MarshalValidation(t *testing.T) {
	t.Parallel()

	cases := []toolCallMarshalCase{
		{
			name: "missing type",
			value: ChatToolCall{
				ID:       "id",
				Function: ChatFunctionCall{Name: "fn", Arguments: "{}"},
			},
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "missing id",
			value: ChatToolCall{
				Type:     "function",
				Function: ChatFunctionCall{Name: "fn", Arguments: "{}"},
			},
			wantErrKind: hferrors.SDKErrorKindConfiguration,
		},
		{
			name: "missing function name",
			value: ChatToolCall{
				ID:       "id",
				Type:     "function",
				Function: ChatFunctionCall{Arguments: "{}"},
			},
			wantErrKind: hferrors.SDKErrorKindValidation,
		},
		{
			name: "missing function arguments",
			value: ChatToolCall{
				ID:       "id",
				Type:     "function",
				Function: ChatFunctionCall{Name: "fn"},
			},
			wantErrKind: hferrors.SDKErrorKindValidation,
		},
	}

	runToolCallMarshalTests(t, cases)
}
