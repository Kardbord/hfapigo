package testutils

import (
	"errors"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
)

// AssertSDKErrorKind fails the test if err is not an SDKError of the expected kind.
func AssertSDKErrorKind(t *testing.T, err error, want internalErrors.SDKErrorKind) {
	t.Helper()

	var sdkErr *internalErrors.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}
	if sdkErr.Kind != want {
		t.Fatalf("expected SDKError kind %q, got %q", want, sdkErr.Kind)
	}
}

// AssertAPIErrorStatus fails the test if err is not an APIError with the expected status.
// It returns the APIError for additional assertions.
func AssertAPIErrorStatus(t *testing.T, err error, want int) *internalErrors.APIError {
	t.Helper()

	var apiErr *internalErrors.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != want {
		t.Fatalf("expected status %d, got %d", want, apiErr.StatusCode)
	}
	return apiErr
}
