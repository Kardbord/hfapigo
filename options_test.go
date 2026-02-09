package hfapigo

import (
	"errors"
	"net/http"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/request"
)

func TestWithHTTPClientFactoryNil(t *testing.T) {
	opts := request.NewRequestOptions().WithHTTPClientFactory(nil)
	if opts.HTTPClient != nil {
		t.Fatal("expected nil http client when factory returns nil")
	}

	_, err := request.Do(opts, http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error when transport is nil")
	}
	var sdkErr *internalErrors.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}
	if sdkErr.Kind != internalErrors.SDKErrorKindConfiguration {
		t.Fatalf("expected configuration SDKError, got %q", sdkErr.Kind)
	}
}

func TestWithDefaultHTTPClient(t *testing.T) {
	opts := request.NewRequestOptions().WithHTTPClientFactory(nil).With(WithDefaultHTTPClient())
	if opts.HTTPClient == nil {
		t.Fatal("expected default http client, got nil")
	}
}
