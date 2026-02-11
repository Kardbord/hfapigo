package request

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Kardbord/hfapigo/v4/internal/errors"
)

// Do performs an HTTP request with the provided options and returns the response.
// It creates a new HTTP request with the given method, path, and body, adds authorization
// and custom headers, and executes the request using the configured HTTP client.
// For HTTP status codes >= 400, it returns an *errors.APIError.
// The caller must close resp.Body on success.
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
		if resp.Body == nil || resp.Body == http.NoBody {
			return nil, &errors.APIError{
				StatusCode: resp.StatusCode,
				Message:    "error response body is missing",
				Method:     method,
				URL:        resp.Request.URL.String(),
				RequestID:  resp.Header.Get("X-Request-ID"),
			}
		}
		b, truncated, readErr := readResponseBodyTruncated(resp.Body, opts.MaxResponseBodyBytes)
		if resp.Body != nil {
			drainAndCloseBody(resp.Body)
		}
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
// The caller must close resp.Body on success.
func DoRaw(
	opts RequestOptions,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	reqURL, err := joinURL(opts.BaseURL, path)
	if err != nil {
		var sdkErr *errors.SDKError
		if stderrors.As(err, &sdkErr) {
			return nil, sdkErr
		}
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

	resp, err := opts.HTTPClient.Do(req)
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
	if resp == nil {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindTransport,
			Message: "http client returned nil response without error",
		}
	}

	if resp.Request == nil {
		resp.Request = req
	}
	if resp.Body == nil {
		resp.Body = http.NoBody
	}

	return resp, nil
}

// joinURL combines a base URL with a relative path while preserving query and fragment.
func joinURL(baseURL string, path string) (string, error) {
	if path == "" {
		return baseURL, nil
	}
	parsedPath, err := url.Parse(path)
	if err != nil {
		return "", &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("invalid path %q", path),
			Err:     err,
		}
	}
	if parsedPath.Scheme != "" || parsedPath.Host != "" {
		return "", &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("path must be relative, got %q", path),
		}
	}
	joined, err := url.JoinPath(baseURL, parsedPath.Path)
	if err != nil {
		return "", &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("failed to join base URL %q with path %q", baseURL, parsedPath.Path),
			Err:     err,
		}
	}
	joinedURL, err := url.Parse(joined)
	if err != nil {
		return "", &errors.SDKError{
			Kind:    errors.SDKErrorKindInternal,
			Message: fmt.Sprintf("failed to parse joined URL %q", joined),
			Err:     err,
		}
	}
	joinedURL.RawQuery = parsedPath.RawQuery
	joinedURL.Fragment = parsedPath.Fragment
	return joinedURL.String(), nil
}

// DoBytes performs an HTTP request with a byte slice body.
// It is a convenience wrapper around Do that converts the byte slice to an io.Reader.
// The caller must close resp.Body on success.
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
// The caller must close resp.Body on success.
func DoBytesRaw(
	opts RequestOptions,
	method string,
	path string,
	data []byte,
) (*http.Response, error) {
	return DoRaw(opts, method, path, bytes.NewReader(data))
}

// readResponseBodyLimited reads up to maxBytes and returns an error if the body is larger.
func readResponseBodyLimited(r io.Reader, maxBytes int64) ([]byte, error) {
	b, truncated, err := readResponseBodyTruncated(r, maxBytes)
	if err != nil {
		return nil, err
	}
	if truncated {
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
			Message: fmt.Sprintf("response body exceeds max size (limit %d bytes)", maxBytes),
		}
	}
	return b, nil
}

// readResponseBodyTruncated reads up to maxBytes and reports if truncation occurred.
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

// drainAndCloseBody drains any remaining data and closes the body.
func drainAndCloseBody(body io.ReadCloser) {
	if body == nil || body == http.NoBody {
		return
	}
	// Drain the remainder so the underlying HTTP connection can be reused.
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
