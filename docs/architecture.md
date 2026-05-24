# hfgo SDK Architecture Documentation

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
   - `TextClassificationService`: Text classification endpoints (Classify, ClassifyBatch)
   - `ZeroShotTextClassificationService`: Zero-shot text classification endpoints (Classify, ClassifyBatch)
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

### Concurrency Safety

This SDK is **safe for concurrent use out of the box**. No explicit synchronization is required when using clients and services from multiple goroutines.

**Concurrency Guarantees**:
- **Clients**: Fully concurrent-safe as immutable value types
- **Services**: Concurrent-safe as lightweight snapshots of client options
- **Per-request calls**: Each call is independent and concurrent-safe
- **Shared HTTP clients**: If you inject an HTTP client via `WithHTTPClientFactory()`, ensure it's either thread-safe by design or properly synchronized externally

**How It Works**:
The SDK achieves concurrency safety through immutability and snapshots:
1. Clients never mutate their options after creation
2. Services capture a snapshot of client options when created
3. Per-request options are applied by value (defensive copies)
4. No shared mutable state between goroutines

**Example**:
```go
// Safe: Single client used by multiple goroutines
client := NewClient(WithToken(token), WithModel("mistral-7b"))

// Each goroutine can safely call methods
go func() {
    resp, err := client.Chat().Complete(req)
    // ...
}()

go func() {
    stream, err := client.Chat().CompleteStream(req)
    // ...
}()
```

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
### Option Precedence

When an option can be specified at multiple levels (client-level, request-level, or in request structures), the following precedence applies (highest to lowest):

1. **Request Structure Fields** (if applicable): Values set directly in request structures (e.g., `ChatRequest.Model`)
2. **Request-Level Options**: Options passed to individual method calls (e.g., `Complete(req, WithModel("..."))`)
3. **Client-Level Options**: Options set when creating the Client (e.g., `NewClient(WithModel("..."))`)

This precedence ensures that more specific (request-level) configurations always override more general (client-level) configurations.

**Example**:
```go
// Client-level Model: "default-model"
client := NewClient(WithModel("default-model"))

// Request-level override: "request-model"
response, err := client.Chat().Complete(
    &ChatRequest{Messages: msgs},
    WithModel("request-model"),
)
// Result: Uses "request-model"

// Request structure field: "structure-model"
response, err := client.Chat().Complete(
    &ChatRequest{
        Model: ptr("structure-model"),
        Messages: msgs,
    },
    WithModel("request-model"),
)
// Result: Uses "structure-model" (highest precedence)
```

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
- `Function ChatFunctionDefinition`: Function definition

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

### RawEvent
Represents a raw SSE event from RawService streaming methods.

- `Data []byte`: Event data payload
- `Event string`: Event type identifier
- `ID string`: Event ID
- `Retry *time.Duration`: Retry duration hint (if provided)

## Services

### ChatService
Created via `client.Chat()`. Methods:

#### Complete(req *ChatRequest, opts ...Option) (ChatResponse, error)
Non-streaming chat completion.

**Model and Provider Precedence**:
The Model field is resolved with the following precedence (highest to lowest):
1. ChatRequest.Model field (if non-nil and non-empty)
2. Per-request options Model override
3. Client-level Model option

The Provider field is applied as a fallback only if the resolved Model does not already contain a provider (indicated by ":" in the model string). If the Model is in the format "model:provider", the Provider option is ignored.

**Behavior**:
- Validates request is not nil
- Applies per-request options to override client defaults
- Rejects requests with Stream=true (use CompleteStream instead)
- Normalizes model and provider fields
- Returns `ChatResponse` with all choices and usage stats
- Returns `SDKError` (kind: Configuration) for invalid requests

#### CompleteStream(req *ChatRequest, opts ...Option) (*ChatStream, error)
Streaming chat completion using SSE.

**Behavior**:
- Validates request is not nil
- Applies per-request options
- Automatically sets Stream=true in request
- Normalizes model and provider fields
- Returns `*ChatStream` for consuming chunks
- Caller must call `Close()` on returned stream
- Returns `SDKError` (kind: Configuration) for invalid requests

