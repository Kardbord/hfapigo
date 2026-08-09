//go:build !integration

package hfgo

import (
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestSummarizationService_Summarize_ResponseVariations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		responseBody string
		wantLen      int
		wantSummary  string
		description  string
	}{
		{
			name:         "single summary",
			responseBody: `[{"summary_text":"A concise summary."}]`,
			wantLen:      1,
			wantSummary:  "A concise summary.",
			description:  "a single summary is returned",
		},
		{
			name:         "empty response",
			responseBody: `[]`,
			wantLen:      0,
			description:  "empty response passes through as an empty list",
		},
		{
			name:         "multiple summaries",
			responseBody: `[{"summary_text":"One."},{"summary_text":"Two."}]`,
			wantLen:      2,
			wantSummary:  "One.",
			description:  "multiple summaries in one response pass through",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(http.StatusOK, tc.responseBody, nil)
			client := NewClient(
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
				WithModel("test-model"),
			)

			result, err := client.Summarize(SummarizationRequest{Input: "Some long text."})
			require.NoError(t, err, tc.description)
			require.NotNil(t, result, tc.description)
			require.Len(t, result, tc.wantLen, tc.description)

			if tc.wantLen > 0 {
				require.Equal(t, tc.wantSummary, result[0].SummaryText, tc.description)
			}
		})
	}
}

func TestSummarizationService_Summarize_ParameterSerialization(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		cleanUp     *bool
		truncation  *string
		generate    map[string]any
		want        map[string]any
		description string
	}{
		{
			name:        "parameters omitted",
			description: "nil parameters are omitted from the request body",
		},
		{
			name:        "clean up tokenization spaces",
			cleanUp:     testutils.Ptr(true),
			want:        map[string]any{"clean_up_tokenization_spaces": true},
			description: "clean_up_tokenization_spaces maps to its JSON key",
		},
		{
			name:        "truncation do_not_truncate",
			truncation:  testutils.Ptr(SummarizationTruncationDoNotTruncate),
			want:        map[string]any{"truncation": SummarizationTruncationDoNotTruncate},
			description: "do_not_truncate constant serializes correctly",
		},
		{
			name:        "truncation longest_first",
			truncation:  testutils.Ptr(SummarizationTruncationLongestFirst),
			want:        map[string]any{"truncation": SummarizationTruncationLongestFirst},
			description: "longest_first constant serializes correctly",
		},
		{
			name:        "truncation only_first",
			truncation:  testutils.Ptr(SummarizationTruncationOnlyFirst),
			want:        map[string]any{"truncation": SummarizationTruncationOnlyFirst},
			description: "only_first constant serializes correctly",
		},
		{
			name:        "truncation only_second",
			truncation:  testutils.Ptr(SummarizationTruncationOnlySecond),
			want:        map[string]any{"truncation": SummarizationTruncationOnlySecond},
			description: "only_second constant serializes correctly",
		},
		{
			name:     "generate parameters",
			generate: map[string]any{"max_new_tokens": 60, "temperature": 0.8},
			want: map[string]any{
				"generate_parameters": map[string]any{
					"max_new_tokens": float64(60),
					"temperature":    0.8,
				},
			},
			description: "generate_parameters maps to its JSON key",
		},
		{
			name:       "all parameters",
			cleanUp:    testutils.Ptr(true),
			truncation: testutils.Ptr(SummarizationTruncationOnlyFirst),
			generate:   map[string]any{"max_new_tokens": 60},
			want: map[string]any{
				"clean_up_tokenization_spaces": true,
				"truncation":                   SummarizationTruncationOnlyFirst,
				"generate_parameters":          map[string]any{"max_new_tokens": float64(60)},
			},
			description: "all parameters serialize together",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(
				http.StatusOK,
				`[{"summary_text":"A concise summary."}]`,
				nil,
			)
			client := NewClient(
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
				WithModel("test-model"),
			)

			req := SummarizationRequest{Input: "Some long text."}
			if tc.cleanUp != nil || tc.truncation != nil || tc.generate != nil {
				req.Parameters = &SummarizationParameters{
					CleanUpTokenizationSpaces: tc.cleanUp,
					Truncation:                tc.truncation,
					GenerateParameters:        tc.generate,
				}
			}

			result, err := client.Summarize(req)
			require.NoError(t, err, tc.description)
			require.NotNil(t, result, tc.description)

			reqBody := testutils.ReadRequestBody(t, mt)
			if tc.want == nil {
				_, ok := reqBody["parameters"]
				require.False(t, ok, tc.description)

				return
			}

			require.Equal(t, tc.want, reqBody["parameters"], tc.description)
		})
	}
}

