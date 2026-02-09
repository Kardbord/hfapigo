// Package hfapigo provides Go bindings for the Hugging Face Inference API.
//
// Design notes:
// - Clients are immutable; options are fixed at creation time, and services capture a snapshot.
// - Per-request options can override client defaults for a single call.
// - HTTP client injection uses a value factory; return a fresh client value to avoid shared state.
// - WithDefaultHTTPClient restores the default client; a nil factory is treated as a configuration error.
// - Concurrency assumes externally supplied objects (for example, transports) are not mutated after use
//   unless callers provide their own synchronization.
package hfapigo
