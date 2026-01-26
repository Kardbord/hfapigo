package request

import "net/http"

type Transport interface {
	Do(*http.Request) (*http.Response, error)
}

type httpTransport struct {
	client *http.Client
}

func NewHTTPTransport(c *http.Client) Transport {
	return &httpTransport{client: c}
}

func (t *httpTransport) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}
