package tests

import (
	"context"
	"os"
	"testing"

	"github.com/revenium/revenium-middleware-fal-go/revenium"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	// Reset global state before test
	revenium.Reset()

	// Set required environment variables
	os.Setenv("FAL_API_KEY", "test-fal-key")
	os.Setenv("REVENIUM_METERING_API_KEY", "hak_test_key")
	defer func() {
		os.Unsetenv("FAL_API_KEY")
		os.Unsetenv("REVENIUM_METERING_API_KEY")
	}()

	err := revenium.Initialize()
	require.NoError(t, err)
	assert.True(t, revenium.IsInitialized())

	client, err := revenium.GetClient()
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestInitializeWithOptions(t *testing.T) {
	revenium.Reset()

	// Clear environment variables that could interfere
	os.Unsetenv("FAL_API_KEY")
	os.Unsetenv("REVENIUM_METERING_API_KEY")

	err := revenium.Initialize(
		revenium.WithFalAPIKey("test-fal-key"),
		revenium.WithReveniumAPIKey("hak_test_key"),
		revenium.WithReveniumOrgID("org-test"),
	)
	require.NoError(t, err)

	client, err := revenium.GetClient()
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "org-test", client.GetConfig().ReveniumOrgID)
}

func TestInitializeMissingAPIKey(t *testing.T) {
	revenium.Reset()

	// Clear environment variables
	os.Unsetenv("FAL_API_KEY")
	os.Unsetenv("REVENIUM_METERING_API_KEY")

	err := revenium.Initialize()
	assert.Error(t, err)
	assert.True(t, revenium.IsConfigError(err))
}

func TestNewReveniumFal(t *testing.T) {
	cfg := &revenium.Config{
		FalAPIKey:       "test-fal-key",
		FalBaseURL:      "https://api.fal.ai",
		ReveniumAPIKey:  "hak_test_key",
		ReveniumBaseURL: "https://api.revenium.ai",
	}

	client, err := revenium.NewReveniumFal(cfg)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, cfg.FalAPIKey, client.GetConfig().FalAPIKey)
}

func TestWithUsageMetadata(t *testing.T) {
	ctx := context.Background()
	metadata := map[string]interface{}{
		"organizationId": "org-123",
		"productId":      "prod-456",
	}

	ctx = revenium.WithUsageMetadata(ctx, metadata)
	retrieved := revenium.GetUsageMetadata(ctx)

	assert.NotNil(t, retrieved)
	assert.Equal(t, "org-123", retrieved["organizationId"])
	assert.Equal(t, "prod-456", retrieved["productId"])
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *revenium.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &revenium.Config{
				FalAPIKey:      "test-key",
				ReveniumAPIKey: "hak_test_key",
			},
			wantErr: false,
		},
		{
			name: "missing fal key",
			config: &revenium.Config{
				ReveniumAPIKey: "hak_test_key",
			},
			wantErr: true,
		},
		{
			name: "missing revenium key",
			config: &revenium.Config{
				FalAPIKey: "test-key",
			},
			wantErr: true,
		},
		{
			name: "invalid revenium key format",
			config: &revenium.Config{
				FalAPIKey:      "test-key",
				ReveniumAPIKey: "invalid_key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizeReveniumBaseURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://api.revenium.ai", "https://api.revenium.ai"},
		{"https://api.revenium.ai/", "https://api.revenium.ai"},
		{"https://api.revenium.ai/meter/v2", "https://api.revenium.ai"},
		{"https://api.revenium.ai/meter", "https://api.revenium.ai"},
		{"", "https://api.revenium.ai"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := revenium.NormalizeReveniumBaseURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeMetadata(t *testing.T) {
	base := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	override := map[string]interface{}{
		"key2": "new_value2",
		"key3": "value3",
	}

	result := revenium.MergeMetadata(base, override)

	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "new_value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
}
