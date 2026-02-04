package hfapigo

import (
	"errors"
	"net/http"
	"testing"

	internalErrors "github.com/Kardbord/hfapigo/v4/internal/errors"
	"github.com/Kardbord/hfapigo/v4/internal/request"
)

func TestWithHTTPClientNil(t *testing.T) {
	opts := request.NewRequestOptions().WithHTTPClient(nil)
	if opts.Transport != nil {
		t.Fatal("expected nil transport when http client is nil")
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
