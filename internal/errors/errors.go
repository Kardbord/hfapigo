package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the HuggingFace API.
// It includes the HTTP status code, error message, response body,
// and request ID if available.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API
	StatusCode int

	// Message is the human-readable error message
	Message string

	// Body contains the raw response body from the API
	Body []byte

	// RequestID is the request identifier from the X-Request-ID header, if available
	RequestID string

	// Method is the HTTP method used for the request
	Method string

	// URL is the URL that was requested
	URL string
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("hf api error (%d): %s [request_id: %s]", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("hf api error (%d): %s", e.StatusCode, e.Message)
}

// IsClientError returns true if the error is a 4xx client error.
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsAuthenticationError returns true if the error is a 401 Unauthorized error.
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsRateLimitError returns true if the error is a 429 Too Many Requests error.
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// ValidationError represents an error that occurs when validating
// request parameters or configuration options.
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string

	// Message is the human-readable error message
	Message string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
