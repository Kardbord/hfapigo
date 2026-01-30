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
	SDKErrorKindValidation    = errors.SDKErrorKindValidation
	SDKErrorKindConfiguration = errors.SDKErrorKindConfiguration
	SDKErrorKindSerialization = errors.SDKErrorKindSerialization
	SDKErrorKindTransport     = errors.SDKErrorKindTransport
	SDKErrorKindInternal      = errors.SDKErrorKindInternal
)
