# hfgo

[![Unit Tests](https://github.com/Kardbord/hfgo/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/Kardbord/hfgo/actions/workflows/unit-tests.yml)
[![CodeQL](https://github.com/Kardbord/hfgo/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Kardbord/hfgo/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kardbord/hfgo)](https://goreportcard.com/report/github.com/Kardbord/hfgo)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kardbord/hfgo.svg)](https://pkg.go.dev/github.com/Kardbord/hfgo)

(Golang) Go bindings for the [Hugging Face Inference API](https://huggingface.co/docs/inference-providers/tasks/index).
Directly call any model available in the [Model Hub](https://huggingface.co/models).

An API key is required for authorized access. To get one, create a [Hugging Face](https://huggingface.co/) profile.

## ⚠️ Notice

**v3** and earlier are significantly out of date from the upstream API.
This project is currently undergoing a major overhaul, at the end of
which **v4** will be released, the module will be renamed to `hfgo`,
and **v3 will be deprecated**. See [#72](https://github.com/Kardbord/hfapigo/issues/72)
for more information.

## Usage

See the [examples](./examples) directory.

> Coming soon!

### Design notes

- `Client` values are immutable; options are fixed at creation time to keep concurrency simple and request behavior predictable.
- Service values capture a snapshot of client options when created for deterministic behavior.
- Clients and services are safe for concurrent use by default or when configured with immutable or synchronized dependencies.
- Per-request options can override client defaults for a single call.
- The SDK favors upstream feature parity and uses DTOs closely aligned to the API; breaking changes are possible as the upstream API evolves.
- `WithHTTPClientFactory` expects a fresh client value; avoid sharing mutable internals like transports unless synchronized to preserve safe concurrency.
- `WithDefaultHTTPClient` restores the default client, while a nil factory is treated as a configuration error.
- RawService exposes both error-interpreting and raw request paths (Do vs DoRaw).
- DTO validation is enforced during JSON marshal/unmarshal. Invalid request payloads surface as configuration errors. For responses, invalid content type surfaces as validation errors, while malformed JSON surfaces as serialization errors.

## Resources

- [Hugging Face](https://huggingface.co/)
- [Inference API JSON Schemas](https://github.com/huggingface/huggingface.js/tree/main/packages/tasks/src/tasks)
- [Model Hub](https://huggingface.co/models)
- [Datasets](https://huggingface.co/datasets)
- [Hugging Face Inference API](https://huggingface.co/docs/inference-providers/tasks/index)
- [HF on GitHub](https://github.com/huggingface)
  - Official [Python bindings](https://github.com/huggingface/huggingface_hub)
  - Official [JavaScript bindings](https://github.com/huggingface/huggingface.js)

