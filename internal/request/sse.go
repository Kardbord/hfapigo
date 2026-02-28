package request

import (
	"bufio"
	"bytes"
	"context"
	stderrs "errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	internalerrors "github.com/Kardbord/hfapigo/v4/internal/errors"
)

// RawEvent represents a single parsed SSE event payload.
type RawEvent struct {
	// Data holds the concatenated data lines for the event.
	Data []byte
	// Event indicates the optional "event:" field name.
	Event string
	// ID is the optional "id:" field value.
	ID string
	// Retry holds the parsed reconnection duration when present.
	Retry *time.Duration
}

// rawResult wraps either an SSE event or an error from the decoder.
// It stays internal so callers only interact with RawEvent via Recv.
type rawResult struct {
	event RawEvent
	err   error
}

// RawStream consumes an SSE response body and exposes parsed events.
type RawStream struct {
	// results receives decoded events/errors from the background goroutine.
	results <-chan rawResult
	// cancel stops the background decoder context.
	cancel context.CancelFunc
	// body is the underlying HTTP response body being consumed.
	body io.ReadCloser
	// bodyOnce ensures the body is closed exactly once.
	bodyOnce sync.Once
	// closeOnce ensures Close() only shuts down the stream a single time.
	closeOnce sync.Once
	// closeError caches the error returned by the underlying body Close.
	closeError error
}

// StreamRaw starts decoding Server-Sent Events from the provided body.
// Callers should close the returned RawStream when they are done consuming it so
// the background goroutine and HTTP body are released promptly (otherwise it
// will only stop once the server closes the stream).
func StreamRaw(ctx context.Context, body io.ReadCloser) (*RawStream, error) {
	if body == nil {
		return nil, &internalerrors.SDKError{
			Kind:    internalerrors.SDKErrorKindConfiguration,
			Message: "sse: body is nil",
			Err:     nil,
		}
	}

	ctx, cancel := context.WithCancel(NormalizeContext(ctx))
	// Buffered with size 1 so the decoder goroutine can enqueue a single event or
	// error without blocking, while still applying backpressure once the caller
	// falls behind (preventing unbounded buffering).
	results := make(chan rawResult, 1)
	stream := &RawStream{
		results:    results,
		cancel:     cancel,
		body:       body,
		bodyOnce:   sync.Once{},
		closeOnce:  sync.Once{},
		closeError: nil,
	}

	go stream.run(ctx, results)

	return stream, nil
}

// Recv blocks until the next event is available, the provided context is canceled,
// or the stream ends. It returns io.EOF when no more events remain.
func (s *RawStream) Recv(ctx context.Context) (RawEvent, error) {
	if s == nil {
		return RawEvent{}, &internalerrors.SDKError{
			Kind:    internalerrors.SDKErrorKindInternal,
			Message: "sse: stream is nil",
			Err:     nil,
		}
	}

	ctx = NormalizeContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return RawEvent{}, ctx.Err()
		case res, ok := <-s.results:
			if !ok {
				return RawEvent{}, io.EOF
			}
			if res.err != nil {
				return RawEvent{}, res.err
			}

			return res.event, nil
		}
	}
}

// Close stops the decoder and releases the underlying body. It is safe to call multiple times.
func (s *RawStream) Close() error {
	if s == nil {
		return &internalerrors.SDKError{
			Kind:    internalerrors.SDKErrorKindInternal,
			Message: "sse: stream is nil",
			Err:     nil,
		}
	}

	s.closeOnce.Do(func() {
		s.cancel()
		s.bodyOnce.Do(func() {
			if err := s.body.Close(); err != nil {
				s.closeError = wrapStreamError(
					err,
					internalerrors.SDKErrorKindTransport,
					"close stream body",
				)
			}
		})
	})

	return s.closeError
}

