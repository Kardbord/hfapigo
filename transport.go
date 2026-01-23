package hfapigo

import "net/http"

type Transport interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPTransport struct {
	client *http.Client
}

func (t *HTTPTransport) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}
