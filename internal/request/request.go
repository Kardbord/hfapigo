package request

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func Do(
	ctx *RequestOptions,
	method string,
	path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {

	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		ctx.BaseURL+path,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ctx.Token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return ctx.Transport.Do(req)
}

func DoBytes(
	ctx *RequestOptions,
	method string,
	path string,
	data []byte,
	headers map[string]string,
) (*http.Response, error) {
	return Do(ctx, method, path, bytes.NewReader(data), headers)
}
