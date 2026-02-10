package request

import (
	"net/http"

	"github.com/Kardbord/hfapigo/v4/internal/testutils"
)

func withMockTransport(opts RequestOptions, mt *testutils.MockTransport) RequestOptions {
	return opts.WithHTTPClientFactory(func() http.Client {
		return testutils.NewMockHTTPClient(mt)
	})
}
