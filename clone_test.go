//go:build !integration

package hfgo

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatRequestClone_Deep(t *testing.T) {
	t.Parallel()

	schema := json.RawMessage(`{"type":"object"}`)
	req := &ChatRequest{
		Model: testutils.Ptr("m"),
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: testutils.Ptr("hello")}},
		},
		Stop: []string{"x"},
		Tools: []ChatTool{
			{Type: "function", Function: ChatFunctionDefinition{Name: "f", Parameters: schema}},
		},
		ResponseFormat: &ChatResponseFormat{
			Type:       ResponseFormatTypeJSONSchema,
			JSONSchema: &ChatJSONSchemaConfig{Name: "n", Schema: schema},
		},
		ToolChoice: &ChatToolChoice{Mode: testutils.Ptr(ToolChoiceModeAuto)},
	}

	cloned := req.Clone()

	*cloned.Messages[0].Content.Text = "changed"
	cloned.Messages[0].Role = "assistant"
	cloned.Stop[0] = "y"
	cloned.Tools[0].Function.Parameters = json.RawMessage(`{"type":"object"}zzz`)
	cloned.ResponseFormat.JSONSchema.Schema = json.RawMessage(`{"type":"object"}xxx`)
	*cloned.ToolChoice.Mode = ToolChoiceModeNone

	require.Equal(t, "hello", *req.Messages[0].Content.Text)
	require.Equal(t, "user", req.Messages[0].Role)
	require.Equal(t, []string{"x"}, req.Stop)
	require.JSONEq(t, `{"type":"object"}`, string(req.Tools[0].Function.Parameters))
	require.JSONEq(t, `{"type":"object"}`, string(req.ResponseFormat.JSONSchema.Schema))
	require.Equal(t, ToolChoiceModeAuto, *req.ToolChoice.Mode)
}

func TestRequestClones_NestedIndependence(t *testing.T) {
	t.Parallel()

	origParams := SummarizationParameters{
		CleanUpTokenizationSpaces: testutils.Ptr(false),
		Truncation:                testutils.Ptr(SummarizationTruncationLongestFirst),
		GenerateParameters:        map[string]any{"max_new_tokens": 5},
	}
	single := SummarizationRequest{Input: "a", Parameters: &origParams}
	batch := SummarizationBatchRequest{Inputs: []string{"a", "b"}, Parameters: &origParams}

	singleCloned := single.Clone()
	batchCloned := batch.Clone()

	*singleCloned.Parameters.Truncation = SummarizationTruncationOnlyFirst
	singleCloned.Parameters.GenerateParameters["max_new_tokens"] = 99
	*singleCloned.Parameters.CleanUpTokenizationSpaces = true
	batchCloned.Inputs[0] = "zzz"
	*batchCloned.Parameters.Truncation = SummarizationTruncationDoNotTruncate

	require.Equal(t, SummarizationTruncationLongestFirst, *origParams.Truncation)
	require.Equal(t, "a", single.Input)
	require.Equal(t, []string{"a", "b"}, batch.Inputs)
	require.False(t, *origParams.CleanUpTokenizationSpaces)
	require.Equal(t, 5, origParams.GenerateParameters["max_new_tokens"])
}

func TestChatRequest_SequentialReuseDoesNotMutate(t *testing.T) {
	t.Parallel()

	text := "hi"
	req := ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	call := func() {
		mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
		client := NewClient(
			WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
			WithModel("default-model"),
		)
		_, err := client.Chat(req)
		require.NoError(t, err)
		require.Nil(t, req.Model)
	}

	call()
	req.Stop = []string{"x"}
	call()
}

func TestChatRequestClone_ConcurrentInvocation(t *testing.T) {
	t.Parallel()

	text := "hi"
	req := ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	const workers = 16
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
			client := NewClient(
				WithHTTPClientFactory(
					func() http.Client { return testutils.NewMockHTTPClient(mt) },
				),
				WithModel("default-model"),
			)
			_, err := client.Chat(req.Clone())
			assert.NoError(t, err)
		})
	}
	wg.Wait()
}
