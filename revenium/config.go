package revenium

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the Revenium middleware
type Config struct {
	// Fal.ai API configuration
	FalAPIKey      string
	FalBaseURL     string
	RequestTimeout time.Duration // HTTP request timeout (default: 1800s / 30 min for video generation)

	// Revenium metering configuration
	ReveniumAPIKey    string
	ReveniumBaseURL   string
	ReveniumOrgID     string
	ReveniumProductID string

	// Logging configuration
	LogLevel       string
	VerboseStartup bool
}

// Option is a functional option for configuring Config
type Option func(*Config)

// WithFalAPIKey sets the Fal.ai API key
func WithFalAPIKey(key string) Option {
	return func(c *Config) {
		c.FalAPIKey = key
	}
}

// WithRequestTimeout sets the HTTP request timeout for Fal.ai API calls
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.RequestTimeout = timeout
	}
}

// WithReveniumAPIKey sets the Revenium API key
func WithReveniumAPIKey(key string) Option {
	return func(c *Config) {
		c.ReveniumAPIKey = key
	}
}

// WithReveniumBaseURL sets the Revenium base URL
func WithReveniumBaseURL(url string) Option {
	return func(c *Config) {
		c.ReveniumBaseURL = url
	}
}

// WithReveniumOrgID sets the Revenium organization ID
func WithReveniumOrgID(id string) Option {
	return func(c *Config) {
		c.ReveniumOrgID = id
	}
}

// WithReveniumProductID sets the Revenium product ID
func WithReveniumProductID(id string) Option {
	return func(c *Config) {
		c.ReveniumProductID = id
	}
}

// loadFromEnv loads configuration from environment variables and .env files
// Only loads values that are not already set programmatically
func (c *Config) loadFromEnv() error {
	// First, try to load .env files automatically
	c.loadEnvFiles()

	// Then load from environment variables (only if not already set)
	if c.FalAPIKey == "" {
		c.FalAPIKey = os.Getenv("FAL_API_KEY")
	}
	if c.FalBaseURL == "" {
		c.FalBaseURL = getEnvOrDefault("FAL_BASE_URL", "https://fal.run")
	}
	if c.RequestTimeout == 0 {
		c.RequestTimeout = parseDurationFromEnv("FAL_REQUEST_TIMEOUT", 1800*time.Second) // 30 min for video generation
	}

	if c.ReveniumAPIKey == "" {
		c.ReveniumAPIKey = os.Getenv("REVENIUM_METERING_API_KEY")
	}
	if c.ReveniumBaseURL == "" {
		baseURL := getEnvOrDefault("REVENIUM_METERING_BASE_URL", "https://api.revenium.ai")
		c.ReveniumBaseURL = NormalizeReveniumBaseURL(baseURL)
	}
	if c.ReveniumOrgID == "" {
		c.ReveniumOrgID = os.Getenv("REVENIUM_ORGANIZATION_ID")
	}
	if c.ReveniumProductID == "" {
		c.ReveniumProductID = os.Getenv("REVENIUM_PRODUCT_ID")
	}

	if c.LogLevel == "" {
		c.LogLevel = getEnvOrDefault("REVENIUM_LOG_LEVEL", "INFO")
	}
	if !c.VerboseStartup {
		c.VerboseStartup = os.Getenv("REVENIUM_VERBOSE_STARTUP") == "true" || os.Getenv("REVENIUM_VERBOSE_STARTUP") == "1"
	}

	// Initialize logger early
	InitializeLogger()

	return nil
}

// loadEnvFiles loads environment variables from .env files
func (c *Config) loadEnvFiles() {
	envFiles := []string{
		".env.local", // Local overrides (highest priority)
		".env",       // Main env file
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Try current directory and parent directories
	searchDirs := []string{
		cwd,
		filepath.Dir(cwd),
		filepath.Join(cwd, ".."),
	}

	for _, dir := range searchDirs {
		for _, envFile := range envFiles {
			envPath := filepath.Join(dir, envFile)
			if _, err := os.Stat(envPath); err == nil {
				godotenv.Load(envPath)
			}
		}
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.FalAPIKey == "" {
		return NewConfigError("FAL_API_KEY is required", nil)
	}

	if c.ReveniumAPIKey == "" {
		return NewConfigError("REVENIUM_METERING_API_KEY is required", nil)
	}

	if !isValidReveniumAPIKey(c.ReveniumAPIKey) {
		return NewConfigError("invalid Revenium API key format (must start with 'hak_')", nil)
	}

	return nil
}

// isValidReveniumAPIKey checks if the API key has a valid format
func isValidReveniumAPIKey(key string) bool {
	if len(key) < 4 {
		return false
	}
	return key[:4] == "hak_"
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NormalizeReveniumBaseURL normalizes the base URL to a consistent format
func NormalizeReveniumBaseURL(baseURL string) string {
	if baseURL == "" {
		return "https://api.revenium.ai"
	}

	// Remove trailing slash if present
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	// Remove legacy endpoints if present
	if len(baseURL) >= 9 && baseURL[len(baseURL)-9:] == "/meter/v2" {
		return baseURL[:len(baseURL)-9]
	}

	if len(baseURL) >= 6 && baseURL[len(baseURL)-6:] == "/meter" {
		return baseURL[:len(baseURL)-6]
	}

	return baseURL
}

// parseDurationFromEnv parses a duration from an environment variable.
// Supports formats: "300s", "5m", "2h", or just "300" (interpreted as seconds).
// Returns defaultValue if the environment variable is not set or invalid.
func parseDurationFromEnv(envKey string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultValue
	}

	// First try standard Go duration format (e.g., "300s", "5m", "2h")
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}

	// Try parsing as plain seconds (e.g., "300")
	value = strings.TrimSpace(value)
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// If all parsing fails, return the default
	return defaultValue
}
