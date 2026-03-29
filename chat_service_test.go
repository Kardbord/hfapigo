//go:build !integration

package hfgo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/request"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

const chatServiceResponseBody = `{"id":"id","created":1,"model":"m","system_fingerprint":"s","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"hi"}}],"usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}}`

func TestChatService_Complete_ModelSelection(t *testing.T) {
	t.Parallel()

	text := "hi"

	cases := []struct {
		name        string
		clientModel string
		optsModel   string
		reqModel    *string
		wantModel   string
	}{
		{
			name:        "uses client model when request and opt model missing",
			clientModel: "default-model",
			wantModel:   "default-model",
		},
		{
			name:        "uses opt model when request missing",
			clientModel: "default-model",
			optsModel:   "explicit-model",
			wantModel:   "explicit-model",
		},
		{
			name:        "respects request model",
			clientModel: "default-model",
			optsModel:   "opts-model",
			reqModel:    testutils.Ptr("explicit-model"),
			wantModel:   "explicit-model",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
			opts := request.NewOptions().
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
				WithModel(tc.clientModel)
			svc := newChatService(opts)
			req := &ChatRequest{
				Model: tc.reqModel,
				Messages: []ChatMessage{
					{Role: "user", Content: ChatMessageContent{Text: &text}},
				},
			}

			var err error
			if tc.optsModel != "" {
				_, err = svc.Complete(req, WithModel(tc.optsModel))
			} else {
				_, err = svc.Complete(req)
			}

			require.NoError(t, err)

			require.NotNil(t, mt.LastRequest)
			require.Equal(t, EndpointChatCompletion, mt.LastRequest.URL.Path)

			body, err := io.ReadAll(mt.LastRequest.Body)
			require.NoError(t, err)
			_ = mt.LastRequest.Body.Close()

			var got map[string]any
			require.NoError(t, json.Unmarshal(body, &got))
			require.Equal(t, tc.wantModel, got["model"])

			if tc.reqModel == nil {
				require.Nil(t, req.Model)
			} else {
				require.NotNil(t, req.Model)
				require.Equal(t, *tc.reqModel, *req.Model)
			}
		})
	}
}

func TestChatService_Complete_ModelValidation(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newChatService(opts)

	text := "hi"
	req := &ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	_, err := svc.Complete(req)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
	require.Nil(t, mt.LastRequest)
}

func TestChatService_Complete_NilRequest(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newChatService(opts)

	_, err := svc.Complete(nil)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
	require.Nil(t, mt.LastRequest)
}

func TestChatService_Complete_StreamNotAllowed(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newChatService(opts)

	text := "hi"
	stream := true
	req := &ChatRequest{
		Stream: &stream,
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	_, err := svc.Complete(req)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
	require.Nil(t, mt.LastRequest)
}

func TestChatService_CompleteStream_Success(t *testing.T) {
	t.Parallel()

	body := "data: {\"id\":\"id\",\"created\":1,\"model\":\"stream-model\",\"system_fingerprint\":\"sig\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"hi\"}}]}\n\n" +
		"data: [DONE]\n\n"
	mt := testutils.NewMockTransport(http.StatusOK, body, nil)
	mt.Response.Header.Set("Content-Type", "text/event-stream")

	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("default-model")
	svc := newChatService(opts)

	text := "hi"
	req := &ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	stream, err := svc.CompleteStream(req)
	require.NoError(t, err)
	defer func() { _ = stream.Close() }()

	chunk, err := stream.Recv(context.Background())
	require.NoError(t, err)
	require.Equal(t, "stream-model", chunk.Model)
	require.Len(t, chunk.Choices, 1)
	require.NotNil(t, chunk.Choices[0].Delta.Content)
	require.Equal(t, "hi", *chunk.Choices[0].Delta.Content)

	_, err = stream.Recv(context.Background())
	require.ErrorIs(t, err, io.EOF)

	require.NotNil(t, mt.LastRequest)
	bodyBytes, err := io.ReadAll(mt.LastRequest.Body)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(bodyBytes, &payload))
	require.Equal(t, true, payload["stream"])
	require.Equal(t, "default-model", payload["model"])
}

