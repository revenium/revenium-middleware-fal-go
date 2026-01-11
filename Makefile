.PHONY: help install test lint fmt clean build-examples
.PHONY: run-basic

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
