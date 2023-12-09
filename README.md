# hfapigo

[![Unit Tests](https://github.com/Kardbord/hfapigo/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/Kardbord/hfapigo/actions/workflows/unit-tests.yml)
[![CodeQL](https://github.com/Kardbord/hfapigo/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Kardbord/hfapigo/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kardbord/hfapigo)](https://goreportcard.com/report/github.com/Kardbord/hfapigo)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kardbord/hfapigo.svg)](https://pkg.go.dev/github.com/Kardbord/hfapigo)

(Golang) Go bindings for the [Hugging Face Inference API](https://api-inference.huggingface.co/docs/python/html/index.html).
Directly call any model available in the [Model Hub](https://huggingface.co/models).

An API key is required for authorized access. To get one, create a [Hugging Face](https://huggingface.co/) profile.

## Usage

See the [examples](./examples) directory.

- [Audio Classification](./examples/audio_classification/main.go)
- [Conversational](./examples/conversational/main.go)
- [Fill Mask](./examples/fill_mask/main.go)
- [Image Classification](./examples/image_classification/main.go)
- [Image Segmentation](./examples/image_segmentation/main.go)
- [Image-To-Text](./examples/image_to_text/main.go)
- [Object Detection](./examples/object_detection/main.go)
- [Question Answering](./examples/question_answering/main.go)
- [Sentence Similarity](./examples/sentence_similarity/main.go)
- [Speech Recognition](./examples/speech_recognition/main.go)
- [Summarization](./examples/summarization/main.go)
- [Table Question Answering](./examples/table_question_answering/main.go)
- [Text Classification](./examples/text_classification/main.go)
- [Text Generation](./examples/text_generation/main.go)
- [Text-To-Image](./examples/text_to_image/main.go)
- [Token Classification](./examples/token_classification/main.go)
- [Translation](./examples/translation/main.go)
- [Zero-shot Classification](./examples/zeroshot/main.go)

## Resources

- [Hugging Face](https://huggingface.co/)
  - [Model Hub](https://huggingface.co/models)
  - [Datasets](https://huggingface.co/datasets)
  - [Hugging Face Inference API](https://api-inference.huggingface.co/docs/python/html/index.html) (HF API)
  - [HF on GitHub](https://github.com/huggingface)
  - [Diffuser Docs](https://huggingface.co/docs/diffusers/using-diffusers/pipeline_overview)
  - Official [Python bindings](https://github.com/huggingface/hfapi) for the HF API
  - Official [JavaScript bindings](https://github.com/huggingface/huggingface.js) for the HF Inference API
