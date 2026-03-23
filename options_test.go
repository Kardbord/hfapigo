package hfgo

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/request"
	"github.com/stretchr/testify/require"
)

func TestWithHTTPClientFactoryNil(t *testing.T) {
	t.Parallel()

	opts := request.NewOptions().WithHTTPClientFactory(nil)
	if opts.HTTPClient != nil {
		t.Fatal("expected nil http client when factory returns nil")
	}

	_, err := request.Do(opts, http.MethodGet, "/test", nil)
	require.Error(t, err)
	var sdkErr *hferrors.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}
	if sdkErr.Kind != hferrors.SDKErrorKindConfiguration {
		t.Fatalf("expected configuration SDKError, got %q", sdkErr.Kind)
	}
}

func TestWithDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	opts := request.NewOptions().WithHTTPClientFactory(nil).With(WithDefaultHTTPClient())
	if opts.HTTPClient == nil {
		t.Fatal("expected default http client, got nil")
	}
}
