# AGENTS.md - hfgo Project Guide

## Project Overview

**hfgo** is a production-quality Go SDK for the [Hugging Face Inference API](https://huggingface.co/docs/inference-providers/tasks/index). It provides Go bindings to perform inference tasks for any supported model available in the Hugging Face Model Hub.

- **Module**: `github.com/Kardbord/hfgo/v4`
- **Go Version**: 1.25+
- **License**: MIT (Copyright 2021 Tanner Kvarfordt)
- **Goal**: Production-ready, follows best practices and idioms, maintains feature parity with upstream API
- **Repository**: https://github.com/Kardbord/hfgo

## Core Architecture

### Client-Centric Design Pattern

The SDK follows a strict immutability pattern for concurrency safety:

1. **Client**: Immutable value type that captures configuration at creation time
   - Options are fixed and never mutated
   - Safe for concurrent use across goroutines
   - Services capture a snapshot of client options when created
   - Lightweight; prefer calling `Chat()` or `Raw()` per use rather than caching service instances

2. **Services**: Lightweight wrappers that snapshot client options
   - `ChatService`: Chat completion endpoints (Complete, CompleteStream)
   - `RawService`: Raw HTTP request handling (Do, DoRaw, Stream, StreamReader)

3. **Per-Request Options**: Can override client defaults for single calls
   - Applied by value with defensive header copies
   - Contexts and HTTP clients are shared references

### Key Design Principles

From `doc.go` and README:

- **Immutability**: Clients are immutable; to change options, create a new Client
- **Concurrency**: Clients and services are safe for concurrent use by default
- **Feature Parity**: SDK favors upstream API feature parity; breaking changes possible as API evolves
- **DTOs**: Request/response types closely aligned to the HuggingFace API
- **Streaming**: Server-Sent Events (SSE) based streaming for chat completions
- **Error Handling**: Distinct APIError vs SDKError types with categorization
- **HTTP Client Injection**: Factory functions return fresh client values; avoid sharing mutable transports unless synchronized
- **Validation**: Enforced during JSON marshal/unmarshal time
- **Context Support**: Full context support for cancellation, timeouts, and request-scoped values
- **Defensive Copies**: Request options applied by value with defensive header copies; contexts and HTTP clients are shared

## Error Handling

### APIError
Represents errors returned by the HuggingFace API. Available at `github.com/Kardbord/hfgo/v4.APIError`.

**Fields**:
- `StatusCode`: HTTP status code
- `Message`: Human-readable error message
- `Body`: Response body as io.ReadCloser (caller responsible for closing)
- `RequestID`: X-Request-ID header value if available
- `Method`: HTTP method used
- `URL`: URL that was requested

**Helper Methods**:
- `IsClientError()`: Returns true for 4xx status codes
- `IsServerError()`: Returns true for 5xx status codes
- `IsAuthenticationError()`: Returns true for 401 Unauthorized
- `IsRateLimitError()`: Returns true for 429 Too Many Requests

**Type Assertion Pattern**:
```go
if apiErr, ok := err.(*hfgo.APIError); ok {
    if apiErr.IsAuthenticationError() {
        // Handle auth error
    }
}
```

### SDKError
Represents client-side SDK errors that occur before API response or during response unmarshaling.
Available at `github.com/Kardbord/hfgo/v4.SDKError`.

**Fields**:
- `Kind`: Error category (SDKErrorKind)
- `Message`: Human-readable error message
- `Err`: Underlying error (if any)

**Error Kinds**:
- `SDKErrorKindValidation`: Validation error in API responses
- `SDKErrorKindConfiguration`: Invalid or missing configuration
- `SDKErrorKindSerialization`: Serialization/deserialization error
- `SDKErrorKindTransport`: Transport-layer failure
- `SDKErrorKindInternal`: Internal SDK error

**Type Assertion Pattern**:
```go
if sdkErr, ok := err.(*hfgo.SDKError); ok {
    fmt.Printf("Kind %s: %s\n", sdkErr.Kind, sdkErr.Message)
}
```

## Configuration Options

All options are functions that return `hfgo.Option`. Applied to clients and per-request.

### Core Options
- `WithBaseURL(url string)`: Base URL for API requests (no query params/fragments)
- `WithToken(token string)`: Bearer authentication token
- `WithModel(model string)`: Model identifier for requests
- `WithProvider(provider string)`: Inference provider

### HTTP & Transport
- `WithHTTPClientFactory(factory func() http.Client)`: Factory for HTTP clients
  - Invoked when options are applied
  - Should return fresh client value
  - Avoid sharing mutable internals like Transport unless synchronized
  - Nil factory results in nil HTTP client
- `WithDefaultHTTPClient()`: Restores default HTTP client
- `WithUserAgentSuffix(suffix string)`: Appends suffix to SDK user agent

### Context & Timeouts
- `WithContext(ctx context.Context)`: Context for cancellation and timeouts
  - Nil context falls back to context.Background()

### Response Handling
- `WithMaxResponseBodyBytes(n int64)`: Max bytes read from response body
  - Values <= 0 fall back to default

### Headers
- `WithHeaders(h http.Header)`: Custom headers applied to all requests
  - Overrides existing values for matching keys
  - Per-request headers can still override
- `WithHeader(key, value string)`: Single header applied to all requests
- `WithDefaultHeader(key, value string)`: Header only if missing or empty

## Core Types

### ChatRequest
Represents a chat completion request. Key fields:

- `Model *string`: Model identifier (required)
- `Messages []ChatMessage`: Conversation history (required)
- `MaxTokens *int`: Max tokens in response (default: 1024, min: 0)
- `Temperature *float64`: Sampling temperature (0-2)
- `TopP *float64`: Nucleus sampling probability mass
- `TopLogProbs *int`: Number of most likely tokens (0-5, requires LogProbs=true)
- `FrequencyPenalty *float64`: Penalty for repeated tokens (-2.0 to 2.0)
- `PresencePenalty *float64`: Penalty for new topics (-2.0 to 2.0)
- `Stop []string`: Stop sequences (up to 4)
- `Seed *int64`: Deterministic sampling seed
- `Stream *bool`: Enable streaming (use CompleteStream, not Complete)
- `StreamOptions *ChatStreamOptions`: SSE stream configuration
- `Tools []ChatTool`: Available tools/functions
- `ToolChoice *ChatToolChoice`: Tool selection behavior (auto, none, required, or function spec)
- `ToolPrompt *string`: Prompt appended before tools
- `LogProbs *bool`: Return log probabilities
- `ResponseFormat *ChatResponseFormat`: Response format (text, json_schema, json_object, or provider-specific)

**Validation**:
- Enforced in `MarshalJSON()` method
- Invalid payloads surface as configuration errors
- Model can be set via request field or client option (request field takes precedence)

### ChatResponse
Response from non-streaming chat completion. Fields:

- `ID string`: Response identifier
- `Model string`: Model used
- `Choices []ChatChoice`: Generated choices
- `Usage`: Token usage statistics

Validation:
- Enforced in `UnmarshalJSON()` method
- Invalid response payloads surface as SDK validation errors

### ChatStream
Wraps streaming chat completion response from `CompleteStream()`.

**Methods**:
- `Recv(ctx context.Context) (ChatStreamResponse, error)`: Blocks until next chunk arrives
  - Returns `io.EOF` when stream ends
  - Merges tool call metadata across deltas
- `Close() error`: Releases underlying HTTP connection and decoder goroutine
  - Must be called to promptly release resources
  - Safe to call on nil stream

**Tool Call Metadata Merging**:
- Automatically caches and merges tool call ID, type, and function name across streaming deltas
- Ensures each delta includes complete tool call metadata

### ChatMessage
Represents a message in conversation history.

- `Role string`: Message role (system, user, assistant)
- `Content ChatMessageContent`: Message content
- `ToolCalls []ChatToolCall`: Tool calls made by assistant (if any)

### ChatTool
Represents a function tool available to the model.

- `Type string`: Tool type (currently "function")
- `Function ChatToolFunction`: Function definition

### ChatToolChoice
Controls tool selection behavior. Supports:
- String values: "auto", "none", "required"
- Object: `{"function":{"name":"..."}}`
- Provider-specific values

### ChatResponseFormat
Response format specification. Known types:
- "text": Plain text response
- "json_schema": JSON with schema validation
- "json_object": JSON object response
- Provider-specific values accepted

### ChatStreamOptions
Configuration for streaming responses.

- `IncludeUsage *bool`: Include token usage in stream

## Services

### ChatService

Created via `client.Chat()`. Methods:

#### Complete(req *ChatRequest, opts ...Option) (ChatResponse, error)
Non-streaming chat completion.

- Validates request is not nil
- Applies per-request options to override client defaults
- Rejects requests with Stream=true (use CompleteStream instead)
- Fills in model from request or client option
- Returns `ChatResponse` with all choices and usage stats
- Returns `SDKError` (kind: Configuration) for invalid requests

#### CompleteStream(req *ChatRequest, opts ...Option) (*ChatStream, error)
Streaming chat completion using SSE.

- Validates request is not nil
- Applies per-request options
- Automatically sets Stream=true in request
- Fills in model from request or client option
- Returns `*ChatStream` for consuming chunks
- Caller must call `Close()` on returned stream
- Returns `SDKError` (kind: Configuration) for invalid requests

**Streaming Pattern**:
```go
stream, err := client.Chat().CompleteStream(req)
if err != nil {
    // Handle error
}
defer stream.Close()

for {
    chunk, err := stream.Recv(ctx)
    if err != nil {
        if errors.Is(err, io.EOF) {
            break
        }
        // Handle error
    }
    // Process chunk
}
```

### RawService

Created via `client.Raw()`. For raw HTTP requests without type-safe JSON handling.

#### Do(body []byte, method, path string, opts ...Option) (*http.Response, error)
Raw request with error interpretation on non-2xx responses.

#### DoRaw(body []byte, method, path string, opts ...Option) (*http.Response, error)
Raw request without error interpretation (allows non-2xx responses).

#### DoReader(body io.Reader, method, path string, opts ...Option) (*http.Response, error)
Streaming body with error interpretation.

#### DoRawReader(body io.Reader, method, path string, opts ...Option) (*http.Response, error)
Streaming body without error interpretation.

#### Stream(body []byte, method, path string, opts ...Option) (*RawStream, error)
SSE stream with error interpretation.

#### StreamReader(body io.Reader, method, path string, opts ...Option) (*RawStream, error)
SSE stream from reader with error interpretation.

**Note**: Caller responsible for closing `resp.Body` or `RawStream`.

## Internal Architecture

### Packages

#### internal/request
Lower-level HTTP, SSE, and JSON utilities.

**Key Functions**:
- `DoJSON[TReq, TResp]`: Type-safe JSON request/response
- `DoJSONStream[TReq, TResp]`: Streaming JSON via SSE
- `Do`: Raw HTTP request with error interpretation
- `DoRaw`: Raw HTTP request without error interpretation
- `StreamRaw`: SSE stream handling
- `Options`: Configuration snapshot applied to requests

**Key Types**:
- `JSONStream[T]`: Consumes JSON SSE events
- `RawStream`: Consumes raw SSE events
- `Options`: Request configuration (immutable snapshot)

#### internal/hferrors
Error types and helper functions.

- `APIError`: API response errors with helper methods
- `SDKError`: Client-side errors with categorization
- `SDKErrorKind`: Error categories

#### internal/chatstream
Tool call metadata handling for streaming responses.

- `ToolCallAccumulator`: Caches and merges tool call metadata across stream deltas

#### internal/sdkversion
Version management.

- `Version`: Current SDK version (semver)
- `UserAgent()`: Returns "hfgo/<version> (Go)"

#### internal/testutils
Test utilities (excluded from linting/coverage requirements).

## Endpoints

### Chat Completions
- **Constant**: `EndpointChatCompletion = "/v1/chat/completions"`
- **Method**: POST
- **Service**: `ChatService.Complete()` or `ChatService.CompleteStream()`

## Quality Assurance

### Testing Strategy

1. **Unit Tests**
   - Run with: `go test ./...`
   - Run with coverage: `go test -cover ./...`
   - Coverage tracked with codecov
   - All non-integration test files excluded from specific linters

2. **Race Condition Detection**
   - Run with: `go test -race ./...`
   - Part of CI/CD pipeline
   - Ensures concurrency safety

3. **Integration Tests**
   - Marked with `//go:build integration` or run with `-tags=integration`
   - Run against live HuggingFace API
   - Require `HUGGING_FACE_TOKEN` secret
   - Retry logic (3 attempts with 10s delays) in CI
   - Automatic issue creation on failures in CI, closure on success

### Linting & Code Quality

**Tools**:
- `golangci-lint`: Comprehensive linting (custom config in .golangci.yml)
- `go vet`: Standard Go vetting
- `gofmt`, `gofumpt`, `goimports`, `golines`: Code formatters
- `CodeQL`: GitHub security analysis

**Key Linter Exclusions**:
- Test files: bodyclose, cyclop, errcheck, exhaustruct, gocognit, goconst, gocyclo, maintidx, varnamelen
- Examples: revive, mnd, exhaustruct, errcheck, godoclint
- Specific checks disabled: arangolint, depguard, err113, ginkgolinter, gocyclo, goheader, gomodguard, lll, noinlineerr, nonamedreturns, paralleltest, promlinter, recvcheck, testpackage, tparallel, whitespace, wrapcheck, wsl, zerologlint

**Linter Configuration Highlights**:
- All linters enabled by default except disabled list
- `exhaustruct` enforced for exported functions
- `funlen`: max 40 statements per function
- JSON tags: snake_case
- Generated files excluded strictly

### CI/CD Workflows

#### unit-tests.yml
- **Triggers**: Push to main/v4-draft, PRs, weekly schedule, manual
- **Platforms**: Windows, macOS, Linux
- **Linux Only**: Coverage upload to codecov

#### lint.yml
- **Triggers**: Push to main/v4-draft, PRs, weekly schedule, manual
- **Steps**: `go vet`, `golangci-lint`

#### integration-tests.yml
- **Triggers**: Push to main/v4-draft, PRs, weekly schedule, manual
- **Concurrency**: Serialized (prevent simultaneous runs)
- **Retry Logic**: 3 attempts with 10s delays
- **Failure Handling**: Auto-creates GitHub issue with label "integration-test-failure"
- **Success Handling**: Auto-closes related issues

#### codeql-analysis.yml, build.yml, release.yml, report-card.yml
- Standard GitHub Actions workflows

### Development Commands

**Build Script** (`tools/build.sh`):
```bash
./tools/build.sh
```

Runs in sequence:
1. `gofmt -s -w .`: Format code
2. `go mod tidy`: Tidy dependencies
3. `go vet ./...`: Vet code
4. `golangci-lint config verify && golangci-lint run --fix --disable godox ./...`: Lint with fixes
5. `go build ./...`: Build
6. `go test ./...`: Unit tests
7. `go test -race ./...`: Race detection
8. `go test -tags=integration ./...`: Integration tests
9. `go test -cover ./...`: Coverage report

**Manual Commands**:
```bash
# Format
gofmt -s -w .

# Tidy
go mod tidy

# Vet
go vet ./...

# Lint
golangci-lint run --fix ./...

# Build
go build ./...

# Unit tests
go test -timeout 600s -v ./...

# Race tests
go test -race -timeout 600s -v ./...

# Integration tests (requires HUGGING_FACE_TOKEN)
go test -tags=integration -timeout 600s -v ./...

# Coverage
go test -coverprofile=coverage.out -covermode=atomic ./...
```

## Examples

Located in `examples/` directory:

### examples/chat/basic/chat_basic.go
Demonstrates basic non-streaming chat completion.

**Key Points**:
- Create client with token and model
- Build ChatRequest with messages
- Call `client.Chat().Complete(request)`
- Access response fields (ID, Model, Choices, Usage)

### examples/chat/streaming/chat_streaming.go
Demonstrates streaming chat completion.

**Key Points**:
- Create client with token and model
- Build ChatRequest with messages
- Call `client.Chat().CompleteStream(request, WithContext(ctx))`
- Loop with `stream.Recv(ctx)` until `io.EOF`
- Call `stream.Close()` when done (or defer)

### examples/chat/convo/convo.go
Multi-message conversation example.

## Best Practices for Development

### 1. Concurrency & Immutability
- Create new Client for different configurations, don't mutate
- Services are lightweight; create per use rather than caching
- When injecting HTTP clients, return fresh values from factory
- Assume externally supplied objects (transports, etc.) aren't mutated unless you synchronize

### 2. Error Handling
- Always type-assert errors to APIError or SDKError
- Use helper methods on APIError (IsAuthenticationError, IsRateLimitError, etc.)
- Check SDKError.Kind for error categorization
- Always close Body on APIError to release resources

### 3. Request Handling
- Always validate requests before sending
- Use per-request options to override defaults for single calls
- Close streams (ChatStream, RawStream) to release resources
- Use context for cancellation and timeouts

### 4. Configuration
- Use immutable patterns; create new Client for different configs
- Options are applied by value (defensive copies)
- Per-request options override client defaults

### 5. Streaming
- Always call Close() on ChatStream or RawStream
- Prefer `defer stream.Close()` to ensure cleanup
- Handle `io.EOF` as normal stream completion
- Tool call metadata is automatically merged in ChatStream

### 6. Testing
- Use `go test -race ./...` to catch concurrency issues
- Mark integration tests with `//go:build integration` or run with `-tags=integration`
- Mark non-integration tests with //go:build !integration
- Ensure tests clean up resources (close bodies, streams, etc.)
- Integration tests require `HUGGING_FACE_TOKEN` environment variable

### 7. Code Quality
- Run `./tools/build.sh` before committing
- Ensure all linters pass (golangci-lint)
- Keep functions under 40 statements
- Document all functions, types, and fields with godoc comments

### 8. Value Receivers vs Pointer Receivers
- Use value receivers for immutable types (Client, ChatRequest, etc.)
- Use pointer receivers for mutable types (ChatStream, RawStream, etc.)
- Defend against mutations in value receivers (copy mutable state like headers)

### 9. Generics
- Leverage Go generics for type-safe request/response handling
- Use consistent naming: `TReq`, `TResp`, `T` for generic types

## Version Management

- **Current Version**: Defined in `internal/sdkversion.Version`
- **User Agent**: `hfgo/<version> (Go)`
- **Semantic Versioning**: Follows semver.org
- **Module Path**: `github.com/Kardbord/hfgo/v4`
- **Go Release Process**: Uses GoReleaser (.goreleaser.yml)

## Dependencies

**Required**:
- Go 1.25+

**Testing Only**:
- `github.com/stretchr/testify v1.9.0`: Assertion library

## Key Files & Structure

### Root Level Files
- `client.go`: Client type and creation
- `options.go`: Option functions
- `option.go`: Option type definition
- `chat_service.go`: ChatService implementation
- `chat_streaming.go`: ChatStream and streaming support
- `chat_request.go`, `chat_response.go`, `chat_common.go`: Chat types and validation
- `raw_service.go`: RawService implementation
- `errors.go`: Error type exports and re-exports
- `version.go`: Version management
- `doc.go`: Package-level documentation

### Test Files
- `*_test.go`: Unit tests (run with `go test ./...`)
- `*_integration_test.go`: Integration tests (run with `-tags=integration`)

### Internal Packages
- `internal/request/`: HTTP, SSE, and JSON utilities
- `internal/hferrors/`: Error types and definitions
- `internal/chatstream/`: Tool call accumulator for streaming
- `internal/sdkversion/`: Version management
- `internal/testutils/`: Test utilities

### Tools & Configuration
- `tools/build.sh`: Comprehensive build/test/lint script
- `.golangci.yml`: Linter configuration (comprehensive, all linters enabled by default)
- `.goreleaser.yml`: Release configuration
- `.github/workflows/`: CI/CD workflows (unit-tests, lint, integration-tests, codeql-analysis, build, release, report-card)
- `examples/`: Example code demonstrating SDK usage

## Important Notes

1. **Breaking Changes**: SDK follows upstream API; breaking changes possible as API evolves
2. **Streaming**: Always close streams to release HTTP connections and decoder goroutines
3. **Context Handling**: Nil contexts fall back to context.Background()
4. **HTTP Client Factory**: Must return fresh clients; avoid sharing mutable internals
5. **DTO Alignment**: Request/response types closely mirror HuggingFace API schema
6. **Validation Timing**: Enforced during JSON marshal/unmarshal, not during type construction
7. **Error Interpretation**: RawService has both error-interpreting (Do) and raw (DoRaw) paths
8. **Content Type Validation**: Invalid response content type surfaces as validation error; malformed JSON surfaces as serialization error
9. **Response Body Limits**: Enforced via `WithMaxResponseBodyBytes` option with default fallback

## Common Patterns

### Creating a Client
```go
client := hfgo.NewClient(
    hfgo.WithToken(token),
    hfgo.WithModel("model-name"),
)
```

### Non-Streaming Chat Completion
```go
response, err := client.Chat().Complete(&hfgo.ChatRequest{
    Messages: []hfgo.ChatMessage{
        {Role: "user", Content: hfgo.ChatMessageContent{Text: &prompt}},
    },
})
```

### Streaming Chat Completion
```go
stream, err := client.Chat().CompleteStream(req, hfgo.WithContext(ctx))
if err != nil {
    return err
}
defer stream.Close()

for {
    chunk, err := stream.Recv(ctx)
    if err != nil {
        if errors.Is(err, io.EOF) {
            break
        }
        return err
    }
    // Process chunk
}
```

### Error Handling
```go
response, err := client.Chat().Complete(req)
if err != nil {
    if apiErr, ok := err.(*hfgo.APIError); ok {
        if apiErr.IsRateLimitError() {
            // Handle rate limit
        }
    } else if sdkErr, ok := err.(*hfgo.SDKError); ok {
        // Handle SDK error
    }
}
```

### Per-Request Options
```go
response, err := client.Chat().Complete(
    req,
    hfgo.WithToken(overrideToken),
    hfgo.WithContext(ctx),
)
```

## Production Readiness Checklist

When working on this library, ensure:
- [ ] All tests pass (`go test ./...` and `go test -race ./...`)
- [ ] Linting passes (`golangci-lint run ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] Dependencies are tidy (`go mod tidy`)
- [ ] Documentation is updated (godoc comments)
- [ ] Examples work correctly
- [ ] Integration tests pass (with valid HF token)
- [ ] No new linter exclusions unless justified
- [ ] Backwards compatibility maintained (or clearly documented breaking changes)
- [ ] Error handling is comprehensive
- [ ] Streaming resources are properly released
- [ ] AGENTS.md is updated
- [ ] Code comments are accurate
