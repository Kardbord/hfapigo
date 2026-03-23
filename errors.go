package hfgo

import "github.com/Kardbord/hfgo/v4/internal/hferrors"

// APIError represents an error returned by the HuggingFace API.
// It includes the HTTP status code, error message, response body,
// and request ID if available.
//
// Users can type-assert errors to *APIError to access additional
// error information and helper methods:
//
//	if apiErr, ok := err.(*hfgo.APIError); ok {
//	    if apiErr.IsAuthenticationError() {
//	        // Handle authentication error
//	    }
//	}
type APIError = hferrors.APIError

// SDKError represents a client-side SDK error that occurred before
// a response was received from the API.
//
// Users can type-assert errors to *SDKError to access
// the error kind and underlying cause:
//
//	if sdkErr, ok := err.(*hfgo.SDKError); ok {
//	    fmt.Printf("Kind %s: %s\n", sdkErr.Kind, sdkErr.Message)
//	}
type SDKError = hferrors.SDKError

// SDKErrorKind represents the category of a client-side SDK error.
type SDKErrorKind = hferrors.SDKErrorKind

const (
	// SDKErrorKindValidation indicates a validation error in API responses.
	SDKErrorKindValidation = hferrors.SDKErrorKindValidation
	// SDKErrorKindConfiguration indicates invalid or missing configuration.
	SDKErrorKindConfiguration = hferrors.SDKErrorKindConfiguration
	// SDKErrorKindSerialization indicates a serialization or deserialization error.
	SDKErrorKindSerialization = hferrors.SDKErrorKindSerialization
	// SDKErrorKindTransport indicates a transport-layer failure.
	SDKErrorKindTransport = hferrors.SDKErrorKindTransport
	// SDKErrorKindInternal indicates an internal SDK error.
	SDKErrorKindInternal = hferrors.SDKErrorKindInternal
)
