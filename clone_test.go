//go:build !integration

package hfgo

import (
	"net/http"
	"sync"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
