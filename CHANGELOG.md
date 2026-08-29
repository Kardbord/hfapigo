# Changelog

## [4.0.0-rc1](https://github.com/Kardbord/hfgo/compare/v4.0.0-rc0...v4.0.0-rc1) (2026-08-29)


### ⚠ BREAKING CHANGES

* `AnswerTableQuestion` now returns a single `TableQuestionAnswer` (renamed from `TableQuestionAnswering`) instead of `[]TableQuestionAnswering`. The HuggingFace API returns a bare JSON object for table question answering, not an array — despite the upstream schema declaring an array response. This change aligns the SDK with the actual API behavior.
* flatten inference services into Client methods; by-value chat requests; add deep-clone methods to request DTOs; update concurrency guarantees ([#156](https://github.com/Kardbord/hfgo/issues/156))

### Features

* add summarization task support ([#155](https://github.com/Kardbord/hfgo/issues/155)) ([94cf386](https://github.com/Kardbord/hfgo/commit/94cf3861d0da28b05c9c59e16b0f5c5a6ea054b2))
* Add support for the Fill Mask infernce task ([#154](https://github.com/Kardbord/hfgo/issues/154)) ([025174f](https://github.com/Kardbord/hfgo/commit/025174f37a22be1486f9acd165bed28fa39b29d4))
* Add support for the question answering inference task ([#165](https://github.com/Kardbord/hfgo/issues/165)) ([818c08f](https://github.com/Kardbord/hfgo/commit/818c08f6995bf0ab2dfe0de497e8541a35af7e79))
* Add support for token classification ([#164](https://github.com/Kardbord/hfgo/issues/164)) ([001fb9b](https://github.com/Kardbord/hfgo/commit/001fb9bfe7a6de917b3ed19eafee7b1f8966679c))
* Add translation task support & fix Bruno request headers ([#157](https://github.com/Kardbord/hfgo/issues/157)) ([9f2e6e0](https://github.com/Kardbord/hfgo/commit/9f2e6e058fc9259673dbb2c18d7f13041f3fe205))
* Table question answering ([#166](https://github.com/Kardbord/hfgo/issues/166)) ([7f015aa](https://github.com/Kardbord/hfgo/commit/7f015aa24e205ff52db6abd858c286e255594dd0))


### Bug Fixes

* **deps:** bump github.com/stretchr/testify from 1.11.1 to 1.12.1 ([#160](https://github.com/Kardbord/hfgo/issues/160)) ([d2da104](https://github.com/Kardbord/hfgo/commit/d2da10478421c1dbc7df637faafbbe997c1679a1))
* **deps:** bump github.com/stretchr/testify from 1.9.0 to 1.11.1 ([#76](https://github.com/Kardbord/hfgo/issues/76)) ([64eee79](https://github.com/Kardbord/hfgo/commit/64eee79059e6f0a3c8cd8b287c9169fb3dbf9ec8))
* Make checkout ref conditional for push/schedule events ([#123](https://github.com/Kardbord/hfgo/issues/123)) ([45264a6](https://github.com/Kardbord/hfgo/commit/45264a6b0b18d82d43ad01afff7456d9c8a25ec3))


### Code Refactoring

* flatten inference services into Client methods; by-value chat requests; add deep-clone methods to request DTOs; update concurrency guarantees ([#156](https://github.com/Kardbord/hfgo/issues/156)) ([f8af23e](https://github.com/Kardbord/hfgo/commit/f8af23eae1d385d3b0a5bd0c00f2b90391065baf))
