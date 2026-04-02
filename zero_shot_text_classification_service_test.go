//go:build !integration

package hfgo

import (
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/request"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestZeroShotTextClassificationService_Classify_SingleInput(t *testing.T) {
	t.Parallel()

	const zeroShotSingleClassificationResponseBody = `[{"label":"positive","score":0.95}]`
	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		zeroShotSingleClassificationResponseBody,
		nil,
	)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newZeroShotTextClassificationService(opts)

	candidateLabels := []string{"positive", "negative", "neutral"}
	req := ZeroShotTextClassificationRequest{
		Input: "This is a great product!",
		Parameters: &ZeroShotTextClassificationParameters{
			CandidateLabels: candidateLabels,
		},
	}

	result, err := svc.Classify(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 1)
	require.Equal(t, "positive", result[0].Label)
	require.InEpsilon(t, 0.95, result[0].Score, 0.001)
}

//nolint:dupl // Similar error test setup for different methods (Classify vs ClassifyBatch)
func TestZeroShotTextClassificationService_Classify_Errors(t *testing.T) {
	t.Parallel()

	const zeroShotSingleClassificationResponseBody = `[{"label":"positive","score":0.95}]`

	cases := []struct {
		name            string
		req             ZeroShotTextClassificationRequest
		withModel       bool
		httpStatusCode  int
		responseBody    string
		expectedErrKind *hferrors.SDKErrorKind
		description     string
	}{
		{
			name: "no candidate labels",
			req: ZeroShotTextClassificationRequest{
				Input: "test text",
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{},
				},
			},
			withModel:       true,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotSingleClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when candidate labels are empty",
		},
		{
			name: "no parameters",
			req: ZeroShotTextClassificationRequest{
				Input: "test text",
			},
			withModel:       true,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotSingleClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when parameters are missing",
		},
		{
			name: "no model configured",
			req: ZeroShotTextClassificationRequest{
				Input: "test text",
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{"positive", "negative"},
				},
			},
			withModel:       false,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotSingleClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when model is missing",
		},
		{
			name: "API error on 404",
			req: ZeroShotTextClassificationRequest{
				Input: "test text",
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{"positive", "negative"},
				},
			},
			withModel:       true,
			httpStatusCode:  http.StatusNotFound,
			responseBody:    `{"error":"Model not found"}`,
			expectedErrKind: nil, // API error, not SDK error
			description:     "API error for nonexistent model",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(tc.httpStatusCode, tc.responseBody, nil)
			opts := request.NewOptions().
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			if tc.withModel {
				opts = opts.WithModel("nonexistent-model")
			}
			svc := newZeroShotTextClassificationService(opts)

			result, err := svc.Classify(tc.req)
			require.Error(t, err, tc.description)
			require.Nil(t, result)

			if tc.expectedErrKind != nil {
				// SDK error expected
				testutils.AssertSDKErrorKind(t, err, *tc.expectedErrKind)
				// Verify no request was made for SDK errors
				require.Nil(t, mt.LastRequest)
			} else {
				// API error expected
				var apiErr *hferrors.APIError
				require.ErrorAs(t, err, &apiErr, tc.description)
				require.Equal(t, tc.httpStatusCode, apiErr.StatusCode)
			}
		})
	}
}

func TestZeroShotTextClassificationService_ClassifyBatch_InputVariations(t *testing.T) {
	t.Parallel()

	const (
		zeroShotBatchClassificationResponseBody = `[{"Sequence":"text1","Labels":["positive","negative","neutral"],"Scores":[0.95,0.03,0.02]}]`
		zeroShotBatchMultipleResponseBody       = `[{"Sequence":"text1","Labels":["positive","negative","neutral"],"Scores":[0.95,0.03,0.02]},{"Sequence":"text2","Labels":["negative","positive","neutral"],"Scores":[0.87,0.10,0.03]},{"Sequence":"text3","Labels":["neutral","positive","negative"],"Scores":[0.75,0.15,0.10]}]`
	)

	cases := []struct {
		name                string
		responseBody        string
		inputs              []string
		expectedOuterLength int
		expectedInnerLength int
		expectedFirstLabel  string
		expectedFirstScore  float64
		expectedSecondLabel string
		expectedSecondScore float64
		description         string
	}{
		{
			name:                "single input",
			responseBody:        zeroShotBatchClassificationResponseBody,
			inputs:              []string{"text1"},
			expectedOuterLength: 1,
			expectedInnerLength: 3,
			expectedFirstLabel:  "positive",
			expectedFirstScore:  0.95,
			description:         "single input classification",
		},
		{
			name:                "multiple inputs",
			responseBody:        zeroShotBatchMultipleResponseBody,
			inputs:              []string{"text1", "text2", "text3"},
			expectedOuterLength: 3,
			expectedInnerLength: 3,
			expectedFirstLabel:  "positive",
			expectedFirstScore:  0.95,
			expectedSecondLabel: "negative",
			expectedSecondScore: 0.87,
			description:         "multiple input classifications",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(http.StatusOK, tc.responseBody, nil)
			opts := request.NewOptions().
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
				WithModel("test-model")
			svc := newZeroShotTextClassificationService(opts)

			candidateLabels := []string{"positive", "negative", "neutral"}
			req := ZeroShotTextClassificationBatchRequest{
				Inputs: tc.inputs,
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: candidateLabels,
				},
			}

			result, err := svc.ClassifyBatch(req)
			require.NoError(t, err, tc.description)
			require.NotNil(t, result)
			require.Len(t, result, tc.expectedOuterLength, tc.description)

			// Check first input classifications
			require.Len(t, result[0], tc.expectedInnerLength)
			require.Equal(t, tc.expectedFirstLabel, result[0][0].Label)
			require.InEpsilon(t, tc.expectedFirstScore, result[0][0].Score, 0.001)

			// Check second input classifications if multiple inputs
			if tc.expectedOuterLength > 1 {
				require.Len(t, result[1], tc.expectedInnerLength)
				require.Equal(t, tc.expectedSecondLabel, result[1][0].Label)
				require.InEpsilon(t, tc.expectedSecondScore, result[1][0].Score, 0.001)
			}
		})
	}
}

