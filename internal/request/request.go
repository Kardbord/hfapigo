package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Kardbord/hfapigo/v4/internal/errors"
)

// Do performs an HTTP request with the provided options and returns the response.
// It creates a new HTTP request with the given method, path, and body, adds authorization
// and custom headers, and executes the request using the configured transport.
// For HTTP status codes >= 400, it returns an *errors.APIError.
func Do(
	opts RequestOptions,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {
	resp, err := DoRaw(opts, method, path, body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		b, truncated, readErr := readResponseBodyTruncated(resp.Body, opts.MaxResponseBodyBytes)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, &errors.SDKError{
				Kind:    errors.SDKErrorKindInternal,
				Message: "failed to read error response body",
				Err:     readErr,
			}
		}
		msg := string(b)
		if truncated {
			msg = msg + " [truncated]"
		}
		return nil, &errors.APIError{
			StatusCode: resp.StatusCode,
			Message:    msg,
			Body:       b,
			Method:     method,
			URL:        resp.Request.URL.String(),
			RequestID:  resp.Header.Get("X-Request-ID"),
		}
	}

	return resp, nil
}

// DoRaw performs an HTTP request with the provided options and returns the response
// without translating non-2xx status codes into SDK errors.
func DoRaw(
	opts RequestOptions,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {
	if opts.Transport == nil {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: "transport is nil",
		}
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	reqURL, err := url.JoinPath(opts.BaseURL, path)
	if err != nil {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("failed to join base URL %q with path %q", opts.BaseURL, path),
			Err:     err,
		}
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		reqURL,
		body,
	)
	if err != nil {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindInternal,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	// Set standard headers
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}
	if opts.Token != "" {
		req.Header.Set("Authorization", "Bearer "+opts.Token)
	}

	// Set custom headers (can override defaults if needed).
	req.Header = overrideHeaders(req.Header, opts.Headers)

	resp, err := opts.Transport.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindTransport,
			Message: "request failed",
			Err:     err,
		}
	}

	if resp != nil && resp.Request == nil {
		resp.Request = req
	}

	return resp, nil
}

// DoBytes performs an HTTP request with a byte slice body.
// It is a convenience wrapper around Do that converts the byte slice to an io.Reader.
func DoBytes(
	opts RequestOptions,
	method string,
	path string,
	data []byte,
) (*http.Response, error) {
	return Do(opts, method, path, bytes.NewReader(data))
}

// DoBytesRaw performs an HTTP request with a byte slice body and returns the response
// without translating non-2xx status codes into SDK errors.
// It is a convenience wrapper around DoRaw that converts the byte slice to an io.Reader.
func DoBytesRaw(
	opts RequestOptions,
	method string,
	path string,
	data []byte,
) (*http.Response, error) {
	return DoRaw(opts, method, path, bytes.NewReader(data))
}

func readResponseBodyLimited(r io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxResponseBodyBytes
	}
	// LimitReader doesn't error on overflow; it just stops at the limit and returns EOF.
	// Read one extra byte so we can detect truncation by checking len(b) > maxBytes.
	limitReader := io.LimitReader(r, maxBytes+1)
	b, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindInternal,
			Message: "response body exceeds max size",
		}
	}
	return b, nil
}

func readResponseBodyTruncated(r io.Reader, maxBytes int64) ([]byte, bool, error) {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxResponseBodyBytes
	}
	// LimitReader doesn't error on overflow; it just stops at the limit and returns EOF.
	// Read one extra byte so we can detect truncation by checking len(b) > maxBytes.
	limitReader := io.LimitReader(r, maxBytes+1)
	b, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, false, err
	}
	if int64(len(b)) > maxBytes {
		return b[:maxBytes], true, nil
	}
	return b, false, nil
}
