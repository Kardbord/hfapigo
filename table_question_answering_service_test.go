//go:build !integration

package hfgo

import (
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestTableQuestionAnsweringService_AnswerTableQuestion_ResponseDecoding(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		responseBody string
		expectedLen  int
		wantAnswer   string
		wantCells    []string
		wantCoords   [][]int
		wantAggr     *string
		description  string
	}{
		{
			name:         "single answer",
			responseBody: `{"answer":"30","cells":["30"],"coordinates":[[1,1]]}`,
			expectedLen:  1,
			wantAnswer:   "30",
			wantCells:    []string{"30"},
			wantCoords:   [][]int{{1, 1}},
			description:  "basic single answer response",
		},
		{
			name:         "answer with aggregator",
			responseBody: `{"answer":"COUNT > 3","cells":["Alice","Bob","Carol"],"coordinates":[[0,0],[0,1],[0,2]],"aggregator":"COUNT"}`,
			expectedLen:  1,
			wantAnswer:   "COUNT > 3",
			wantCells:    []string{"Alice", "Bob", "Carol"},
			wantCoords:   [][]int{{0, 0}, {0, 1}, {0, 2}},
			wantAggr:     testutils.Ptr("COUNT"),
			description:  "answer with aggregator field",
		},
		{
			name:         "answer with multiple cells",
			responseBody: `{"answer":"Alice and Bob","cells":["Alice","Bob"],"coordinates":[[0,0],[0,1]]}`,
			expectedLen:  1,
			wantAnswer:   "Alice and Bob",
			wantCells:    []string{"Alice", "Bob"},
			wantCoords:   [][]int{{0, 0}, {0, 1}},
			description:  "answer spanning multiple cells",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(http.StatusOK, tc.responseBody, nil)
			client := NewClient(
				WithHTTPClientFactory(
					func() http.Client { return testutils.NewMockHTTPClient(mt) },
				),
				WithModel("test-model"),
			)

			req := TableQuestionAnsweringRequest{
				Input: TableQuestionAnsweringInputData{
					Question: "How old is Bob?",
					Table: map[string][]string{
						"Name": {"Alice", "Bob", "Carol"},
						"Age":  {"25", "30", "35"},
					},
				},
			}

			result, err := client.AnswerTableQuestion(req)
			require.NoError(t, err, tc.description)
			require.NotNil(t, result)
			require.Len(t, result, tc.expectedLen, tc.description)
			require.Equal(t, tc.wantAnswer, result[0].Answer)
			require.Equal(t, tc.wantCells, result[0].Cells)
			require.Equal(t, tc.wantCoords, result[0].Coordinates)
			if tc.wantAggr != nil {
				require.NotNil(t, result[0].Aggregator)
				require.Equal(t, *tc.wantAggr, *result[0].Aggregator)
			} else {
				require.Nil(t, result[0].Aggregator)
			}
		})
	}
}

func TestTableQuestionAnsweringService_AnswerTableQuestion_WithParameters(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		`{"answer":"30","cells":["30"],"coordinates":[[1,1]]}`,
		nil,
	)
	client := NewClient(
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
		WithModel("test-model"),
	)

	padding := TableQuestionAnsweringPaddingMaxLength
	sequential := false
	truncation := true
	req := TableQuestionAnsweringRequest{
		Input: TableQuestionAnsweringInputData{
			Question: "How old is Bob?",
			Table: map[string][]string{
				"Name": {"Alice", "Bob", "Carol"},
				"Age":  {"25", "30", "35"},
			},
		},
		Parameters: &TableQuestionAnsweringParameters{
			Padding:    &padding,
			Sequential: &sequential,
			Truncation: &truncation,
		},
	}

	result, err := client.AnswerTableQuestion(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	reqBody := testutils.ReadRequestBody(t, mt)
	params, ok := reqBody["parameters"].(map[string]any)
	require.True(t, ok, "parameters should be a map")
	require.Equal(t, "max_length", params["padding"])
	require.Equal(t, false, params["sequential"])
	require.Equal(t, true, params["truncation"])
}

func TestTableQuestionAnsweringService_AnswerTableQuestion_Errors(t *testing.T) {
	t.Parallel()

	runErrorCases(t,
		[]errorCase{
			{
				name:         "no model configured",
				statusCode:   http.StatusOK,
				responseBody: `{"answer":"30","cells":["30"],"coordinates":[[1,1]]}`,
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
			{
				name:         "API error on 503",
				withModel:    true,
				statusCode:   http.StatusServiceUnavailable,
				responseBody: `{"error":"Model loading"}`,
				want:         testutils.WantErrAPI,
				description:  "API error for model not yet loaded",
			},
		},
		func(opts ...Option) ([]TableQuestionAnswering, error) {
			return NewClient(opts...).AnswerTableQuestion(TableQuestionAnsweringRequest{
				Input: TableQuestionAnsweringInputData{
					Question: "How old is Bob?",
					Table: map[string][]string{
						"Name": {"Alice", "Bob", "Carol"},
						"Age":  {"25", "30", "35"},
					},
				},
			})
		},
	)
}

func TestTableQuestionAnsweringService_AnswerTableQuestion_ModelFromOptions(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		`{"answer":"30","cells":["30"],"coordinates":[[1,1]]}`,
		nil,
	)
	client := NewClient(
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }),
	)

	req := TableQuestionAnsweringRequest{
		Input: TableQuestionAnsweringInputData{
			Question: "How old is Bob?",
			Table: map[string][]string{
				"Name": {"Alice", "Bob", "Carol"},
				"Age":  {"25", "30", "35"},
			},
		},
	}

	result, err := client.AnswerTableQuestion(req, WithModel("override-model"))
	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotNil(t, mt.LastRequest)
	require.Contains(t, mt.LastRequest.URL.Path, "override-model")

	reqBody := testutils.ReadRequestBody(t, mt)
	inputs, ok := reqBody["inputs"].(map[string]any)
	require.True(t, ok, "inputs should be an object")
	require.Equal(t, "How old is Bob?", inputs["question"])

	table, ok := inputs["table"].(map[string]any)
	require.True(t, ok, "table should be an object")
	nameCol, ok := table["Name"].([]any)
	require.True(t, ok, "Name column should be an array")
	require.Len(t, nameCol, 3)
	require.Equal(t, "Alice", nameCol[0])
}
