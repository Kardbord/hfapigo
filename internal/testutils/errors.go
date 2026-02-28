package testutils

import (
	"errors"
	"testing"

	"github.com/Kardbord/hfapigo/v4/internal/hferrors"
)

// AssertSDKErrorKind fails the test if err is not an SDKError of the expected kind.
func AssertSDKErrorKind(t *testing.T, err error, want hferrors.SDKErrorKind) {
	t.Helper()

	var sdkErr *hferrors.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}
	if sdkErr.Kind != want {
		t.Fatalf("expected SDKError kind %q, got %q", want, sdkErr.Kind)
	}
}

// AssertAPIErrorStatus fails the test if err is not an APIError with the expected status.
// It returns the APIError for additional assertions.
func AssertAPIErrorStatus(t *testing.T, err error, want int) *hferrors.APIError {
	t.Helper()

	var apiErr *hferrors.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != want {
		t.Fatalf("expected status %d, got %d", want, apiErr.StatusCode)
	}

	return apiErr
}
