# Revenium Middleware - Fal.ai (Go)

Official Go middleware for integrating [Fal.ai](https://fal.ai) with [Revenium](https://revenium.ai) AI FinOps metering.

## Features

- Support for Fal.ai image generation models (Flux, SDXL, etc.)
- Support for Fal.ai video generation models (Kling, Mochi, etc.)
- Automatic metering to Revenium API
- Context-based metadata tracking
- Configurable via environment variables or code
- Thread-safe concurrent operations
- Comprehensive error handling
- Detailed logging support

## Installation

```bash
go get github.com/revenium/revenium-middleware-fal-go
```

## Quick Start

### 1. Set up environment variables

Create a `.env` file:

```bash
# Fal.ai Configuration
FAL_API_KEY=your-fal-api-key-here

# Revenium Configuration
REVENIUM_METERING_API_KEY=hak_your_api_key_here
```

### 2. Initialize the middleware

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/revenium/revenium-middleware-fal-go/revenium"
)

func main() {
    // Initialize middleware (loads from .env automatically)
    if err := revenium.Initialize(); err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }

    // Get the client
    client, err := revenium.GetClient()
    if err != nil {
        log.Fatalf("Failed to get client: %v", err)
    }

    // Create a request
    ctx := context.Background()
    request := &revenium.FalRequest{
        Prompt:    "A serene mountain landscape at sunset",
        ImageSize: "landscape_4_3",
        NumImages: 1,
    }

    // Generate image
    resp, err := client.GenerateImage(ctx, "flux/dev", request)
    if err != nil {
        log.Fatalf("Failed to generate image: %v", err)
    }

    fmt.Printf("Generated %d image(s)\n", len(resp.Images))
    for _, img := range resp.Images {
        fmt.Printf("Image URL: %s\n", img.URL)
    }
}
```

## Usage

### Image Generation

```go
// Create image generation request
request := &revenium.FalRequest{
    Prompt:         "A futuristic city at night",
    ImageSize:      "landscape_4_3",
    NumImages:      2,
    GuidanceScale:  7.5,
    NumInferenceSteps: 30,
    EnableSafetyChecker: true,
}

// Generate with Flux model
resp, err := client.GenerateImage(ctx, "flux/dev", request)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Access generated images
for _, img := range resp.Images {
    fmt.Printf("URL: %s, Size: %dx%d\n", img.URL, img.Width, img.Height)
}
```

### Video Generation

```go
// Create video generation request
request := &revenium.FalRequest{
    Prompt: "A drone flying over a futuristic city",
}

// Generate with Kling model
resp, err := client.GenerateVideo(ctx, "kling-video/v1/standard/text-to-video", request)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Access generated video
fmt.Printf("Video URL: %s\n", resp.Video.URL)
fmt.Printf("Duration: %.2f seconds\n", resp.Video.Duration)
```

### Adding Metadata for Tracking

```go
// Create context with business metadata
metadata := map[string]interface{}{
    "organizationId": "org-123",
    "productId":      "product-456",
    "subscriber": map[string]interface{}{
        "id":    "user-789",
        "email": "user@example.com",
    },
    "taskType": "image-generation",
}
ctx := revenium.WithUsageMetadata(context.Background(), metadata)

// Make request with metadata
resp, err := client.GenerateImage(ctx, "flux/dev", request)
```

## Supported Fal.ai Models

### Image Generation Models

- `flux/dev` - Flux image generation
- `flux-pro` - Flux Pro (higher quality)
- `stable-diffusion-xl` - Stable Diffusion XL

### Video Generation Models

- `kling-video/v1/standard/text-to-video` - Kling video generation
- `mochi-v1` - Mochi video generation

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `FAL_API_KEY` | Yes | - | Your Fal.ai API key |
| `FAL_BASE_URL` | No | `https://api.fal.ai` | Fal.ai API base URL |
| `REVENIUM_METERING_API_KEY` | Yes | - | Your Revenium API key (starts with `hak_`) |
| `REVENIUM_METERING_BASE_URL` | No | `https://api.revenium.ai` | Revenium API base URL |
| `REVENIUM_ORGANIZATION_ID` | No | - | Your organization ID |
| `REVENIUM_PRODUCT_ID` | No | - | Your product ID |
| `REVENIUM_LOG_LEVEL` | No | `INFO` | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) |

### Programmatic Configuration

```go
err := revenium.Initialize(
    revenium.WithFalAPIKey("your-fal-key"),
    revenium.WithReveniumAPIKey("hak_your_revenium_key"),
    revenium.WithReveniumOrgID("org-123"),
    revenium.WithReveniumProductID("product-456"),
)
```

## Metering

The middleware automatically sends metering data to Revenium for:

### Image Generation
- Image count
- Image dimensions (width/height)
- Model used
- Request duration

### Video Generation
- Video duration
- Video dimensions (width/height)
- Model used
- Request duration

Metering is done asynchronously in the background and won't block your API calls.

## Error Handling

```go
resp, err := client.GenerateImage(ctx, "flux/dev", request)
if err != nil {
    // Check error type
    if revenium.IsConfigError(err) {
        log.Println("Configuration error:", err)
    } else if revenium.IsMeteringError(err) {
        log.Println("Metering error:", err)
    } else {
        log.Println("Other error:", err)
    }
    return
}
```

## Logging

Set log level via environment:

```bash
export REVENIUM_LOG_LEVEL=DEBUG
```

Or programmatically:

```go
revenium.SetLogLevel(revenium.LogLevelDebug)
```

## Examples

See the [examples/](examples/) directory for complete working examples:

- `examples/basic/main.go` - Basic image and video generation

## Requirements

- Go 1.21 or higher
- Fal.ai API key ([Get one here](https://fal.ai))
- Revenium API key ([Sign up here](https://revenium.ai))

## Support

- Documentation: [https://docs.revenium.ai](https://docs.revenium.ai)
- Issues: [GitHub Issues](https://github.com/revenium/revenium-middleware-fal-go/issues)
- Email: support@revenium.io

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
