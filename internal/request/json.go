package request

import (
	"encoding/json"
	"fmt"
	"io"
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
// or the response status code is 300 or greater.
func DoJSON[TReq any, TResp any](
	opts RequestOptions,
	method string,
	path string,
	reqBody TReq,
) (TResp, error) {

	var zero TResp

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return zero, err
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

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("hf api error (%d): %s", resp.StatusCode, string(b))
	}

	var out TResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return zero, err
	}

	return out, nil
}
