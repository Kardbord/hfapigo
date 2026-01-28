package request

import (
	"bytes"
	"io"
	"net/http"
)

// Do performs an HTTP request with the provided options and returns the response.
// It creates a new HTTP request with the given method, path, and body, adds authorization
// and custom headers, and executes the request using the configured transport.
func Do(
	opts RequestOptions,
	method string,
	path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {

	req, err := http.NewRequestWithContext(
		opts.Ctx,
		method,
		opts.BaseURL+path,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+opts.Token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return opts.Transport.Do(req)
}

// DoBytes performs an HTTP request with a byte slice body.
// It is a convenience wrapper around Do that converts the byte slice to an io.Reader.
func DoBytes(
	opts RequestOptions,
	method string,
	path string,
	data []byte,
	headers map[string]string,
) (*http.Response, error) {
	return Do(opts, method, path, bytes.NewReader(data), headers)
}
