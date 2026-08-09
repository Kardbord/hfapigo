//go:build !integration

package hfgo

import (
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestTranslationService_Translate_ResponseVariations(t *testing.T) {
	t.Parallel()

	runSingleResponseVariations(t,
		[]singleResponseVariationCase{
			{
				name:         "single translation",
				responseBody: `[{"translation_text":"Bonjour, comment allez-vous aujourd'hui ?"}]`,
				wantLen:      1,
				wantText:     "Bonjour, comment allez-vous aujourd'hui ?",
				description:  "a single translation is returned",
			},
			{
				name:         "empty response",
				responseBody: `[]`,
				wantLen:      0,
				description:  "empty response passes through as an empty list",
			},
			{
				name:         "multiple translations",
				responseBody: `[{"translation_text":"Un."},{"translation_text":"Deux."}]`,
				wantLen:      2,
				wantText:     "Un.",
				description:  "multiple translations in one response pass through",
			},
		},
		func() TranslationRequest { return TranslationRequest{Input: "Hello, how are you?"} },
		func(c Client, req TranslationRequest) ([]Translation, error) { return c.Translate(req) },
		func(tr Translation) string { return tr.TranslationText },
	)
}

func TestTranslationService_Translate_ParameterSerialization(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		cleanUp     *bool
		srcLang     *string
		tgtLang     *string
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
			name:        "source language",
			srcLang:     testutils.Ptr("en"),
			want:        map[string]any{"src_lang": "en"},
			description: "src_lang maps to its JSON key",
		},
		{
			name:        "target language",
			tgtLang:     testutils.Ptr("fr"),
			want:        map[string]any{"tgt_lang": "fr"},
			description: "tgt_lang maps to its JSON key",
		},
		{
			name:        "truncation do_not_truncate",
			truncation:  testutils.Ptr(TranslationTruncationDoNotTruncate),
			want:        map[string]any{"truncation": TranslationTruncationDoNotTruncate},
			description: "do_not_truncate constant serializes correctly",
		},
		{
			name:        "truncation longest_first",
			truncation:  testutils.Ptr(TranslationTruncationLongestFirst),
			want:        map[string]any{"truncation": TranslationTruncationLongestFirst},
			description: "longest_first constant serializes correctly",
		},
		{
			name:        "truncation only_first",
			truncation:  testutils.Ptr(TranslationTruncationOnlyFirst),
			want:        map[string]any{"truncation": TranslationTruncationOnlyFirst},
			description: "only_first constant serializes correctly",
		},
		{
			name:        "truncation only_second",
			truncation:  testutils.Ptr(TranslationTruncationOnlySecond),
			want:        map[string]any{"truncation": TranslationTruncationOnlySecond},
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
			srcLang:    testutils.Ptr("en"),
			tgtLang:    testutils.Ptr("de"),
			truncation: testutils.Ptr(TranslationTruncationOnlyFirst),
			generate:   map[string]any{"max_new_tokens": 60},
			want: map[string]any{
				"clean_up_tokenization_spaces": true,
				"src_lang":                     "en",
				"tgt_lang":                     "de",
				"truncation":                   TranslationTruncationOnlyFirst,
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
				`[{"translation_text":"Test translation."}]`,
				nil,
			)
			client := NewClient(
				WithHTTPClientFactory(
					func() http.Client { return testutils.NewMockHTTPClient(mt) },
				),
				WithModel("test-model"),
			)

			req := TranslationRequest{Input: "Hello, how are you?"}
			if tc.cleanUp != nil || tc.srcLang != nil || tc.tgtLang != nil ||
				tc.truncation != nil || tc.generate != nil {
				req.Parameters = &TranslationParameters{
					CleanUpTokenizationSpaces: tc.cleanUp,
					SrcLang:                   tc.srcLang,
					TgtLang:                   tc.tgtLang,
					Truncation:                tc.truncation,
					GenerateParameters:        tc.generate,
				}
			}

			result, err := client.Translate(req)
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

func TestTranslationService_Translate_Errors(t *testing.T) {
	t.Parallel()

	runErrorCases(t,
		[]errorCase{
			{
				name:         "no model configured",
				statusCode:   http.StatusOK,
				responseBody: `[{"translation_text":"Test translation."}]`,
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
		func(opts ...Option) ([]Translation, error) {
			return NewClient(opts...).Translate(TranslationRequest{
				Input: "Hello, how are you?",
			})
		},
	)
}

func TestTranslationService_TranslateBatch_ResponseVariations(t *testing.T) {
	t.Parallel()

	runBatchResponseVariations(
		t,
		[]batchResponseVariationCase{
			{
				name:         "single input",
				responseBody: `[{"translation_text":"Bonjour."}]`,
				want:         []string{"Bonjour."},
				description:  "a single batched input returns one translation",
			},
			{
				name:         "multiple inputs",
				responseBody: `[{"translation_text":"Bonjour."},{"translation_text":"À demain."}]`,
				want:         []string{"Bonjour.", "À demain."},
				description:  "each batched input returns its own flat translation",
			},
			{
				name:         "empty response",
				responseBody: `[]`,
				want:         []string{},
				description:  "empty response passes through as an empty list",
			},
		},
		func() TranslationBatchRequest {
			return TranslationBatchRequest{Inputs: []string{"Hello.", "See you tomorrow."}}
		},
		func(c Client, req TranslationBatchRequest) ([]Translation, error) { return c.TranslateBatch(req) },
		func(tr Translation) string { return tr.TranslationText },
	)
}

func TestTranslationService_TranslateBatch_Errors(t *testing.T) {
	t.Parallel()

	runErrorCases(t,
		[]errorCase{
			{
				name:         "no model configured",
				statusCode:   http.StatusOK,
				responseBody: `[{"translation_text":"Bonjour."}]`,
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
		func(opts ...Option) ([]Translation, error) {
			return NewClient(opts...).TranslateBatch(TranslationBatchRequest{
				Inputs: []string{"Hello."},
			})
		},
	)
}

func TestTranslationService_TranslateBatch_ModelFromOptions(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		`[{"translation_text":"Bonjour."}]`,
		nil,
	)
	client := NewClient(
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
	)

	result, err := client.TranslateBatch(TranslationBatchRequest{
		Inputs: []string{"Hello."},
	}, WithModel("override-model"))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the correct model was used in the request
	require.NotNil(t, mt.LastRequest)
	require.Contains(t, mt.LastRequest.URL.Path, "override-model")
}
