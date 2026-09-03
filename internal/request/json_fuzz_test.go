//go:build !integration

package request

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Kardbord/hfgo/v4/internal/testutils"
)

type fuzzDoJSONReq struct {
	Inputs string `json:"inputs"`
}

type fuzzDoJSONResp struct {
	GeneratedText string `json:"generated_text"`
}

func FuzzDoJSON(f *testing.F) {
	f.Add([]byte(`{"generated_text":"hello"}`))
	f.Add([]byte(``))
	f.Add([]byte(`{not valid json}`))
	f.Add([]byte(
		`{"generated_text":"` + strings.Repeat("a", int(DefaultMaxResponseBodyBytes)+1) + `"}`))
	f.Fuzz(func(_ *testing.T, respBody []byte) {
		mt := &testutils.MockTransport{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(respBody)),
				Header:     make(http.Header),
			},
		}
		opts := NewOptions().WithHTTPClientFactory(func() http.Client {
			return testutils.NewMockHTTPClient(mt)
		})
		_, _ = DoJSON[fuzzDoJSONReq, fuzzDoJSONResp](
			opts, http.MethodPost, "/test", fuzzDoJSONReq{Inputs: "test"},
		)
	})
}

type fuzzDoJSONStreamReq struct {
	Prompt string `json:"prompt"`
}

type fuzzDoJSONStreamResp struct {
	Text string `json:"text"`
}

func FuzzDoJSONStream(f *testing.F) {
	f.Add([]byte("data: {\"text\":\"hello\"}\n\ndata: [DONE]\n\n"))
	f.Add([]byte(``))
	f.Add([]byte("data: [DONE]\n\n"))
	f.Add([]byte("data: {\"text\":\"hello\"}\n\n"))
	f.Add([]byte("not sse"))
	f.Fuzz(func(_ *testing.T, body []byte) {
		mt := &testutils.MockTransport{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     make(http.Header),
			},
		}
		mt.Response.Header.Set("Content-Type", "text/event-stream")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		opts := NewOptions().
			WithContext(ctx).
			WithHTTPClientFactory(func() http.Client {
				return testutils.NewMockHTTPClient(mt)
			})

		stream, err := DoJSONStream[fuzzDoJSONStreamReq, fuzzDoJSONStreamResp](
			opts, http.MethodPost, "/stream", fuzzDoJSONStreamReq{Prompt: "test"},
		)
		if err != nil {
			return
		}
		defer func() { _ = stream.Close() }()

		for {
			_, err = stream.Recv(ctx)
			if err != nil {
				return
			}
		}
	})
}
