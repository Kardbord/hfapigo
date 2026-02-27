package request

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	internalerrors "github.com/Kardbord/hfapigo/v4/internal/errors"
)

func TestStreamRaw_BasicEvents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		payload  string
		validate func(t *testing.T, ev RawEvent)
	}{
		{
			name:    "single event",
			payload: "data: hello\n\n",
			validate: func(t *testing.T, ev RawEvent) {
				if string(ev.Data) != "hello" {
					t.Fatalf("unexpected data: %q", string(ev.Data))
				}
			},
		},
		{
			name:    "multi line event",
			payload: "data: foo\ndata: bar\n\n",
			validate: func(t *testing.T, ev RawEvent) {
				if string(ev.Data) != "foo\nbar" {
					t.Fatalf("unexpected data: %q", string(ev.Data))
				}
			},
		},
		{
			name:    "metadata event",
			payload: "event: chunk\nid: 42\nretry: 1000\ndata: hi\n\n",
			validate: func(t *testing.T, ev RawEvent) {
				if ev.Event != "chunk" || ev.ID != "42" {
					t.Fatalf("unexpected metadata: %+v", ev)
				}
				if ev.Retry == nil || *ev.Retry != time.Second {
					t.Fatalf("unexpected retry: %+v", ev.Retry)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body := io.NopCloser(strings.NewReader(tc.payload))
			stream, err := StreamRaw(context.Background(), body)
			if err != nil {
				t.Fatalf("StreamRaw: %v", err)
			}
			defer func() { _ = stream.Close() }()

			ev, err := stream.Recv(context.Background())
			if err != nil {
				t.Fatalf("Recv: %v", err)
			}
			tc.validate(t, ev)

			if _, err := stream.Recv(context.Background()); !errors.Is(err, io.EOF) {
				t.Fatalf("expected EOF, got %v", err)
			}
		})
	}
}

func TestStreamRaw_CloseCancels(t *testing.T) {
	t.Parallel()

	bodyReader, bodyWriter := io.Pipe()
	stream, err := StreamRaw(context.Background(), bodyReader)
	if err != nil {
		t.Fatalf("StreamRaw: %v", err)
	}

	// Done channel to ensure the goroutine below completes before
	// this function exits.
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		_ = stream.Close()
	}()

	_, err = stream.Recv(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}

	<-done
	_ = bodyWriter.Close()
}

func TestStreamRaw_Errors(t *testing.T) {
	t.Parallel()

	t.Run("nil body", func(t *testing.T) {
		t.Parallel()

		stream, err := StreamRaw(context.Background(), nil)
		if stream != nil || err == nil {
			t.Fatalf("expected error, got stream=%v err=%v", stream, err)
		}
		var sdkErr *internalerrors.SDKError
		if !errors.As(err, &sdkErr) || sdkErr.Kind != internalerrors.SDKErrorKindConfiguration {
			t.Fatalf("expected configuration SDK error, got %#v", err)
		}
	})

	t.Run("body close error", func(t *testing.T) {
		t.Parallel()

		body := errorReadCloser{
			Reader: strings.NewReader("data: hi\n\n"),
			Err:    errors.New("boom"),
		}
		stream, err := StreamRaw(context.Background(), body)
		if err != nil {
			t.Fatalf("unexpected error creating stream: %v", err)
		}
		if err := stream.Close(); err == nil {
			t.Fatalf("expected close error")
		} else {
			var sdkErr *internalerrors.SDKError
			if !errors.As(err, &sdkErr) || sdkErr.Kind != internalerrors.SDKErrorKindTransport {
				t.Fatalf("expected transport SDK error, got %#v", err)
			}
		}
	})
}

type errorReadCloser struct {
	io.Reader
	Err error
}

func (e errorReadCloser) Close() error {
	return e.Err
}
