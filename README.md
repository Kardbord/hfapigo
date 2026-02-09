# hfapigo

[![Unit Tests](https://github.com/Kardbord/hfapigo/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/Kardbord/hfapigo/actions/workflows/unit-tests.yml)
[![CodeQL](https://github.com/Kardbord/hfapigo/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Kardbord/hfapigo/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kardbord/hfapigo)](https://goreportcard.com/report/github.com/Kardbord/hfapigo)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kardbord/hfapigo.svg)](https://pkg.go.dev/github.com/Kardbord/hfapigo)

(Golang) Go bindings for the [Hugging Face Inference API](https://huggingface.co/docs/inference-providers/tasks/index).
Directly call any model available in the [Model Hub](https://huggingface.co/models).

An API key is required for authorized access. To get one, create a [Hugging Face](https://huggingface.co/) profile.

## Usage

See the [examples](./examples) directory.

> Coming soon!

### Design notes

- `Client` values are immutable; options are fixed at creation time to keep concurrency simple and request behavior predictable.
- Service values capture a snapshot of client options when created for deterministic behavior.
- Clients and services are safe for concurrent use by default or when configured with immutable or synchronized dependencies.
- Per-request options can override client defaults for a single call.
- `WithHTTPClientFactory` expects a fresh client value; avoid sharing mutable internals like transports unless synchronized to preserve safe concurrency.
- `WithDefaultHTTPClient` restores the default client, while a nil factory is treated as a configuration error.
- RawService exposes both error-interpreting and raw request paths (Do vs DoRaw).

## Resources

- [Hugging Face](https://huggingface.co/)
- [Inference API JSON Schemas](https://github.com/huggingface/huggingface.js/tree/main/packages/tasks/src/tasks)
- [Model Hub](https://huggingface.co/models)
- [Datasets](https://huggingface.co/datasets)
- [Hugging Face Inference API](https://huggingface.co/docs/inference-providers/tasks/index)
- [HF on GitHub](https://github.com/huggingface)
  - Official [Python bindings](https://github.com/huggingface/huggingface_hub)
  - Official [JavaScript bindings](https://github.com/huggingface/huggingface.js)
