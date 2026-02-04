package hfapigo

import (
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

// Do performs an HTTP request and applies SDK error interpretation on non-2xx responses.
func (r RawService) Do(
	requestBody []byte,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return request.DoBytes(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}

// DoRaw performs an HTTP request without translating non-2xx responses into SDK errors.
func (r RawService) DoRaw(
	requestBody []byte,
	method string,
	path string,
	opts ...request.RequestOption,
) (*http.Response, error) {
	return request.DoBytesRaw(
		r.opts.With(opts...),
		method,
		path,
		requestBody,
	)
}
