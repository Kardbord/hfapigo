package testutils

import (
	"encoding/json"
	"testing"
)

// AssertUnmarshalError unmarshals JSON data and asserts the error expectation.
func AssertUnmarshalError[T any](t *testing.T, data []byte, wantErr bool) {
	t.Helper()
	var out T
	err := json.Unmarshal(data, &out)
	AssertError(t, err, wantErr)
}
