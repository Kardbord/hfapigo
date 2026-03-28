// Package hfgo provides Go bindings for the Hugging Face Inference API.
//
// Design notes:
//   - Clients are immutable; options are fixed at creation time, and services capture a snapshot.
//   - Clients and services are safe for concurrent use by default or when configured with immutable or synchronized dependencies.
//   - Per-request options can override client defaults for a single call.
//   - Request options are applied by value with defensive header copies; contexts and HTTP clients are shared.
//   - HTTP client injection uses a value factory; return a fresh client value to avoid shared state.
//   - The SDK favors upstream feature parity and uses DTOs closely aligned to the API; breaking changes are possible as the upstream API evolves.
//   - WithDefaultHTTPClient restores the default client; a nil factory is treated as a configuration error.
//   - RawService exposes both error-interpreting and raw request paths (Do vs DoRaw).
//   - DTO validation is enforced during JSON marshal/unmarshal. Invalid request
//     payloads surface as configuration errors. For responses, invalid content
//     type surfaces as validation errors, while malformed JSON surfaces as
//     serialization errors.
//   - Concurrency assumes externally supplied objects (for example, transports) are not mutated after use
//     unless callers provide their own synchronization.
package hfgo

// TODO: Wire up provider options. Currently provider is unused.