func TestSummarizationService_Summarize_Errors(t *testing.T) {
	t.Parallel()

	runErrorCases(t,
		[]errorCase{
			{
				name:         "no model configured",
				statusCode:   http.StatusOK,
				responseBody: `[{"summary_text":"A concise summary."}]`,
				want:         testutils.WantErrSDK,
				sdkErrKind:   hferrors.SDKErrorKindConfiguration,
				description:  "SDK error when model is missing",
			},
			{
				name:         "API error on 404",
				withModel:    true,
				statusCode:   http.StatusNotFound,
				responseBody: `{"error":"Model not found"}`,
				want:         testutils.WantErrAPI,
				description:  "API error for nonexistent model",
			},
		},
		func(opts ...Option) ([]Summarization, error) {
			return NewClient(opts...).Summarize(SummarizationRequest{
				Input: "Some long text that should be summarized.",
			})
		},
	)
}

func TestSummarizationService_SummarizeBatch_ResponseVariations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		responseBody string
		want         []string
		description  string
	}{
		{
			name:         "single input",
			responseBody: `[{"summary_text":"Summary one."}]`,
			want:         []string{"Summary one."},
			description:  "a single batched input returns one summary",
		},
		{
			name:         "multiple inputs",
			responseBody: `[{"summary_text":"Summary one."},{"summary_text":"Summary two."}]`,
			want:         []string{"Summary one.", "Summary two."},
			description:  "each batched input returns its own flat summary",
		},
		{
			name:         "empty response",
			responseBody: `[]`,
			want:         []string{},
			description:  "empty response passes through as an empty list",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(http.StatusOK, tc.responseBody, nil)
			client := NewClient(
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
				WithModel("test-model"),
			)

			result, err := client.SummarizeBatch(SummarizationBatchRequest{
				Inputs: []string{"Long text one.", "Long text two."},
			})
			require.NoError(t, err, tc.description)
			require.NotNil(t, result, tc.description)
			require.Len(t, result, len(tc.want), tc.description)

			for j, summary := range result {
				require.Equal(t, tc.want[j], summary.SummaryText, tc.description)
			}
		})
	}
}

func TestSummarizationService_SummarizeBatch_Errors(t *testing.T) {
	t.Parallel()

	runErrorCases(t,
		[]errorCase{
			{
				name:         "no model configured",
				statusCode:   http.StatusOK,
				responseBody: `[{"summary_text":"Summary one."}]`,
				want:         testutils.WantErrSDK,
				sdkErrKind:   hferrors.SDKErrorKindConfiguration,
				description:  "SDK error when model is missing",
			},
			{
				name:         "API error on 404",
				withModel:    true,
				statusCode:   http.StatusNotFound,
				responseBody: `{"error":"Model not found"}`,
				want:         testutils.WantErrAPI,
				description:  "API error for nonexistent model",
			},
		},
		func(opts ...Option) ([]Summarization, error) {
			return NewClient(opts...).SummarizeBatch(SummarizationBatchRequest{
				Inputs: []string{"Long text one."},
			})
		},
	)
}

func TestSummarizationService_SummarizeBatch_ModelFromOptions(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		`[{"summary_text":"Summary one."}]`,
		nil,
	)
	client := NewClient(
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
	)

	result, err := client.SummarizeBatch(SummarizationBatchRequest{
		Inputs: []string{"Long text one."},
	}, WithModel("override-model"))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the correct model was used in the request
	require.NotNil(t, mt.LastRequest)
	require.Contains(t, mt.LastRequest.URL.Path, "override-model")
}
