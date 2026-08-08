package testutils

import (
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/stretchr/testify/require"
)

// WantErr selects which SDK error type a test expects from an API call.
type WantErr int

const (
	// WantErrSDK indicates that a test should expect an SDKError.
	WantErrSDK WantErr = iota

	// WantErrAPI indicates that a test should expect an APIError.
	WantErrAPI
)

// AssertErrorType asserts err is an SDKError of sdkErrKind (when want is
// WantErrSDK) or an APIError with the given status (when want is WantErrAPI).
func AssertErrorType(
	t *testing.T,
	err error,
	want WantErr,
	sdkErrKind hferrors.SDKErrorKind,
	apiStatus int,
) {
	t.Helper()

	switch want {
	case WantErrSDK:
		AssertSDKErrorKind(t, err, sdkErrKind)
	case WantErrAPI:
		AssertAPIErrorStatus(t, err, apiStatus)
	default:
		t.Fatalf("unknown error expectation: %v", want)
	}
}

// AssertSDKErrorKind fails the test if err is not an SDKError of the expected kind.
func AssertSDKErrorKind(t *testing.T, err error, want hferrors.SDKErrorKind) {
	t.Helper()

	var sdkErr *hferrors.SDKError
	require.ErrorAs(t, err, &sdkErr)
	require.Equal(t, want, sdkErr.Kind)
}

// AssertAPIErrorStatus fails the test if err is not an APIError with the expected status.
// It returns the APIError for additional assertions.
func AssertAPIErrorStatus(t *testing.T, err error, want int) *hferrors.APIError {
	t.Helper()

	var apiErr *hferrors.APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, want, apiErr.StatusCode)

	return apiErr
}