func TestChatStream_Recv_MergesToolCallMetadata(t *testing.T) {
	t.Parallel()

	body := strings.Join([]string{
		`data: {"id":"id","created":1,"model":"stream-model","system_fingerprint":"sig","choices":[{"index":0,"delta":{"tool_calls":[{"id":"call_0","type":"function","index":0,"function":{"name":"fn","arguments":""}}]}}]}`,
		``,
		`data: {"id":"id","created":1,"model":"stream-model","system_fingerprint":"sig","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"foo\":1}"}}]}}]}`,
		``,
		"data: [DONE]",
		``,
		``,
	}, "\n")
	assertToolCallStream(t, body, func(chunks []ChatStreamResponse) {
		require.Len(t, chunks, 2)
		first, second := chunks[0], chunks[1]
		require.Equal(t, "call_0", first.Choices[0].Delta.ToolCalls[0].ID)
		require.Equal(t, "call_0", second.Choices[0].Delta.ToolCalls[0].ID)
		require.Equal(t, `{"foo":1}`, second.Choices[0].Delta.ToolCalls[0].Function.Arguments)
	})
}

func TestChatStream_Recv_MergesAcrossChoices(t *testing.T) {
	t.Parallel()

	body := strings.Join([]string{
		`data: {"id":"id","created":1,"model":"stream-model","system_fingerprint":"sig","choices":[{"index":0,"delta":{"tool_calls":[{"id":"call_0","type":"function","index":0,"function":{"name":"fn","arguments":""}}]}},{"index":1,"delta":{"tool_calls":[{"id":"call_1","type":"function","index":0,"function":{"name":"fn2","arguments":""}}]}}]}`,
		``,
		`data: {"id":"id","created":1,"model":"stream-model","system_fingerprint":"sig","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"foo\":1}"}}]}},{"index":1,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"bar\":2}"}}]}}]}`,
		``,
		"data: [DONE]",
		``,
		``,
	}, "\n")

	assertToolCallStream(t, body, func(chunks []ChatStreamResponse) {
		require.Len(t, chunks, 2)
		first, second := chunks[0], chunks[1]
		require.Equal(t, "call_0", first.Choices[0].Delta.ToolCalls[0].ID)
		require.Equal(t, "call_1", first.Choices[1].Delta.ToolCalls[0].ID)
		require.Equal(t, "call_0", second.Choices[0].Delta.ToolCalls[0].ID)
		require.Equal(t, "call_1", second.Choices[1].Delta.ToolCalls[0].ID)
		require.Equal(t, `{"foo":1}`, second.Choices[0].Delta.ToolCalls[0].Function.Arguments)
		require.Equal(t, `{"bar":2}`, second.Choices[1].Delta.ToolCalls[0].Function.Arguments)
	})
}

func TestChatStream_Recv_InvalidJSONError(t *testing.T) {
	t.Parallel()

	body := strings.Join([]string{
		`data: {"id":"id","created":1,"model":"stream-model","system_fingerprint":"sig","choices":[{"index":0,"delta":{"tool_calls":[{"id":"call_0","type":"function","index":0,"function":{"name":"fn","arguments":""}}]}}]}`,
		``,
		"data: {not json}",
		``,
	}, "\n")

	mt := testutils.NewMockTransport(http.StatusOK, body, nil)
	mt.Response.Header.Set("Content-Type", "text/event-stream")

	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("default-model")
	svc := newChatService(opts)

	text := "hi"
	req := &ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	stream, err := svc.CompleteStream(req)
	require.NoError(t, err)
	defer func() { _ = stream.Close() }()

	first, err := stream.Recv(context.Background())
	require.NoError(t, err)
	require.Equal(t, "call_0", first.Choices[0].Delta.ToolCalls[0].ID)

	_, err = stream.Recv(context.Background())
	require.Error(t, err)
}

// assertToolCallStream streams the provided SSE body through ChatStream and
// passes all decoded chunks to the supplied assertion callback.
func assertToolCallStream(
	t *testing.T,
	body string,
	assertions func(chunks []ChatStreamResponse),
) {
	t.Helper()

	mt := testutils.NewMockTransport(http.StatusOK, body, nil)
	mt.Response.Header.Set("Content-Type", "text/event-stream")

	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("default-model")
	svc := newChatService(opts)

	text := "hi"
	req := &ChatRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: ChatMessageContent{Text: &text}},
		},
	}

	stream, err := svc.CompleteStream(req)
	require.NoError(t, err)
	defer func() { _ = stream.Close() }()

	var chunks []ChatStreamResponse
	for {
		chunk, err := stream.Recv(context.Background())
		if err != nil {
			require.ErrorIs(t, err, io.EOF)

			break
		}
		chunks = append(chunks, chunk)
	}
	assertions(chunks)
}

