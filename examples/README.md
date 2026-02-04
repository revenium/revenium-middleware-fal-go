# Examples - Revenium Middleware for Fal.ai (Go)

## Prerequisites

- Go 1.21+
- Fal.ai API key ([fal.ai](https://fal.ai/dashboard/keys))
- Revenium API key ([app.revenium.ai](https://app.revenium.ai))

## Setup

```bash
git clone https://github.com/revenium/revenium-middleware-fal-go.git
cd revenium-middleware-fal-go
go mod download
cp .env.example .env
# Edit .env with your API keys
```

## Examples

| Example | Description | Run |
|---------|-------------|-----|
| `basic/` | Image and video generation with Flux/Kling | `go run examples/basic/main.go` |
| `tracing/` | Distributed tracing, retry tracking, environment/region | `go run examples/tracing/main.go` |

## Metadata Fields

All metadata fields are optional. Add business context to your requests:

| Field | Type | Description |
|-------|------|-------------|
| `organizationName` | string | Human-readable organization name |
| `productName` | string | Human-readable product name |
| `taskType` | string | Task category (e.g., "image-generation", "video-generation") |
| `totalCost` | number | Cost override when provider pricing unavailable |
| `subscriber` | object | End-user identification (`id`, `email`) |
| `agent` | string | AI agent or workflow identifier |

### Trace Visualization Fields

| Field | Type | Description |
|-------|------|-------------|
| `traceId` | string | Link related requests together |
| `parentTransactionId` | string | Your ID for the parent operation (links child to parent) |
| `environment` | string | Deployment environment (e.g., "production") |
| `region` | string | Cloud region (e.g., "us-east-1") |
| `traceType` | string | Workflow category for grouping traces |
| `traceName` | string | Human-readable trace label |
| `credentialAlias` | string | API key identifier for auditing |
| `retryNumber` | int | Retry attempt (0 = first attempt) |
| `taskId` | string | Unique task identifier for job tracking |
| `videoJobId` | string | Video generation job ID (Fal.ai specific) |
| `audioJobId` | string | Audio generation job ID (Fal.ai specific) |

## Environment Variables

```bash
# Required
FAL_API_KEY=your_fal_api_key
REVENIUM_METERING_API_KEY=your_revenium_key

# Optional
REVENIUM_METERING_BASE_URL=https://api.revenium.ai
REVENIUM_CAPTURE_PROMPTS=true  # Enable prompt capture for analytics (default: false)
```

## Prompt Capture Feature

Enable prompt capture to track generation prompts in the Revenium dashboard:

```go
// Via environment variable
// REVENIUM_CAPTURE_PROMPTS=true

// Or programmatically
revenium.Initialize(
    revenium.WithCapturePrompts(true),
)
```

When enabled, the following fields are added to metering payloads:

| Field | Description |
|-------|-------------|
| `inputMessages` | JSON array: `[{"role": "user", "content": "<prompt>"}]` |
| `outputResponse` | Generated content URL(s) |
| `promptsTruncated` | `true` if prompt exceeded 50K characters |

## Supported Models

### Image Generation
- `flux/dev`, `flux-pro`, `flux/schnell`

### Video Generation
- `kling-video/v1/standard/text-to-video`
- `kling-video/v1/pro/text-to-video`

## Support

- [Documentation](https://docs.revenium.io)
- [Issues](https://github.com/revenium/revenium-middleware-fal-go/issues)
