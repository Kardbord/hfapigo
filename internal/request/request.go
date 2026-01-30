package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	reqURL, err := url.JoinPath(opts.BaseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to join base URL %q with path %q: %w", opts.BaseURL, path, err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		reqURL,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set standard headers
	req.Header.Set("User-Agent", opts.UserAgent)
	if opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+opts.Token)
	}

	// Set custom headers (can override defaults if needed)
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
