package errors

import (
	"net/http"
	"strings"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name       string
		err        *APIError
		wantSubstr []string
	}{
		{
			name: "basic error message",
			err: &APIError{
				StatusCode: http.StatusBadRequest,
				Message:    "Bad Request",
			},
			wantSubstr: []string{"400", "Bad Request"},
		},
		{
			name: "error with request ID",
			err: &APIError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				RequestID:  "test-req-12345",
			},
			wantSubstr: []string{"500", "Internal Server Error", "test-req-12345"},
		},
		{
			name: "error with URL and method",
			err: &APIError{
				StatusCode: http.StatusNotFound,
				Message:    "Not Found",
			},
			wantSubstr: []string{"404", "Not Found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("APIError.Error() = %q, want to contain %q", got, substr)
				}
			}
		})
	}
}

func TestAPIError_IsClientError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"400 is client error", http.StatusBadRequest, true},
		{"404 is client error", http.StatusNotFound, true},
		{"499 is client error", 499, true},
		{"200 is not client error", http.StatusOK, false},
		{"500 is not client error", http.StatusInternalServerError, false},
		{"300 is not client error", http.StatusMultipleChoices, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			if got := err.IsClientError(); got != tt.want {
				t.Errorf("APIError.IsClientError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"500 is server error", http.StatusInternalServerError, true},
		{"503 is server error", http.StatusServiceUnavailable, true},
		{"599 is server error", 599, true},
		{"400 is not server error", http.StatusBadRequest, false},
		{"200 is not server error", http.StatusOK, false},
		{"600 is not server error", 600, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			if got := err.IsServerError(); got != tt.want {
				t.Errorf("APIError.IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsAuthenticationError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"401 is authentication error", http.StatusUnauthorized, true},
		{"400 is not authentication error", http.StatusBadRequest, false},
		{"403 is not authentication error", http.StatusForbidden, false},
		{"500 is not authentication error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			if got := err.IsAuthenticationError(); got != tt.want {
				t.Errorf("APIError.IsAuthenticationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsRateLimitError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"429 is rate limit error", http.StatusTooManyRequests, true},
		{"400 is not rate limit error", http.StatusBadRequest, false},
		{"500 is not rate limit error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			if got := err.IsRateLimitError(); got != tt.want {
				t.Errorf("APIError.IsRateLimitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name       string
		err        *ValidationError
		wantSubstr []string
	}{
		{
			name: "basic validation error",
			err: &ValidationError{
				Field:   "Token",
				Message: "cannot be empty",
			},
			wantSubstr: []string{"Token", "cannot be empty", "validation error"},
		},
		{
			name: "validation error with special characters",
			err: &ValidationError{
				Field:   "BaseURL",
				Message: "must be a valid URL",
			},
			wantSubstr: []string{"BaseURL", "must be a valid URL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("ValidationError.Error() = %q, want to contain %q", got, substr)
				}
			}
		})
	}
}
