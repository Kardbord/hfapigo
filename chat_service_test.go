package hfapigo

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/request"
	"github.com/Kardbord/hfapigo/v4/internal/testutils"
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
			opts := request.NewRequestOptions().
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

			if mt.LastRequest == nil {
				t.Fatal("expected request to be sent")
			}
			if mt.LastRequest.URL.Path != EndpointChatCompletion {
				t.Fatalf("unexpected path: %s", mt.LastRequest.URL.Path)
			}

			body, err := io.ReadAll(mt.LastRequest.Body)
			require.NoError(t, err)
			_ = mt.LastRequest.Body.Close()

			var got map[string]any
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("unexpected json: %v", err)
			}
			if got["model"] != tc.wantModel {
				t.Fatalf("unexpected model: %#v", got["model"])
			}

			if tc.reqModel == nil {
				if req.Model != nil {
					t.Fatalf("expected request model to remain nil, got %#v", req.Model)
				}
			} else if req.Model == nil || *req.Model != *tc.reqModel {
				t.Fatalf("unexpected request model: %#v", req.Model)
			}
		})
	}
}

func TestChatService_Complete_ModelValidation(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
	opts := request.NewRequestOptions().
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
	testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
	if mt.LastRequest != nil {
		t.Fatalf("expected no request, got %#v", mt.LastRequest)
	}
}

func TestChatService_Complete_NilRequest(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, chatServiceResponseBody, nil)
	opts := request.NewRequestOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newChatService(opts)

	_, err := svc.Complete(nil)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, internalErrors.SDKErrorKindConfiguration)
	if mt.LastRequest != nil {
		t.Fatalf("expected no request, got %#v", mt.LastRequest)
	}
}
