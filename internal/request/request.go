package request

import (
	"bytes"
	"io"
	"net/http"
)

func Do(
	opts *RequestOptions,
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

func DoBytes(
	opts *RequestOptions,
	method string,
	path string,
	data []byte,
	headers map[string]string,
) (*http.Response, error) {
	return Do(opts, method, path, bytes.NewReader(data), headers)
}
