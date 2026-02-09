package hfapigo

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// RawService provides methods for sending raw HTTP requests, with optional SDK error interpretation.
type rawService struct {
	opts request.RequestOptions
}

func newRawService(opts request.RequestOptions) rawService {
	return rawService{opts: opts}
}

// Do performs a raw HTTP request with a byte slice body and applies SDK error interpretation on non-2xx responses.
// The caller must close resp.Body on success.
func (r rawService) Do(
	requestBody []byte,
	method string,
	path string,
	opts ...RequestOption,
) (*http.Response, error) {
	return r.DoReader(bytes.NewReader(requestBody), method, path, opts...)
}

// DoRaw performs a raw HTTP request with a byte slice body without translating non-2xx responses into SDK errors.
// The caller must close resp.Body on success.
func (r rawService) DoRaw(
	requestBody []byte,
	method string,
	path string,
	opts ...RequestOption,
) (*http.Response, error) {
	return r.DoRawReader(bytes.NewReader(requestBody), method, path, opts...)
}

// DoReader performs a raw HTTP request with a streaming body and applies SDK error interpretation on non-2xx responses.
// The caller must close resp.Body on success.
func (r rawService) DoReader(
	requestBody io.Reader,
	method string,
	path string,
	opts ...RequestOption,
) (*http.Response, error) {
	return request.Do(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}

// DoRawReader performs a raw HTTP request with a streaming body without translating non-2xx responses into SDK errors.
// The caller must close resp.Body on success.
func (r rawService) DoRawReader(
	requestBody io.Reader,
	method string,
	path string,
	opts ...RequestOption,
) (*http.Response, error) {
	return request.DoRaw(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}
