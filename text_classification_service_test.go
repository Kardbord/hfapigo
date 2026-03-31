//go:build !integration

package hfgo

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/request"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

const textClassificationResponseBody = `[[{"label":"positive","score":0.95}]]`

func TestTextClassificationService_Classify_SingleInput(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.Classify(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 1)
	require.Equal(t, "positive", result[0].Label)
	require.InEpsilon(t, 0.95, result[0].Score, 0.001)
}

func TestTextClassificationService_Classify_WithParameters(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	topK := 2
	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
		Parameters: &TextClassificationParameters{
			TopK: &topK,
		},
	}

	result, err := svc.Classify(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the request was made correctly
	require.NotNil(t, mt.LastRequest)
	body, err := io.ReadAll(mt.LastRequest.Body)
	require.NoError(t, err)
	_ = mt.LastRequest.Body.Close()

	var reqBody map[string]any
	err = json.Unmarshal(body, &reqBody)
	require.NoError(t, err)

	params, ok := reqBody["parameters"].(map[string]any)
	require.True(t, ok, "parameters should be a map")
	require.InEpsilon(t, float64(2), params["top_k"], 0.001)
}

func TestTextClassificationService_Classify_MultipleInputsError(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"text1", "text2"},
	}

	result, err := svc.Classify(req)
	require.Error(t, err)
	require.Nil(t, result)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)

	// Verify no request was made
	require.Nil(t, mt.LastRequest)
}

func TestTextClassificationService_ClassifyBatch_SingleInput(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.ClassifyBatch(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 1)
	require.Len(t, result[0], 1)
	require.Equal(t, "positive", result[0][0].Label)
	require.InEpsilon(t, 0.95, result[0][0].Score, 0.001)
}

func TestTextClassificationService_ClassifyBatch_MultipleInputs(t *testing.T) {
	t.Parallel()

	batchResponseBody := `[[{"label":"positive","score":0.95}],[{"label":"negative","score":0.87}],[{"label":"neutral","score":0.75}]]`
	mt := testutils.NewJSONMockTransport(http.StatusOK, batchResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"text1", "text2", "text3"},
	}

	result, err := svc.ClassifyBatch(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 3)

	require.Equal(t, "positive", result[0][0].Label)
	require.InEpsilon(t, 0.95, result[0][0].Score, 0.001)

	require.Equal(t, "negative", result[1][0].Label)
	require.InEpsilon(t, 0.87, result[1][0].Score, 0.001)

	require.Equal(t, "neutral", result[2][0].Label)
	require.InEpsilon(t, 0.75, result[2][0].Score, 0.001)
}

func TestTextClassificationService_ClassifyBatch_NoModel(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.ClassifyBatch(req)
	require.Error(t, err)
	require.Nil(t, result)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindConfiguration)

	// Verify no request was made
	require.Nil(t, mt.LastRequest)
}

func TestTextClassificationService_ClassifyBatch_ModelFromOptions(t *testing.T) {
	t.Parallel()

	mt := testutils.NewJSONMockTransport(http.StatusOK, textClassificationResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) })
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.ClassifyBatch(req, WithModel("override-model"))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the correct model was used in the request
	require.NotNil(t, mt.LastRequest)
	require.Contains(t, mt.LastRequest.URL.Path, "override-model")
}

func TestTextClassificationService_ClassifyBatch_APIError(t *testing.T) {
	t.Parallel()

	errorResponse := `{"error":"Model not found"}`
	mt := testutils.NewJSONMockTransport(http.StatusNotFound, errorResponse, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("nonexistent-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.ClassifyBatch(req)
	require.Error(t, err)
	require.Nil(t, result)

	// Should be an API error, not SDK error
	var apiErr *hferrors.APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}

func TestTextClassificationService_ClassifyBatch_EmptyResponse(t *testing.T) {
	t.Parallel()

	emptyResponseBody := `[[]]`
	mt := testutils.NewJSONMockTransport(http.StatusOK, emptyResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"test text"},
	}

	result, err := svc.Classify(req)
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestTextClassificationService_ClassifyBatch_MultipleClassificationsPerInput(t *testing.T) {
	t.Parallel()

	multiResponseBody := `[[{"label":"positive","score":0.95},{"label":"negative","score":0.05}],[{"label":"negative","score":0.87},{"label":"positive","score":0.13}]]`
	mt := testutils.NewJSONMockTransport(http.StatusOK, multiResponseBody, nil)
	opts := request.NewOptions().
		WithHTTPClientFactory(func() http.Client { return testutils.NewMockHTTPClient(mt) }).
		WithModel("test-model")
	svc := newTextClassificationService(opts)

	topK := 2
	req := TextClassificationRequest{
		Inputs: TextClassificationInput{"text1", "text2"},
		Parameters: &TextClassificationParameters{
			TopK: &topK,
		},
	}

	result, err := svc.ClassifyBatch(req)
	require.NoError(t, err)
	require.Len(t, result, 2)

	// First input classifications
	require.Len(t, result[0], 2)
	require.Equal(t, "positive", result[0][0].Label)
	require.InEpsilon(t, 0.95, result[0][0].Score, 0.001)
	require.Equal(t, "negative", result[0][1].Label)
	require.InEpsilon(t, 0.05, result[0][1].Score, 0.001)

	// Second input classifications
	require.Len(t, result[1], 2)
	require.Equal(t, "negative", result[1][0].Label)
	require.InEpsilon(t, 0.87, result[1][0].Score, 0.001)
	require.Equal(t, "positive", result[1][1].Label)
	require.InEpsilon(t, 0.13, result[1][1].Score, 0.001)
}
