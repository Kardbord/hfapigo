package hfapigo

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/request"
)

// RawService provides methods for sending raw HTTP requests, with optional SDK error interpretation.
type RawService struct {
	opts request.RequestOptions
}

func newRawService(opts request.RequestOptions) RawService {
	return RawService{opts: opts}
}

// Do performs a raw HTTP request with a byte slice body and applies SDK error interpretation on non-2xx responses.
func (r RawService) Do(
	requestBody []byte,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return r.DoReader(bytes.NewReader(requestBody), method, path, opts...)
}

// DoRaw performs a raw HTTP request with a byte slice body without translating non-2xx responses into SDK errors.
func (r RawService) DoRaw(
	requestBody []byte,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return r.DoRawReader(bytes.NewReader(requestBody), method, path, opts...)
}

// DoReader performs a raw HTTP request with a streaming body and applies SDK error interpretation on non-2xx responses.
func (r RawService) DoReader(
	requestBody io.Reader,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return request.Do(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}

// DoRawReader performs a raw HTTP request with a streaming body without translating non-2xx responses into SDK errors.
func (r RawService) DoRawReader(
	requestBody io.Reader,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return request.DoRaw(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}
