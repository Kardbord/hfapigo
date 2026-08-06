# hfgo - HuggingFace Inference API Bruno Collection

A collection of [Bruno](https://www.usebruno.com/) requests for exercising the
HuggingFace Inference API (the `router.huggingface.co` endpoint used by the
hfgo Go SDK).

## Getting started

1. Open the Bruno **desktop app** and open this `hfgo` folder as a collection.
2. Select the `default` environment.
3. Set the environment variables you need (see below) and run any request.

## Environment variables

Defined in `environments/`:

- `HF_TOKEN` - your HuggingFace access token, sent as `Authorization: Bearer`.
- `HF_BASE_URL` - the inference host (`https://router.huggingface.co`).
- `model_*` - the model used per task (one per task folder).
- `sample_png` / `sample_audio` - file paths used by image/audio tasks.

Requests use single, batch, and parametrized variants where the API supports
them. The image/audio tasks upload a sample file via multipart form data.

## Structure

- top-level folders - one per inference task
- `environments/` - environment configurations
- `*.bru` - Bruno request files

## CLI (optional)

Requests can also be run with the Bruno CLI:

```
bruno-cli run --env default --output junit
```
