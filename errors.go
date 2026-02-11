package hfapigo

import "github.com/Kardbord/hfapigo/v4/internal/errors"

// APIError represents an error returned by the HuggingFace API.
// It includes the HTTP status code, error message, response body,
// and request ID if available.
//
// Users can type-assert errors to *APIError to access additional
// error information and helper methods:
//
//	if apiErr, ok := err.(*hfapigo.APIError); ok {
//	    if apiErr.IsAuthenticationError() {
//	        // Handle authentication error
//	    }
//	}
type APIError = errors.APIError

// SDKError represents a client-side SDK error that occurred before
// a response was received from the API.
//
// Users can type-assert errors to *SDKError to access
// the error kind and underlying cause:
//
//	if sdkErr, ok := err.(*hfapigo.SDKError); ok {
//	    fmt.Printf("Kind %s: %s\n", sdkErr.Kind, sdkErr.Message)
//	}
type SDKError = errors.SDKError

// SDKErrorKind represents the category of a client-side SDK error.
type SDKErrorKind = errors.SDKErrorKind

const (
	// SDKErrorKindValidation indicates a validation error in API responses.
	SDKErrorKindValidation = errors.SDKErrorKindValidation
	// SDKErrorKindConfiguration indicates invalid or missing configuration.
	SDKErrorKindConfiguration = errors.SDKErrorKindConfiguration
	// SDKErrorKindSerialization indicates a serialization or deserialization error.
	SDKErrorKindSerialization = errors.SDKErrorKindSerialization
	// SDKErrorKindTransport indicates a transport-layer failure.
	SDKErrorKindTransport = errors.SDKErrorKindTransport
	// SDKErrorKindInternal indicates an internal SDK error.
	SDKErrorKindInternal = errors.SDKErrorKindInternal
)
