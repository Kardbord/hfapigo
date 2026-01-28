package request

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Kardbord/hfapigo/v4/internal/errors"
)

// DoJSON performs an HTTP request with a JSON request body and expects a JSON response.
// It marshals the request body to JSON, sends the request, and unmarshals the response
// into the specified response type. The function uses Go generics to provide type-safe
// request and response handling.
//
// Type parameters:
//   - TReq: The type of the request body
//   - TResp: The type of the response body
//
// Returns an error if JSON marshaling/unmarshaling fails, the HTTP request fails,
// or the response status code is 400 or greater. For HTTP errors, returns an *errors.APIError
// which includes the status code, response body, and other metadata.
func DoJSON[TReq any, TResp any](
	opts RequestOptions,
	method string,
	path string,
	reqBody TReq,
) (TResp, error) {

	var zero TResp

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := DoBytes(
		opts,
		method,
		path,
		buf,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return zero, fmt.Errorf("failed to execute request to %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	// Check for HTTP error status codes (4xx and 5xx)
	if resp.StatusCode >= 400 {
		b, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return zero, &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("failed to read error response body: %v", readErr),
				Method:     method,
				URL:        opts.BaseURL + path,
				RequestID:  resp.Header.Get("X-Request-ID"),
			}
		}
		return zero, &errors.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(b),
			Body:       b,
			Method:     method,
			URL:        opts.BaseURL + path,
			RequestID:  resp.Header.Get("X-Request-ID"),
		}
	}

	var out TResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return zero, fmt.Errorf("failed to decode response body: %w", err)
	}

	return out, nil
}