// run drives the background goroutine that parses the SSE stream and forwards
// data or errors onto the results channel.
func (s *RawStream) run(ctx context.Context, results chan<- rawResult) {
	defer close(results)
	defer func() { _ = s.Close() }()

	reader := bufio.NewReader(s.body)
	state := &sseState{
		data:      bytes.Buffer{},
		eventType: "",
		eventID:   "",
		retrySet:  false,
		retry:     0,
	}

	for {
		line, err := reader.ReadString('\n')
		hasData := line != ""
		if hasData {
			line = strings.TrimRight(line, "\r\n")
			state.handleLine(line, results)
		}

		if err != nil {
			if stderrs.Is(err, io.EOF) {
				state.emit(results)

				return
			}
			if ctxErr := ctx.Err(); ctxErr != nil {
				err = ctxErr
			}
			// Any data parsed before the error has already been emitted, so
			// callers just observe the terminal error once they drain the queue.
			results <- rawResult{
				event: RawEvent{
					Data:  nil,
					Event: "",
					ID:    "",
					Retry: nil,
				},
				err: wrapStreamError(err, internalerrors.SDKErrorKindTransport, "read sse line"),
			}

			return
		}
	}
}

func parseSSEField(line string) (string, string) {
	// parseSSEField splits an SSE line into its field/value components while
	// trimming the optional leading space defined by the spec.
	const sseFieldSplitParts = 2

	parts := strings.SplitN(line, ":", sseFieldSplitParts)
	field := parts[0]
	if len(parts) == 1 {
		return field, ""
	}
	value := parts[1]
	if value != "" && value[0] == ' ' {
		value = value[1:]
	}

	return field, value
}

// wrapStreamError converts stream read errors into SDK errors, preserving caller
// cancellation errors as-is while retaining the calling context message.
func wrapStreamError(err error, kind internalerrors.SDKErrorKind, contextMsg string) error {
	switch {
	case err == nil:
		return nil
	case stderrs.Is(err, context.Canceled), stderrs.Is(err, context.DeadlineExceeded):
		return err
	default:
		return &internalerrors.SDKError{
			Kind:    kind,
			Message: "sse: failed to " + contextMsg,
			Err:     err,
		}
	}
}

// sseState holds SSE field values between flushes.
type sseState struct {
	data      bytes.Buffer
	eventType string
	eventID   string
	retrySet  bool
	retry     time.Duration
}

// handleLine processes a single SSE line and emits an event if needed.
func (s *sseState) handleLine(line string, results chan<- rawResult) {
	switch {
	case line == "":
		s.emit(results)
	case strings.HasPrefix(line, ":"):
		// comment, ignore
	default:
		field, value := parseSSEField(line)
		switch field {
		case "data":
			s.data.WriteString(value)
			s.data.WriteByte('\n')
		case "event":
			s.eventType = value
		case "id":
			s.eventID = value
		case "retry":
			if ms, err := strconv.Atoi(value); err == nil {
				s.retrySet = true
				s.retry = time.Duration(ms) * time.Millisecond
			}
		}
	}
}

// emit forwards the buffered event (if any) to the results channel.
func (s *sseState) emit(results chan<- rawResult) {
	if s.isEmpty() {
		return
	}

	payload := s.data.Bytes()
	if len(payload) > 0 && payload[len(payload)-1] == '\n' {
		payload = payload[:len(payload)-1]
	}
	data := append([]byte(nil), payload...)

	var retryPtr *time.Duration
	if s.retrySet {
		d := s.retry
		retryPtr = &d
	}

	results <- rawResult{
		event: RawEvent{
			Data:  data,
			Event: s.eventType,
			ID:    s.eventID,
			Retry: retryPtr,
		},
		err: nil,
	}

	s.reset()
}

func (s *sseState) isEmpty() bool {
	return s.data.Len() == 0 && s.eventType == "" && s.eventID == "" && !s.retrySet
}

func (s *sseState) reset() {
	s.data.Reset()
	s.eventType = ""
	s.eventID = ""
	s.retrySet = false
	s.retry = 0
}
