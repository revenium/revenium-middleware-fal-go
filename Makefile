.PHONY: help install test lint fmt clean build-examples
.PHONY: run-basic
.PHONY: test-unit test-race coverage ci vet

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod tidy

test: ## Run tests
	go test -v ./...

lint: ## Run basic linters (vet + fmt)
	go vet ./...
	go fmt ./...

fmt: ## Format code
	go fmt ./...

clean: ## Clean build artifacts
	go clean
	rm -rf bin/

# Fal Examples (Image/Video Generation)
run-basic: ## Run basic image generation example
	go run examples/basic/main.go

build-examples: ## Build all examples
	@mkdir -p bin
	go build -o bin/basic examples/basic/main.go

# CI/CD Testing Targets
vet: ## Run go vet
	go vet ./...

test-unit: ## Run unit tests only (CI-safe, no API keys required)
	@echo "Running unit tests..."
	go test -v -short ./... -coverprofile=coverage.out

test-race: ## Run tests with race detector
	@echo "Running tests with race detection..."
	go test -race -short ./...

coverage: ## Generate test coverage report
	go test -v -short -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

ci: fmt vet lint test-unit coverage ## Run full CI pipeline
	@echo "CI pipeline completed successfully"
