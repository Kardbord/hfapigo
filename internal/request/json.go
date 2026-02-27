package request

import (
	"bytes"
	"context"
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

	buf, err := marshalJSONRequestBody(reqBody)
	if err != nil {
		return zero, err
	}

	opts, err = prepareJSONRequestOptions(opts, "application/json")
	if err != nil {
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

// DoJSONStream performs an HTTP request with a JSON body and returns a streaming JSON response.
// The response body must be a Server-Sent Events (SSE) stream where each data chunk contains JSON.
// Callers are responsible for closing the returned stream to release resources.
func DoJSONStream[TReq any, TResp any](
	opts RequestOptions,
	method string,
	path string,
	reqBody TReq,
) (*JSONStream[TResp], error) {
	buf, err := marshalJSONRequestBody(reqBody)
	if err != nil {
		return nil, err
	}

	opts, err = prepareJSONRequestOptions(opts, "text/event-stream")
	if err != nil {
		return nil, err
	}

	resp, err := DoBytes(
		opts,
		method,
		path,
		buf,
	)
	if err != nil {
		return nil, err
	}

	if err := validateEventStreamResponseContentType(resp.Header); err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	raw, err := StreamRaw(ctx, resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	return &JSONStream[TResp]{raw: raw}, nil
}

// JSONStream consumes JSON SSE events produced by DoJSONStream.
type JSONStream[T any] struct {
	raw *RawStream
}

// Recv blocks until the next JSON event is available or the stream ends.
// It skips keepalive events, treats data: [DONE] as EOF, and unmarshals each chunk into T.
func (s *JSONStream[T]) Recv(ctx context.Context) (T, error) {
	var zero T
	if s == nil || s.raw == nil {
		return zero, &errors.SDKError{
			Kind:    errors.SDKErrorKindInternal,
			Message: "json stream is nil",
		}
	}

	for {
		event, err := s.raw.Recv(ctx)
		if err != nil {
			return zero, err
		}
		data := bytes.TrimSpace(event.Data)
		if len(data) == 0 {
			continue
		}
		if bytes.Equal(data, []byte("[DONE]")) {
			_ = s.raw.Close()
			return zero, io.EOF
		}
		var out T
		if err := json.Unmarshal(data, &out); err != nil {
			return zero, &errors.SDKError{
				Kind:    errors.SDKErrorKindSerialization,
				Message: "failed to decode stream event",
				Err:     err,
			}
		}

		return out, nil
	}
}

// Close releases the underlying stream resources.
func (s *JSONStream[T]) Close() error {
	if s == nil || s.raw == nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindInternal,
			Message: "json stream is nil",
		}
	}

	return s.raw.Close()
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
			Kind:    errors.SDKErrorKindConfiguration,
			Message: "invalid Content-Type header",
			Err:     err,
		}
	}
	if mediatype != "application/json" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindConfiguration,
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

// validateEventStreamResponseContentType ensures the response advertises text/event-stream.
func validateEventStreamResponseContentType(headers http.Header) error {
	ct := headers.Get("Content-Type")
	if ct == "" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "response Content-Type must be text/event-stream",
		}
	}
	mediatype, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "invalid Content-Type header on response",
			Err:     err,
		}
	}
	if mediatype != "text/event-stream" {
		return &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "response Content-Type must be text/event-stream",
		}
	}
	return nil
}

// marshalJSONRequestBody serializes the payload and normalizes errors to SDK errors.
func marshalJSONRequestBody(payload any) ([]byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		var sdkErr *errors.SDKError
		if stderrors.As(err, &sdkErr) {
			return nil, sdkErr
		}
		return nil, &errors.SDKError{
			Kind:    errors.SDKErrorKindSerialization,
			Message: "failed to marshal request body",
			Err:     err,
		}
	}

	return buf, nil
}

// prepareJSONRequestOptions sets standard headers and validates Content-Type for JSON requests.
func prepareJSONRequestOptions(opts RequestOptions, accept string) (RequestOptions, error) {
	opts = opts.WithDefaultHeader("Content-Type", "application/json")
	opts = opts.WithDefaultHeader("Accept", accept)
	if err := validateJSONRequestContentType(opts.Headers); err != nil {
		return RequestOptions{}, err
	}

	return opts, nil
}

// isJSONMediaType reports whether the media type is JSON or a +json subtype.
func isJSONMediaType(mediatype string) bool {
	if mediatype == "application/json" {
		return true
	}
	return strings.HasSuffix(mediatype, "+json")
}
