// E2E Test for Fal.ai Go Middleware
// Run with: FAL_API_KEY=xxx REVENIUM_METERING_API_KEY=xxx go test -v -run TestE2E
package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/revenium/revenium-middleware-fal-go/revenium"
)

func TestE2EImageGeneration(t *testing.T) {
	// Check required env vars
	falKey := os.Getenv("FAL_API_KEY")
	revKey := os.Getenv("REVENIUM_METERING_API_KEY")
	if falKey == "" || revKey == "" {
		t.Skip("Skipping E2E test: FAL_API_KEY and REVENIUM_METERING_API_KEY required")
	}

	// Set up environment for DEV
	os.Setenv("REVENIUM_METERING_BASE_URL", "https://api.dev.hcapp.io")
	os.Setenv("REVENIUM_DEBUG", "true")

	// Initialize middleware
	if err := revenium.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	// Create context with ALL metadata fields for comprehensive testing
	ctx := context.Background()
	metadata := map[string]interface{}{
		// Business context
		"organizationId": "org-e2e-test",
		"productId":      "product-fal-image-test",
		"subscriptionId": "sub-e2e-123",
		"taskType":       "e2e-image-test",
		"agent":          "e2e-test-agent",
		// Subscriber info
		"subscriber": map[string]interface{}{
			"id":    "user-e2e-test",
			"email": "e2e@test.revenium.io",
			"tier":  "enterprise",
		},
		// Tracing fields (ALL 10)
		"traceId":             "trace-e2e-" + time.Now().Format("20060102-150405"),
		"parentTransactionId": "parent-txn-123",
		"traceType":           "E2E_TEST",
		"traceName":           "Fal.ai Image E2E Test",
		"environment":         "dev",
		"region":              "us-west-2",
		"retryNumber":         0,
		"credentialAlias":     "fal-e2e-test",
		"taskId":              "task-e2e-" + time.Now().Format("150405"),
		// Multimodal job IDs
		"videoJobId": "vjob-e2e-test",
		"audioJobId": "ajob-e2e-test",
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create image generation request
	request := &revenium.FalRequest{
		Prompt:              "A simple test image of a red circle on white background",
		ImageSize:           "square",
		NumImages:           1,
		EnableSafetyChecker: true,
	}

	t.Log("Generating image with Fal.ai...")
	startTime := time.Now()

	// Generate image using flux/schnell (fastest/cheapest)
	// Model path without fal-ai/ prefix since client adds it
	resp, err := client.GenerateImage(ctx, "flux/schnell", request)
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	duration := time.Since(startTime)
	t.Logf("Image generated in %v", duration)
	t.Logf("Generated %d image(s)", len(resp.Images))

	for i, img := range resp.Images {
		urlPreview := img.URL
		if len(urlPreview) > 80 {
			urlPreview = urlPreview[:80] + "..."
		}
		t.Logf("  Image %d: %dx%d - %s", i+1, img.Width, img.Height, urlPreview)
	}

	// Wait for metering to complete (async)
	t.Log("Waiting for metering data to be sent...")
	time.Sleep(3 * time.Second)

	t.Log("SUCCESS: Image generated and metering data sent to Revenium DEV")
	t.Log("Check dashboard: https://app.dev.hcapp.io -> AI Transaction Log")
}

func TestE2EVideoGeneration(t *testing.T) {
	// Check required env vars
	falKey := os.Getenv("FAL_API_KEY")
	revKey := os.Getenv("REVENIUM_METERING_API_KEY")
	if falKey == "" || revKey == "" {
		t.Skip("Skipping E2E test: FAL_API_KEY and REVENIUM_METERING_API_KEY required")
	}

	// Set up environment for DEV
	os.Setenv("REVENIUM_METERING_BASE_URL", "https://api.dev.hcapp.io")
	os.Setenv("REVENIUM_DEBUG", "true")

	// Initialize middleware
	if err := revenium.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	// Create context with ALL metadata fields
	ctx := context.Background()
	metadata := map[string]interface{}{
		"organizationId":      "org-e2e-test",
		"productId":           "product-fal-video-test",
		"subscriptionId":      "sub-e2e-456",
		"taskType":            "e2e-video-test",
		"agent":               "e2e-test-agent",
		"traceId":             "trace-video-" + time.Now().Format("20060102-150405"),
		"parentTransactionId": "parent-video-txn",
		"traceType":           "E2E_VIDEO_TEST",
		"traceName":           "Fal.ai Video E2E Test",
		"environment":         "dev",
		"region":              "us-west-2",
		"retryNumber":         0,
		"credentialAlias":     "fal-e2e-video",
		"taskId":              "vtask-" + time.Now().Format("150405"),
		"videoJobId":          "fal-vjob-" + time.Now().Format("150405"),
		"audioJobId":          "",
		"subscriber": map[string]interface{}{
			"id":    "user-e2e-video",
			"email": "video-e2e@test.revenium.io",
		},
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create video generation request (using Kling for cheapest)
	request := &revenium.FalRequest{
		Prompt: "A simple animation of a bouncing ball",
	}

	t.Log("Generating video with Fal.ai (Kling)...")
	t.Log("NOTE: Video generation takes 1-2 minutes...")
	startTime := time.Now()

	// Generate video
	resp, err := client.GenerateVideo(ctx, "fal-ai/kling-video/v1/standard/text-to-video", request)
	if err != nil {
		t.Fatalf("Failed to generate video: %v", err)
	}

	duration := time.Since(startTime)
	t.Logf("Video generated in %v", duration)
	t.Logf("Video URL: %s", resp.Video.URL[:80]+"...")
	t.Logf("Duration: %.2f seconds", resp.Video.Duration)

	// Wait for metering
	time.Sleep(3 * time.Second)

	t.Log("SUCCESS: Video generated and metering data sent to Revenium DEV")
}

func main() {
	fmt.Println("Run with: go test -v -run TestE2E")
}