//nolint:dupl // Similar error test setup for different methods (Classify vs ClassifyBatch)
func TestZeroShotTextClassificationService_ClassifyBatch_Errors(t *testing.T) {
	t.Parallel()

	const zeroShotBatchClassificationResponseBody = `[{"Sequence":"text1","Labels":["positive","negative","neutral"],"Scores":[0.95,0.03,0.02]}]`

	cases := []struct {
		name            string
		req             ZeroShotTextClassificationBatchRequest
		withModel       bool
		httpStatusCode  int
		responseBody    string
		expectedErrKind *hferrors.SDKErrorKind
		description     string
	}{
		{
			name: "no candidate labels",
			req: ZeroShotTextClassificationBatchRequest{
				Inputs: []string{"test text"},
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{},
				},
			},
			withModel:       true,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotBatchClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when candidate labels are empty",
		},
		{
			name: "no parameters",
			req: ZeroShotTextClassificationBatchRequest{
				Inputs: []string{"test text"},
			},
			withModel:       true,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotBatchClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when parameters are missing",
		},
		{
			name: "no model configured",
			req: ZeroShotTextClassificationBatchRequest{
				Inputs: []string{"test text"},
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{"positive", "negative"},
				},
			},
			withModel:       false,
			httpStatusCode:  http.StatusOK,
			responseBody:    zeroShotBatchClassificationResponseBody,
			expectedErrKind: testutils.Ptr(hferrors.SDKErrorKindConfiguration),
			description:     "SDK error when model is missing",
		},
		{
			name: "API error on 404",
			req: ZeroShotTextClassificationBatchRequest{
				Inputs: []string{"test text"},
				Parameters: &ZeroShotTextClassificationParameters{
					CandidateLabels: []string{"positive", "negative"},
				},
			},
			withModel:       true,
			httpStatusCode:  http.StatusNotFound,
			responseBody:    `{"error":"Model not found"}`,
			expectedErrKind: nil, // API error, not SDK error
			description:     "API error for nonexistent model",
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			mt := testutils.NewJSONMockTransport(tc.httpStatusCode, tc.responseBody, nil)
			opts := request.NewOptions().
				WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
			if tc.withModel {
				opts = opts.WithModel("nonexistent-model")
			}
			svc := newZeroShotTextClassificationService(opts)

			result, err := svc.ClassifyBatch(tc.req)
			require.Error(t, err, tc.description)
			require.Nil(t, result)

			if tc.expectedErrKind != nil {
				// SDK error expected
				testutils.AssertSDKErrorKind(t, err, *tc.expectedErrKind)
				// Verify no request was made for SDK errors
				require.Nil(t, mt.LastRequest)
			} else {
				// API error expected
				var apiErr *hferrors.APIError
				require.ErrorAs(t, err, &apiErr, tc.description)
				require.Equal(t, tc.httpStatusCode, apiErr.StatusCode)
			}
		})
	}
}

func TestZeroShotTextClassificationService_ClassifyBatch_ModelFromOptions(t *testing.T) {
	t.Parallel()

	const zeroShotBatchSingleTestTextResponseBody = `[{"Sequence":"test text","Labels":["positive","negative","neutral"],"Scores":[0.95,0.03,0.02]}]`

	mt := testutils.NewJSONMockTransport(
		http.StatusOK,
		zeroShotBatchSingleTestTextResponseBody,
		nil,
	)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newZeroShotTextClassificationService(opts)

	candidateLabels := []string{"positive", "negative"}
	req := ZeroShotTextClassificationBatchRequest{
		Inputs: []string{"test text"},
		Parameters: &ZeroShotTextClassificationParameters{
			CandidateLabels: candidateLabels,
		},
	}

	result, err := svc.ClassifyBatch(req, WithModel("override-model"))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the correct model was used in the request
	require.NotNil(t, mt.LastRequest)
	require.Contains(t, mt.LastRequest.URL.Path, "override-model")
}
