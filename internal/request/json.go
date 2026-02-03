package request

import (
	"encoding/json"
	stderrors "errors"
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
		return zero, &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "failed to marshal request body",
			Err:     err,
		}
	}

	resp, err := DoBytes(
		opts,
		method,
		path,
		buf,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

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
