// Deprecated: This module has moved to github.com/revenium/revenium-go-sdk/fal.
module github.com/revenium/revenium-middleware-fal-go

go 1.21

require github.com/joho/godotenv v1.5.1

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.4 // Deprecation tombstone: use revenium-go-sdk/fal instead.
	v1.0.3 // Consolidated into revenium-go-sdk/fal.
	v1.0.2 // Consolidated into revenium-go-sdk/fal.
	v1.0.1 // Consolidated into revenium-go-sdk/fal.
	v1.0.0 // Consolidated into revenium-go-sdk/fal.
	v0.0.1 // Consolidated into revenium-go-sdk/fal.
)