func TestChatService_CompleteStream_NilRequest(t *testing.T) {
	t.Parallel()

	mt := testutils.NewMockTransport(http.StatusOK, "", nil)
	mt.Response.Header.Set("Content-Type", "text/event-stream")
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newChatService(opts)

	_, err := svc.CompleteStream(nil)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)
	require.Nil(t, mt.LastRequest)
}

func TestApplyProvider(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		model       *string
		provider    string
		wantModel   *string
		description string
	}{
		{
			name:        "applies provider to model without provider",
			model:       testutils.Ptr("mistral-7b"),
			provider:    "huggingface",
			wantModel:   testutils.Ptr("mistral-7b:huggingface"),
			description: "model + provider → model:provider",
		},
		{
			name:        "ignores provider when model already has provider",
			model:       testutils.Ptr("mistral-7b:mistral"),
			provider:    "huggingface",
			wantModel:   testutils.Ptr("mistral-7b:mistral"),
			description: "model:provider + different provider → unchanged",
		},
		{
			name:        "returns nil model when model is nil",
			model:       nil,
			provider:    "huggingface",
			wantModel:   nil,
			description: "nil model → nil",
		},
		{
			name:        "returns nil model when model is empty string",
			model:       testutils.Ptr(""),
			provider:    "huggingface",
			wantModel:   testutils.Ptr(""),
			description: "empty model → empty",
		},
		{
			name:        "returns model unchanged when provider is empty",
			model:       testutils.Ptr("mistral-7b"),
			provider:    "",
			wantModel:   testutils.Ptr("mistral-7b"),
			description: "model + empty provider → model unchanged",
		},
		{
			name:        "handles provider with special characters",
			model:       testutils.Ptr("mistral-7b"),
			provider:    "provider-name",
			wantModel:   testutils.Ptr("mistral-7b:provider-name"),
			description: "provider with hyphens",
		},
		{
			name:        "handles provider with underscores",
			model:       testutils.Ptr("mistral-7b"),
			provider:    "provider_name",
			wantModel:   testutils.Ptr("mistral-7b:provider_name"),
			description: "provider with underscores",
		},
		{
			name:        "handles provider with dots",
			model:       testutils.Ptr("mistral-7b"),
			provider:    "provider.com",
			wantModel:   testutils.Ptr("mistral-7b:provider.com"),
			description: "provider with dots",
		},
		{
			name:        "ignores provider when model has multiple colons",
			model:       testutils.Ptr("org:model:variant"),
			provider:    "huggingface",
			wantModel:   testutils.Ptr("org:model:variant"),
			description: "model with multiple colons",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := applyProvider(tc.model, tc.provider)
			// Verify result matches expectation
			if tc.wantModel == nil {
				require.Nil(t, got, tc.description)

				return
			}
			require.NotNil(t, got, tc.description)
			require.Equal(t, *tc.wantModel, *got, tc.description)
		})
	}
}

func TestResolveModel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		reqModel       *string
		clientModel    string
		clientProvider string
		optsModel      string
		optsProvider   string
		wantModel      *string
		description    string
	}{
		{
			name:        "uses request model when provided",
			reqModel:    testutils.Ptr("request-model"),
			clientModel: "client-model",
			wantModel:   testutils.Ptr("request-model"),
			description: "request model takes precedence",
		},
		{
			name:        "uses options model when request model is nil",
			clientModel: "client-model",
			optsModel:   "opts-model",
			wantModel:   testutils.Ptr("opts-model"),
			description: "options model used when request is nil",
		},
		{
			name:        "uses client model when request and options are nil",
			clientModel: "client-model",
			wantModel:   testutils.Ptr("client-model"),
			description: "client model used as fallback",
		},
		{
			name:           "applies provider to resolved model",
			clientModel:    "mistral-7b",
			clientProvider: "huggingface",
			wantModel:      testutils.Ptr("mistral-7b:huggingface"),
			description:    "provider applied to client model",
		},
		{
			name:           "applies options provider to request model",
			reqModel:       testutils.Ptr("mistral-7b"),
			clientModel:    "client-model",
			clientProvider: "client-provider",
			optsProvider:   "opts-provider",
			wantModel:      testutils.Ptr("mistral-7b:opts-provider"),
			description:    "options provider applied to request model",
		},
		{
			name:         "request model with provider ignores provider option",
			reqModel:     testutils.Ptr("mistral-7b:mistral"),
			clientModel:  "client-model",
			optsProvider: "huggingface",
			wantModel:    testutils.Ptr("mistral-7b:mistral"),
			description:  "existing provider in model not overridden",
		},
		{
			name:           "empty request model falls back to options model",
			reqModel:       testutils.Ptr(""),
			optsModel:      "opts-model",
			clientProvider: "client-provider",
			optsProvider:   "opts-provider",
			wantModel:      testutils.Ptr("opts-model:opts-provider"),
			description:    "empty string treated as nil for fallback",
		},
		{
			name:        "no provider applied when provider is empty",
			clientModel: "mistral-7b",
			wantModel:   testutils.Ptr("mistral-7b"),
			description: "model unchanged when no provider",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := &ChatRequest{Model: tc.reqModel}
			// Manually set client-level options by creating a base options and applying overrides
			baseOpts := request.NewOptions().
				WithModel(tc.clientModel).
				WithProvider(tc.clientProvider)

			// For this test, we'll apply the options with the client defaults
			var optsOverride request.Options
			if tc.optsModel == "" && tc.optsProvider == "" {
				// Use only client options
				optsOverride = baseOpts
			} else {
				// Apply client defaults first, then override with opts
				optsOverride = baseOpts.With(
					WithModel(tc.optsModel),
					WithProvider(tc.optsProvider),
				)
			}

			resolveModel(payload, optsOverride)

			if tc.wantModel == nil {
				if payload.Model != nil {
					t.Fatalf("expected nil, got %#v", payload.Model)
				}
			} else {
				if payload.Model == nil {
					t.Fatalf("expected %#v, got nil", tc.wantModel)
				}
				if *payload.Model != *tc.wantModel {
					t.Fatalf("expected %#v, got %#v", *tc.wantModel, *payload.Model)
				}
			}
		})
	}
}

