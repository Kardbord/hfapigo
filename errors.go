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

// ValidationError represents an error that occurs when validating
// request parameters or configuration options.
//
// Users can type-assert errors to *ValidationError to access
// the specific field that failed validation:
//
//	if valErr, ok := err.(*hfapigo.ValidationError); ok {
//	    fmt.Printf("Field %s failed validation: %s\n", valErr.Field, valErr.Message)
//	}
type ValidationError = errors.ValidationError
