package revenium

import "context"

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	usageMetadataKey contextKey = "revenium_usage_metadata"
)

// WithUsageMetadata adds usage metadata to the context
func WithUsageMetadata(ctx context.Context, metadata map[string]interface{}) context.Context {
	return context.WithValue(ctx, usageMetadataKey, metadata)
}

// GetUsageMetadata retrieves usage metadata from the context
func GetUsageMetadata(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	metadata, ok := ctx.Value(usageMetadataKey).(map[string]interface{})
	if !ok {
		return nil
	}

	return metadata
}

// MergeMetadata merges two metadata maps, with priority to the second map
func MergeMetadata(base, override map[string]interface{}) map[string]interface{} {
	if base == nil && override == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Copy base metadata
	if base != nil {
		for k, v := range base {
			result[k] = v
		}
	}

	// Override with new metadata
	if override != nil {
		for k, v := range override {
			result[k] = v
		}
	}

	return result
}