### TextClassificationService
Created via `client.ClassifyText()`. For text classification tasks like sentiment analysis.

#### Classify(req TextClassificationRequest, opts ...Option) ([]TextClassification, error)
Single text classification.

**Behavior**:
- Validates request contains exactly one input
- Returns error if multiple inputs provided (use ClassifyBatch instead)
- Applies per-request options
- Returns flat array of classifications for the single input
- Automatically unwraps batch response to get single input result

#### ClassifyBatch(req TextClassificationRequest, opts ...Option) ([][]TextClassification, error)
Batch text classification for multiple inputs.

**API Response Format Normalization**:
The service handles a quirk in the HuggingFace API where the response format differs based on whether the `TopK` parameter is explicitly set:
- **When TopK is explicitly set**: Returns `[[classifications for input1], [classifications for input2], ...]` (per-input format)
- **When TopK is unset (nil)**: Returns `[[all classifications together]]` (flat format)

This inconsistency is handled transparently by the `normalizeTextClassificationResponse()` helper function.

### ZeroShotTextClassificationService
Created via `client.ZeroShotClassifyText()`. For zero-shot text classification tasks.

#### Classify(req ZeroShotTextClassificationRequest, opts ...Option) ([]ZeroShotTextClassification, error)
Single input zero-shot text classification.

**Behavior**:
- Validates that candidate labels are provided in parameters
- Returns error if candidate labels are missing or empty
- Applies per-request options
- Returns flat array of classifications for the single input, ordered by score (descending)

#### ClassifyBatch(req ZeroShotTextClassificationBatchRequest, opts ...Option) ([][]ZeroShotTextClassification, error)
Batch zero-shot text classification for multiple inputs.

**API Response Normalization**:
The HuggingFace API returns batched zero-shot results in a different format than single inputs. The service transparently normalizes responses via `normalizeZeroShotTextClassificationResponse()`.

### RawService
Created via `client.Raw()`. For raw HTTP requests without type-safe JSON handling.

#### Do(body []byte, method, path string, opts ...Option) (*http.Response, error)
Raw request with error interpretation on non-2xx responses.

#### DoRaw(body []byte, method, path string, opts ...Option) (*http.Response, error)
Raw request without error interpretation (allows non-2xx responses).

#### Stream(body []byte, method, path string, opts ...Option) (*RawStream, error)
SSE stream with error interpretation.

#### StreamRaw(body []byte, method, path string, opts ...Option) (*RawStream, error)
SSE stream without error interpretation (allows non-2xx responses).

## Endpoints

### Chat Completions
- **Constant**: `EndpointChatCompletion = "/v1/chat/completions"`
- **Method**: POST
- **Service**: `ChatService.Complete()` or `ChatService.CompleteStream()`

## Quality Assurance

### Testing Strategy

1. **Unit Tests**: Run with `go test ./...`
2. **Race Condition Detection**: Run with `go test -race ./...`
3. **Integration Tests**: Run with `-tags=integration`

### Linting & Code Quality
- `golangci-lint`: Comprehensive linting with custom config
- Test files excluded from specific linters (bodyclose, cyclop, errcheck, etc.)
- Examples excluded from revive, mnd, exhaustruct, errcheck, godoclint

## Development Commands

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

## Best Practices for Development

### 1. Concurrency & Immutability
- Create new Client for different configurations, don't mutate
- Services are lightweight; create per use rather than caching

### 2. Error Handling
- Always type-assert errors to APIError or SDKError
- Use helper methods on APIError
- Always close Body on APIError

### 3. Request Mutation Safety
- Do not mutate ChatRequest after passing it to Complete or CompleteStream
- For concurrent requests, create a new ChatRequest for each call

### 4. Streaming
- Always call Close() on ChatStream or RawStream
- Prefer `defer stream.Close()` to ensure cleanup

### 5. Value Receivers vs Pointer Receivers
- Use value receivers for immutable types (Client, ChatRequest, etc.)
- Use pointer receivers for mutable types (ChatStream, RawStream, etc.)

### 6. Generics
- Leverage Go generics for type-safe request/response handling

