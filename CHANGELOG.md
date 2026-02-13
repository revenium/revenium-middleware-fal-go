# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3] - 2026-02-08

### Fixed
- BACK-452: Align model name and provider with LiteLLM naming convention
  - Model names now use full `fal_ai/fal-ai/{endpoint}` format matching backend's LiteLLM sync
  - Provider field changed from `fal` to `fal_ai` to match COALESCE pricing query
  - `normalizeModelName()` handles all input formats idempotently
  - Fixes multimodal pricing lookup failures caused by provider/model mismatch

### Changed
- BACK-456: Updated examples and documentation to use `organizationName` and `productName` as preferred fields
  - Human-readable names provide better UX in Revenium dashboard
  - Legacy `organizationId`/`productId` fields still supported for backwards compatibility

### Added
- Unit tests for `normalizeModelName()` covering all input format variations and idempotency

## [1.0.1] - 2026-01-22

### Added
- Opt-in prompt capture for analytics via `WithCapturePrompts(true)` option
  - `inputMessages`: JSON array with role/content format for generation prompts
  - `outputResponse`: Generated content URLs from Fal.ai
  - `promptsTruncated`: Flag when prompt exceeds 50K character limit
- Comprehensive examples demonstrating UsageMetadata fields

### Fixed
- Model name format now includes `fal-ai/` prefix for LiteLLM pricing compatibility
  - Fal.ai models like `flux/dev` are now correctly sent as `fal-ai/flux/dev`
  - Applies to both image and video generation endpoints

## [1.0.0] - 2026-01-09

### Added
- Initial release of Revenium Fal.ai Go middleware
- Support for Fal.ai image generation models (Flux, SDXL)
- Support for Fal.ai video generation models (Kling, Mochi)
- Automatic metering to Revenium API for images and videos
- Context-based metadata tracking
- Environment variable configuration
- Programmatic configuration via Options pattern
- Comprehensive error handling with typed errors
- Structured logging with configurable log levels
- Thread-safe concurrent operations
- Asynchronous metering (fire-and-forget)
- Complete examples and documentation

### Features
- Image generation with customizable parameters
- Video generation with model-specific endpoints
- Automatic metering of image count, dimensions, and duration
- Business context tracking (organizationId, productId, subscriber, etc.)
- Configurable timeouts and retry logic
- Validation of API keys and configuration
