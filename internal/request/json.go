package request

import (
	"encoding/json"
	stderrors "errors"
	"io"
	"mime"
	"net/http"
	"strings"

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
// Returns an error if JSON marshaling/unmarshaling fails or the HTTP request fails.
// For HTTP errors, Do returns an *errors.APIError which includes the status code,
// response body, and other metadata.
func DoJSON[TReq any, TResp any](
	opts RequestOptions,
	method string,
	path string,
	reqBody TReq,
) (TResp, error) {

	var zero TResp

	buf, err := json.Marshal(reqBody)
	if err != nil {
		var sdkErr *errors.SDKError
		if stderrors.As(err, &sdkErr) {
			return zero, sdkErr
		}
		return zero, &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "failed to marshal request body",
			Err:     err,
		}
	}

	opts = opts.WithDefaultHeader("Content-Type", "application/json")
	opts = opts.WithDefaultHeader("Accept", "application/json")

	if err := validateJSONRequestContentType(opts.Headers); err != nil {
		return zero, err
	}

	resp, err := DoBytes(
		opts,
		method,
		path,
		buf,
	)
	if err != nil {
		return zero, err
	}
	defer drainAndCloseBody(resp.Body)

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusResetContent {
		return zero, nil
	}
	if err := validateJSONResponseContentType(resp.Header); err != nil {
		return zero, err
	}

	var out TResp
	body, err := readResponseBodyLimited(resp.Body, opts.MaxResponseBodyBytes)
	if err != nil {
		return zero, err
	}
	if len(body) == 0 {
		return zero, &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "empty response body",
		}
	}
	if err := json.Unmarshal(body, &out); err != nil {
		if stderrors.Is(err, io.EOF) {
			return zero, &errors.SDKError{
				Kind:    errors.SDKErrorKindSerialization,
				Message: "empty response body",
			}
		}
		return zero, &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "failed to decode response body",
			Err:     err,
		}
	}

	return out, nil
}

// ensureHeader returns a copy of headers with a default value set when missing or empty.
func ensureHeader(h http.Header, key string, value string) http.Header {
	out := cloneHeader(h)
	if out == nil {
		out = make(http.Header, 1)
	}
	if v := out.Get(key); v == "" {
		out.Set(key, value)
	}
	return out
}

// validateJSONRequestContentType validates that Content-Type is application/json when provided.
func validateJSONRequestContentType(headers http.Header) error {
	ct := headers.Get("Content-Type")
	if ct == "" {
		return nil
	}
	mediatype, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "invalid Content-Type header",
			Err:     err,
		}
	}
	if mediatype != "application/json" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "Content-Type must be application/json for DoJSON requests",
		}
	}
	return nil
}

// validateJSONResponseContentType validates that the response Content-Type indicates JSON.
func validateJSONResponseContentType(headers http.Header) error {
	ct := headers.Get("Content-Type")
	if ct == "" {
		return nil
	}
	mediatype, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "invalid Content-Type header on response",
			Err:     err,
		}
	}
	if !isJSONMediaType(mediatype) {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "response Content-Type must be application/json",
		}
	}
	return nil
}

// isJSONMediaType reports whether the media type is JSON or a +json subtype.
func isJSONMediaType(mediatype string) bool {
	if mediatype == "application/json" {
		return true
	}
	return strings.HasSuffix(mediatype, "+json")
}
