package testutils

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadRequestBody reads and closes the last captured request body from mt and
// unmarshals it as JSON into a map. It fails the test on any error.
func ReadRequestBody(t *testing.T, mt *MockTransport) map[string]any {
	t.Helper()

	require.NotNil(t, mt.LastRequest)
	body, err := io.ReadAll(mt.LastRequest.Body)
	require.NoError(t, err)
	require.NoError(t, mt.LastRequest.Body.Close())

	var reqBody map[string]any
	require.NoError(t, json.Unmarshal(body, &reqBody))

	return reqBody
}
