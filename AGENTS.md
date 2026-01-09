# AGENTS.md

> Machine-readable instructions for AI agents. Human docs: [README.md](README.md)

## Project Context

**Type**: Go middleware for Fal.ai API
**Purpose**: Add Revenium metering to Fal.ai image/video generation
**Stack**: Go 1.21+
**Module**: `github.com/revenium/revenium-middleware-fal-go`

## Commands

```bash
# Build
go build ./...

# Test
go test ./...

# Run example
cd examples && go run getting_started.go
```

## Architecture

```
revenium/
├── client.go      # Fal.ai client wrapper
├── config.go      # Configuration and validation
├── context.go     # Context metadata handling
├── errors.go      # Error types
├── logger.go      # Logging utilities
├── metering.go    # Revenium metering (fire-and-forget)
├── middleware.go  # Core middleware logic
└── version.go     # Dynamic version detection
```

## Critical Constraints

1. **Dynamic versioning** - Use `GetMiddlewareSource()` for metering payloads (never hardcode)
2. **Billing fields at TOP LEVEL** - `actualImageCount`, `durationSeconds` NOT in attributes
3. **Fire-and-forget metering** - Use goroutines, never block main request
4. **Auth header** - Use `x-api-key` (not `Authorization: Bearer`)

## Environment Variables

```bash
FAL_API_KEY=...                    # Required: Fal.ai API key
REVENIUM_METERING_API_KEY=hak_...  # Required: Revenium key (starts with hak_)
REVENIUM_DEBUG=true                # Optional: Enable debug logging
```

## Metering Endpoints

| Type | Endpoint | Key Fields |
|------|----------|------------|
| Image | `/meter/v2/ai/images` | `actualImageCount`, `requestedImageCount` |
| Video | `/meter/v2/ai/video` | `durationSeconds` |

## Common Errors

| Error | Fix |
|-------|-----|
| `package not found` | `go mod tidy` |
| `metering not tracking` | Check `REVENIUM_METERING_API_KEY` |
| `Fal auth failed` | Check `FAL_API_KEY` |

## References

- [Revenium Docs](https://docs.revenium.io)
- [Fal.ai Docs](https://fal.ai/docs)
- [AGENTS.md Spec](https://agents.md/)