func TestChatService_ProviderFallback(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name           string
		clientModel    string
		clientProvider string
		optsModel      *string
		optsProvider   *string
		reqModel       *string
		wantModel      string
		description    string
	}

	cases := []TestCase{
		{
			name:           "applies client provider to client model",
			clientModel:    "mistral-7b",
			clientProvider: "huggingface",
			wantModel:      "mistral-7b:huggingface",
			description:    "client model + client provider",
		},
		{
			name:         "applies request provider to request model",
			clientModel:  "default-model",
			reqModel:     testutils.Ptr("mistral-7b"),
			optsProvider:   testutils.Ptr("huggingface"),
			wantModel:    "mistral-7b:huggingface",
			description:  "request model + request provider",
		},
		{
			name:           "request provider overrides client provider",
			clientModel:    "mistral-7b",
			clientProvider: "mistral",
			optsProvider:   testutils.Ptr("huggingface"),
			wantModel:      "mistral-7b:huggingface",
			description:    "request provider overrides client provider",
		},
		{
			name:           "ignores provider when model already has provider",
			clientModel:    "mistral-7b:mistral",
			clientProvider: "huggingface",
			wantModel:      "mistral-7b:mistral",
			description:    "model with provider + provider option (ignored)",
		},
		{
			name:           "request model takes precedence over provider",
			clientModel:    "default-model",
			clientProvider: "huggingface",
			reqModel:       testutils.Ptr("mistral-7b:mistral"),
			optsProvider:   testutils.Ptr("huggingface"),
			wantModel:      "mistral-7b:mistral",
			description:    "request model with provider + request provider (model wins)",
		},
		{
			name:           "applies request provider to request model without provider",
			clientModel:    "default-model",
			clientProvider: "client-provider",
			reqModel:       testutils.Ptr("mistral-7b"),
			optsProvider:   testutils.Ptr("huggingface"),
			wantModel:      "mistral-7b:huggingface",
			description:    "request model + request provider (overrides client)",
		},
		{
			name:        "no provider applied when both missing",
			clientModel: "mistral-7b",
			wantModel:   "mistral-7b",
			description: "model without provider, no provider option",
		},
		{
			name:           "client provider applied to request model when no request provider",
			clientModel:    "default-model",
			clientProvider: "client-provider",
			reqModel:       testutils.Ptr("mistral-7b"),
			wantModel:      "mistral-7b:client-provider",
			description:    "client provider with request model",
		},
		{
			name:           "model with multiple colons not modified by provider",
			clientModel:    "org:model:variant",
			clientProvider: "huggingface",
			wantModel:      "org:model:variant",
			description:    "multiple colons in model",
		},
		{
			name:           "optsModel without optsProvider uses clientProvider",
			clientModel:    "default-model",
			clientProvider: "client-provider",
			optsModel:      testutils.Ptr("opts-model"),
			wantModel:      "opts-model:client-provider",
			description:    "opts model with client provider fallback",
		},
		{
			name:           "both optsModel and optsProvider set override client defaults",
			clientModel:    "client-model",
			clientProvider: "client-provider",
			optsModel:      testutils.Ptr("opts-model"),
			optsProvider:   testutils.Ptr("opts-provider"),
			wantModel:      "opts-model:opts-provider",
			description:    "both request-level options override client defaults",
		},
		{
			name:           "empty reqModel falls back to optsModel",
			clientModel:    "default-model",
			clientProvider: "client-provider",
			reqModel:       testutils.Ptr(""),
			optsModel:      testutils.Ptr("opts-model"),
			optsProvider:   testutils.Ptr("opts-provider"),
			wantModel:      "opts-model:opts-provider",
			description:    "empty request model treated as nil for fallback",
		},
		{
			name:           "model with trailing colon not modified by provider",
			clientModel:    "model:",
			clientProvider: "provider",
			wantModel:      "model:",
			description:    "model with trailing colon treated as having provider",
		},
		{
			name:           "override client provider with empty provider",
			clientModel:    "mistral-7b",
			clientProvider: "huggingface",
			optsProvider:   testutils.Ptr(""),
			wantModel:      "mistral-7b",
			description:    "explicitly passing empty provider removes provider",
		},
		{
			name:           "override optsModel with empty, keep optsProvider",
			clientModel:    "default-model",
			clientProvider: "client-provider",
			optsProvider:   testutils.Ptr("opts-provider"),
			wantModel:      "default-model:opts-provider",
			description:    "request provider override without model override uses client model",
		},
		{
			name:           "override client provider with opts provider, no model override",
			clientModel:    "mistral-7b",
			clientProvider: "huggingface",
			optsProvider:   testutils.Ptr("mistral"),
			wantModel:      "mistral-7b:mistral",
			description:    "request provider overrides client provider without model change",
		},
	}

	// testImpl is a helper that tests model and provider resolution
	// for both Complete and CompleteStream methods. It handles all setup, method invocation,
	// and verification of the resolved model in the request.
	testImpl := func(
		t *testing.T,
		tc TestCase,
		mtFactory func() *testutils.MockTransport,
		methodCall func(*ChatService, *ChatRequest, []Option) error,
	) {
		t.Helper()

		mt := mtFactory()
		opts := request.NewOptions().
			WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
			WithModel(tc.clientModel).
			WithProvider(tc.clientProvider)
		svc := newChatService(opts)

		text := "hi"
		req := &ChatRequest{
			Model: tc.reqModel,
			Messages: []ChatMessage{
				{Role: "user", Content: ChatMessageContent{Text: &text}},
			},
		}

		optsToPass := []Option{}
		if tc.optsModel != nil {
			optsToPass = append(optsToPass, WithModel(*tc.optsModel))
		}
		if tc.optsProvider != nil {
			optsToPass = append(optsToPass, WithProvider(*tc.optsProvider))
		}

		err := methodCall(&svc, req, optsToPass)
		require.NoError(t, err)

		require.NotNil(t, mt.LastRequest)
		require.Equal(t, EndpointChatCompletion, mt.LastRequest.URL.Path)

		body, err := io.ReadAll(mt.LastRequest.Body)
		require.NoError(t, err)
		_ = mt.LastRequest.Body.Close()

		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		require.Equal(t, tc.wantModel, got["model"], tc.description)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Complete method
			testImpl(t, tc,
				func() *testutils.MockTransport {
					return testutils.NewJSONMockTransport(
						http.StatusOK,
						chatServiceResponseBody,
						nil,
					)
				},
				func(svc *ChatService, req *ChatRequest, opts []Option) error {
					_, err := svc.Complete(req, opts...)

					return err
				},
			)

			// Test CompleteStream method
			sseBody := "data: {\"id\":\"id\",\"created\":1,\"model\":\"" + tc.wantModel + "\",\"system_fingerprint\":\"sig\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"hi\"}}]}\n\n" +
				"data: [DONE]\n\n"
			testImpl(t, tc,
				func() *testutils.MockTransport {
					mt := testutils.NewMockTransport(http.StatusOK, sseBody, nil)
					mt.Response.Header.Set("Content-Type", "text/event-stream")

					return mt
				},
				func(svc *ChatService, req *ChatRequest, opts []Option) error {
					stream, err := svc.CompleteStream(req, opts...)
					if err != nil {
						return err
					}
					defer func() { _ = stream.Close() }()

					chunk, err := stream.Recv(context.Background())
					if err != nil {
						return err
					}
					require.Equal(t, tc.wantModel, chunk.Model)

					return nil
				},
			)
		})
	}
}
