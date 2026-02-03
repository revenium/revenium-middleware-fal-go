# Revenium Middleware for Fal.ai (Go)

A lightweight, production-ready middleware that adds **Revenium metering and tracking** to Fal.ai API calls.

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![Documentation](https://img.shields.io/badge/docs-revenium.io-blue)](https://docs.revenium.io)
[![Website](https://img.shields.io/badge/website-revenium.ai-blue)](https://www.revenium.ai)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Seamless Integration** - Drop-in middleware with minimal code changes
- **Automatic Metering** - Tracks all API calls with detailed usage metrics
- **Image Generation** - Full support for Flux, SDXL, and other image models
- **Video Generation** - Support for Kling, Mochi, and other video models
- **Custom Metadata** - Add custom tracking metadata to any request
- **Production Ready** - Battle-tested and optimized for production use
- **Type Safe** - Built with Go's strong typing system

## Getting Started (5 minutes)

### Step 1: Create Your Project

```bash
mkdir my-fal-app
cd my-fal-app
go mod init my-fal-app
```

### Step 2: Install Dependencies

```bash
go get github.com/revenium/revenium-middleware-fal-go
go mod tidy
```

### Step 3: Create Environment File

Create `.env` file in your project root:

```bash
# Required - Get from https://fal.ai/dashboard/keys
FAL_API_KEY=your_fal_api_key_here

# Required - Get from Revenium dashboard (https://app.revenium.ai)
REVENIUM_METERING_API_KEY=your_revenium_api_key_here

# Optional - Revenium API base URL (defaults to production)
REVENIUM_METERING_BASE_URL=https://api.revenium.ai
```

**Replace the API keys with your actual keys!**

> **Automatic .env Loading**: The middleware automatically loads `.env` files from your project directory. No need to manually export environment variables!

## Examples

This repository includes runnable examples demonstrating how to use the Revenium middleware with Fal.ai:

- **[Examples Guide](./examples/README.md)** - Detailed guide for running examples
- **Go Examples**: `examples/basic/`

**Run examples after setup:**

```bash
# Clone this repository:
git clone https://github.com/revenium/revenium-middleware-fal-go.git
cd revenium-middleware-fal-go
go mod download
go mod tidy

# Run examples:
go run examples/basic/main.go
```

See **[Examples Guide](./examples/README.md)** for detailed setup instructions and what each example demonstrates.

## What Gets Tracked

The middleware automatically captures:

- **Image Count**: Number of images generated per request
- **Image Dimensions**: Width and height of generated images
- **Video Duration**: Length of generated videos in seconds
- **Request Duration**: Total time for each API call
- **Model Information**: Which Fal.ai model was used
- **Custom Metadata**: Business context you provide
- **Error Tracking**: Failed requests and error details

### Optional: Prompt Capture for Analytics

Enable prompt capture to track generation prompts in the Revenium dashboard:

```go
// Enable prompt capture (opt-in, default: false)
revenium.Initialize(
    revenium.WithCapturePrompts(true),
)
```

When enabled, the following fields are added to metering payloads:

| Field | Description |
|-------|-------------|
| `inputMessages` | JSON array with `[{"role": "user", "content": "<prompt>"}]` format |
| `outputResponse` | Generated content URL(s) |
| `promptsTruncated` | `true` if prompt exceeded 50,000 characters |

**Privacy Note**: Prompt capture is opt-in by default. Only enable if your use case requires prompt analytics.

### Custom Metadata

Add business context to requests using metadata:

```go
metadata := map[string]interface{}{
    "organizationName": "my-company",
    "productName":      "my-app",
    "totalCost":        0.05,  // Override cost when provider pricing unavailable
    "subscriber": map[string]interface{}{
        "id":    "user-123",
        "email": "user@example.com",
    },
}
ctx = revenium.WithUsageMetadata(ctx, metadata)
```

| Field | Type | Description |
|-------|------|-------------|
| `organizationName` | string | Human-readable organization name |
| `productName` | string | Human-readable product name |
| `totalCost` | number | Cost override (float64 or int accepted) |
| `subscriber` | object | End-user identification |

## Environment Variables

### Required

```bash
FAL_API_KEY=your_fal_api_key_here
REVENIUM_METERING_API_KEY=your_revenium_api_key_here
```

### Optional

```bash
# Fal.ai API base URL (defaults to production)
FAL_BASE_URL=https://fal.run

# Request timeout for Fal.ai API calls (default: 30m)
# Video generation (Kling) can take 2-10+ minutes
# Supports formats: "30m", "300s", "1800" (interpreted as seconds)
FAL_REQUEST_TIMEOUT=30m

# Revenium API base URL (defaults to production)
REVENIUM_METERING_BASE_URL=https://api.revenium.ai

# Default metadata for all requests (human-readable names preferred)
REVENIUM_ORGANIZATION_NAME=my-company
REVENIUM_PRODUCT_NAME=my-app

# Enable prompt capture for analytics (opt-in, default: false)
# When enabled, generation prompts and output URLs are sent to Revenium
REVENIUM_CAPTURE_PROMPTS=false

# Debug logging
REVENIUM_LOG_LEVEL=INFO
```

## Configuration Options

The middleware supports both environment variables and programmatic configuration:

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| Fal.ai API Key | `FAL_API_KEY` | (required) | Your Fal.ai API key |
| Fal.ai Base URL | `FAL_BASE_URL` | `https://fal.run` | Fal.ai API endpoint |
| Request Timeout | `FAL_REQUEST_TIMEOUT` | `30m` | HTTP request timeout |
| Revenium API Key | `REVENIUM_METERING_API_KEY` | (required) | Your Revenium API key |
| Revenium Base URL | `REVENIUM_METERING_BASE_URL` | `https://api.revenium.ai` | Revenium API endpoint |
| Organization Name | `REVENIUM_ORGANIZATION_NAME` | (optional) | Human-readable organization name (preferred) |
| Product Name | `REVENIUM_PRODUCT_NAME` | (optional) | Human-readable product name (preferred) |
| Capture Prompts | `REVENIUM_CAPTURE_PROMPTS` | `false` | Enable prompt analytics |
| Log Level | `REVENIUM_LOG_LEVEL` | `INFO` | Logging verbosity |
| Verbose Startup | `REVENIUM_VERBOSE_STARTUP` | `false` | Show startup details |

### Programmatic Configuration

```go
// Configure with functional options
err := revenium.Initialize(
    revenium.WithFalAPIKey("fal-api-key"),
    revenium.WithReveniumAPIKey("hak_your_key"),
    revenium.WithCapturePrompts(true),  // Enable prompt capture for analytics
    revenium.WithRequestTimeout(10 * time.Minute),
)
```

## Supported Models

### Image Generation

- `flux/dev` - Flux image generation
- `flux-pro` - Flux Pro (higher quality)
- `stable-diffusion-xl` - Stable Diffusion XL

### Video Generation

- `kling-video/v1/standard/text-to-video` - Kling video generation
- `mochi-v1` - Mochi video generation

## Troubleshooting

### Metering data not appearing in Revenium dashboard

**Problem**: Your app runs successfully but no data appears in Revenium.

**Solution**: The middleware sends metering data asynchronously in the background. If your program exits too quickly, the data won't be sent. Add a delay before exit:

```go
// At the end of your main() function
time.Sleep(2 * time.Second)
```

### "Failed to initialize" error

Check your API keys:

```bash
echo $FAL_API_KEY
echo $REVENIUM_METERING_API_KEY
```

### Enable debug logging

```bash
export REVENIUM_LOG_LEVEL=DEBUG
go run main.go
```

## Requirements

- **Go**: 1.21 or higher
- **Fal.ai API Key**: Get from [fal.ai/dashboard/keys](https://fal.ai/dashboard/keys)
- **Revenium API Key**: Get from [app.revenium.ai](https://app.revenium.ai)

## Documentation

For more information and advanced usage:

- [Revenium Documentation](https://docs.revenium.io)
- [Fal.ai Documentation](https://fal.ai/docs)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## Security

See [SECURITY.md](SECURITY.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, feature requests, or contributions:

- **GitHub Repository**: [revenium/revenium-middleware-fal-go](https://github.com/revenium/revenium-middleware-fal-go)
- **Issues**: [Report bugs or request features](https://github.com/revenium/revenium-middleware-fal-go/issues)
- **Documentation**: [docs.revenium.io](https://docs.revenium.io)
- **Contact**: Reach out to the Revenium team for additional support

---

**Built by Revenium**
