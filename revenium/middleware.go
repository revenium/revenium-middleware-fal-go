package revenium

import (
	"context"
	"sync"
	"time"
)

// ReveniumFal is the main middleware client that wraps Fal.ai API calls with metering
type ReveniumFal struct {
	config         *Config
	falClient      *FalClient
	meteringClient *MeteringClient
	mu             sync.RWMutex
	wg             sync.WaitGroup
}

var (
	globalClient *ReveniumFal
	globalMu     sync.RWMutex
	initialized  bool
)

// Initialize sets up the global Revenium middleware with configuration
func Initialize(opts ...Option) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if initialized {
		return nil
	}

	// Initialize logger first
	InitializeLogger()
	Info("Initializing Revenium Fal.ai middleware...")

	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Load from environment if not provided
	if err := cfg.loadFromEnv(); err != nil {
		Warn("Failed to load configuration from environment: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Create clients
	falClient, err := NewFalClient(cfg)
	if err != nil {
		return err
	}

	meteringClient, err := NewMeteringClient(cfg)
	if err != nil {
		return err
	}

	globalClient = &ReveniumFal{
		config:         cfg,
		falClient:      falClient,
		meteringClient: meteringClient,
	}

	initialized = true
	Info("Revenium Fal.ai middleware initialized successfully")
	return nil
}

// IsInitialized checks if the middleware is properly initialized
func IsInitialized() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return initialized
}

// GetClient returns the global Revenium client
func GetClient() (*ReveniumFal, error) {
	globalMu.RLock()
	defer globalMu.RUnlock()

	if !initialized {
		return nil, NewConfigError("middleware not initialized, call Initialize() first", nil)
	}

	return globalClient, nil
}

// NewReveniumFal creates a new Revenium client with explicit configuration
func NewReveniumFal(cfg *Config) (*ReveniumFal, error) {
	if cfg == nil {
		return nil, NewConfigError("config cannot be nil", nil)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	falClient, err := NewFalClient(cfg)
	if err != nil {
		return nil, err
	}

	meteringClient, err := NewMeteringClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ReveniumFal{
		config:         cfg,
		falClient:      falClient,
		meteringClient: meteringClient,
	}, nil
}

// GetConfig returns the configuration
func (r *ReveniumFal) GetConfig() *Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.config
}

// GenerateImage generates images using Fal.ai with automatic metering
func (r *ReveniumFal) GenerateImage(ctx context.Context, model string, request *FalRequest) (*FalImageResponse, error) {
	// Extract metadata from context
	metadata := GetUsageMetadata(ctx)

	// Record start time
	startTime := time.Now()

	// Call Fal.ai API
	resp, err := r.falClient.GenerateImage(ctx, model, request)
	if err != nil {
		return nil, err
	}

	// Calculate duration
	duration := time.Since(startTime)

	// Send metering data asynchronously (fire-and-forget)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.sendImageMetering(resp, model, metadata, duration, startTime)
	}()

	return resp, nil
}

// GenerateVideo generates a video using Fal.ai with automatic metering
func (r *ReveniumFal) GenerateVideo(ctx context.Context, model string, request *FalRequest) (*FalVideoResponse, error) {
	// Extract metadata from context
	metadata := GetUsageMetadata(ctx)

	// Record start time
	startTime := time.Now()

	// Call Fal.ai API
	resp, err := r.falClient.GenerateVideo(ctx, model, request)
	if err != nil {
		return nil, err
	}

	// Calculate duration
	duration := time.Since(startTime)

	// Capture the requested duration before the goroutine
	// Guard against nil request for defensive programming
	var requestedDuration string
	if request != nil {
		requestedDuration = request.Duration
	}

	// Send metering data asynchronously (fire-and-forget)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.sendVideoMetering(resp, model, metadata, duration, startTime, requestedDuration)
	}()

	return resp, nil
}

// sendImageMetering sends image metering data in the background
func (r *ReveniumFal) sendImageMetering(resp *FalImageResponse, model string, metadata map[string]interface{}, duration time.Duration, startTime time.Time) {
	defer func() {
		if rec := recover(); rec != nil {
			Error("Metering goroutine panic: %v", rec)
		}
	}()

	payload := buildImageMeteringPayload(model, resp, metadata, duration, startTime)

	if err := r.meteringClient.SendImageMetering(payload); err != nil {
		Error("Failed to send image metering data: %v", err)
	}
}

// sendVideoMetering sends video metering data in the background
func (r *ReveniumFal) sendVideoMetering(resp *FalVideoResponse, model string, metadata map[string]interface{}, duration time.Duration, startTime time.Time, requestedDuration string) {
	defer func() {
		if rec := recover(); rec != nil {
			Error("Metering goroutine panic: %v", rec)
		}
	}()

	payload := buildVideoMeteringPayload(model, resp, metadata, duration, startTime, requestedDuration)

	if err := r.meteringClient.SendVideoMetering(payload); err != nil {
		Error("Failed to send video metering data: %v", err)
	}
}

// Flush waits for all pending metering goroutines to complete.
// Call this before application shutdown to ensure all metering data is sent.
func (r *ReveniumFal) Flush() {
	r.wg.Wait()
}

// Close closes the client and cleans up resources.
// It calls Flush() to ensure all pending metering operations complete.
func (r *ReveniumFal) Close() error {
	r.Flush()
	r.mu.Lock()
	defer r.mu.Unlock()
	return nil
}

// Reset resets the global middleware state (for testing)
func Reset() {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalClient != nil {
		globalClient.Close()
		globalClient = nil
	}

	initialized = false
}
