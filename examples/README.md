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

## Environment Variables

```bash
FAL_KEY=your_fal_api_key
REVENIUM_METERING_API_KEY=your_revenium_key
REVENIUM_METERING_BASE_URL=https://api.revenium.ai
```

## Supported Models

### Image Generation
- `flux/dev`, `flux-pro`, `flux/schnell`

### Video Generation
- `kling-video/v1/standard/text-to-video`
- `kling-video/v1/pro/text-to-video`

## Support

- [Documentation](https://docs.revenium.io)
- [Issues](https://github.com/revenium/revenium-middleware-fal-go/issues)
