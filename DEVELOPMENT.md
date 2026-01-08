# Development Guide - Revenium Middleware for Fal.ai (Go)

## Quick Start for Developers & QA

### Prerequisites

1. **Go 1.21+** installed
2. **Git** for version control
3. **API Keys** (for examples and integration tests)

---

## Setup Instructions

### 1. Clone and Setup

```bash
git clone <repository-url>
cd revenium-middleware-fal-go
go mod download
```

### 2. Environment Configuration

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your actual API keys
# Required:
FAL_API_KEY=your_fal_key_here
REVENIUM_METERING_API_KEY=hak_your_key_here
```

---

## Testing Commands (Copy & Paste Ready)

### Build & Verification

```bash
# Build project (must pass)
go build ./...

# Clean dependencies
go mod tidy

# Verify no compilation errors
go vet ./...

# Format code
go fmt ./...
```

### Unit Tests (No API Keys Required)

```bash
# Test middleware core
go test -v ./tests/middleware_test.go

# All unit tests at once
go test -v ./tests/...
```

### Examples Testing (Requires .env with API Keys)

```bash
# Basic example
go run examples/basic/main.go
```

### Integration Tests (Requires API Keys)

```bash
# Full integration test suite
go test -v ./tests/e2e/...
```

---

## Expected Results

### Unit Tests Should Show:

```
=== RUN   TestMiddlewareInitialization
--- PASS: TestMiddlewareInitialization (0.00s)
...
PASS
```

### Examples Should Show:

```
=== Revenium Fal.ai Middleware - Basic Example ===

Request ID: req_xxx
Status: completed
Images generated: 1

Basic example completed successfully!
```

### Build Should Show:

```
# No output = success
go build ./...
```

---

## Development Workflow

### Daily Development:

```bash
# 1. Pull latest changes
git pull

# 2. Build and verify
go build ./...
go mod tidy

# 3. Run unit tests
go test -v ./tests/...

# 4. Test your changes with examples
go run examples/basic/main.go
```

### Before Committing:

```bash
# 1. Format code
go fmt ./...

# 2. Verify code
go vet ./...

# 3. Run all unit tests
go test -v ./tests/...

# 4. Test examples
go run examples/basic/main.go

# 5. Build final verification
go build ./...
```

---

## QA Testing Checklist

### Environment Setup

- [ ] Go 1.21+ installed
- [ ] Repository cloned
- [ ] Dependencies downloaded (`go mod download`)
- [ ] `.env` file created with valid API keys

### Build Verification

- [ ] `go build ./...` - No errors
- [ ] `go mod tidy` - No changes needed
- [ ] `go vet ./...` - No warnings

### Unit Tests (No API Keys)

- [ ] Middleware tests pass
- [ ] All tests pass

### Examples (With API Keys)

- [ ] Basic example works

### Expected Behaviors

- [ ] Image generation completes
- [ ] Metadata included in payloads
- [ ] No manual `export` commands needed
- [ ] Fire-and-forget metering working
- [ ] Dynamic version detection working

---

## Troubleshooting

### Common Issues:

#### "API key not found"

```bash
# Check .env file exists and has correct format
cat .env
# Should show: FAL_API_KEY=...
```

#### "Tests failing"

```bash
# Check if .env is interfering with unit tests
# Unit tests should work without .env
mv .env .env.backup
go test -v ./tests/...
mv .env.backup .env
```

#### "Examples not working"

```bash
# Verify .env file is loaded
go run examples/basic/main.go
# Should show: "Fal.ai API key loaded"
```

---

## Project Structure

```
revenium-middleware-fal-go/
├── revenium/           # Core middleware code
├── examples/           # Working examples
│   └── basic/         # Basic example
├── tests/             # Unit tests
│   └── e2e/           # End-to-end tests
├── .env.example       # Environment template
├── .env               # Your API keys (create this)
├── go.mod             # Go dependencies
├── README.md          # User documentation
└── DEVELOPMENT.md     # This file
```

---

## Success Criteria

**The project is working correctly when:**

1. All unit tests pass
2. All examples run without errors
3. Image generation completes successfully
4. Metadata is included in metering payloads
5. No manual `export` commands needed
6. Build completes without errors
7. Fire-and-forget metering works asynchronously
8. Dynamic version detection returns correct version

**Ready for production when all items above are complete**
