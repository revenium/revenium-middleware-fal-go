// Package main demonstrates a comprehensive E2E test that populates ALL available
// metering fields with realistic enterprise values. This example serves to:
//
// 1. Verify all metering fields are transmitted correctly to Revenium
// 2. Document all available metadata options for enterprise customers
// 3. Provide a template for customers who need full tracing/observability
// 4. Demonstrate prompt capture functionality (inputMessages, outputResponse, promptsTruncated)
//
// Run with: FAL_API_KEY=... REVENIUM_METERING_API_KEY=hak_... go run main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/revenium/revenium-middleware-fal-go/revenium"
)

// ComprehensiveMetadata contains ALL possible metering fields with realistic enterprise values.
// This struct documents every field available in the Revenium metering API for Fal.ai middleware.
type ComprehensiveMetadata struct {
	// Business Context - Core identifiers for billing and attribution
	// BACK-456: Use organizationName and productName (human-readable) as preferred fields
	OrganizationName string                 // Human-readable organization name (preferred)
	ProductName      string                 // Human-readable product name (preferred)
	SubscriptionID string                 // Subscription plan identifier
	TaskType       string                 // Type of AI task (image-generation, video-generation, etc.)
	Agent          string                 // Agent/worker identifier for distributed systems

	// Distributed Tracing - Full observability support
	TraceID             string // Root trace identifier (correlates all related operations)
	ParentTransactionID string // Parent operation ID (for hierarchical tracing)
	TraceType           string // Type of trace (distributed, local, batch)
	TraceName           string // Human-readable name for the trace/pipeline
	TaskID              string // Specific task identifier within a workflow

	// Environment & Infrastructure
	Environment     string // Deployment environment (production, staging, development)
	Region          string // Cloud region or data center
	CredentialAlias string // Alias for the API credential used (for multi-key setups)
	RetryNumber     int    // Retry attempt number (0 = first attempt)

	// Subscriber Information - End-user attribution
	Subscriber map[string]interface{} // Rich subscriber data for per-user billing

	// Multimodal Job Identifiers - For complex pipelines
	VideoJobID string // Video generation job ID (for async video pipelines)
	AudioJobID string // Audio generation job ID (for audio synthesis pipelines)

	// Quality Metrics
	ResponseQualityScore float64 // Quality score for response (0.0-1.0)
}

// GetRealisticEnterpriseMetadata returns metadata populated with realistic enterprise values.
// These values represent what a real production customer might use.
func GetRealisticEnterpriseMetadata() *ComprehensiveMetadata {
	return &ComprehensiveMetadata{
		// Business Context - Human-readable names (BACK-456)
		OrganizationName: "Acme Corporation",
		ProductName:      "Creative Suite Enterprise",
		SubscriptionID: "sub-enterprise-annual-2026-q1",
		TaskType:       "creative-asset-generation",
		Agent:          "media-worker-us-west-2-node-07",

		// Distributed Tracing
		TraceID:             "trace-f47ac10b-58cc-4372-a567-0e02b2c3d479",
		ParentTransactionID: "txn-parent-batch-job-20260122-143052",
		TraceType:           "distributed",
		TraceName:           "creative-asset-generation-pipeline",
		TaskID:              "task-media-gen-20260122-143052-00042",

		// Environment & Infrastructure
		Environment:     "production",
		Region:          "us-west-2",
		CredentialAlias: "fal-production-key-primary",
		RetryNumber:     0, // First attempt

		// Subscriber Information - Rich end-user data
		Subscriber: map[string]interface{}{
			"id":           "usr-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"email":        "creative.team@acme-corporation.com",
			"name":         "ACME Creative Team",
			"plan":         "enterprise-unlimited",
			"department":   "marketing",
			"costCenter":   "CC-MARKETING-2026",
			"teamId":       "team-creative-north-america",
			"accountTier":  "platinum",
			"billingCode":  "BC-2026-CREATIVE-042",
			"projectCode":  "PROJ-SPRING-CAMPAIGN-2026",
			"customFields": map[string]interface{}{
				"campaignId":   "CAMP-2026-SPRING-LAUNCH",
				"assetType":    "hero-image",
				"approvalFlow": "marketing-review-required",
			},
		},

		// Multimodal Job Identifiers
		VideoJobID: "vjob-render-20260122-143052-preview",
		AudioJobID: "ajob-voiceover-20260122-143052-narration",

		// Quality Metrics
		ResponseQualityScore: 0.95, // 95% quality score
	}
}

// ToContextMetadata converts ComprehensiveMetadata to the map format expected by WithUsageMetadata
func (m *ComprehensiveMetadata) ToContextMetadata() map[string]interface{} {
	return map[string]interface{}{
		// Business Context - Human-readable names (BACK-456)
		"organizationName": m.OrganizationName,
		"productName":      m.ProductName,
		"subscriptionId": m.SubscriptionID,
		"taskType":       m.TaskType,
		"agent":          m.Agent,

		// Distributed Tracing
		"traceId":             m.TraceID,
		"parentTransactionId": m.ParentTransactionID,
		"traceType":           m.TraceType,
		"traceName":           m.TraceName,
		"taskId":              m.TaskID,

		// Environment & Infrastructure
		"environment":     m.Environment,
		"region":          m.Region,
		"credentialAlias": m.CredentialAlias,
		"retryNumber":     m.RetryNumber,

		// Subscriber Information
		"subscriber": m.Subscriber,

		// Multimodal Job Identifiers
		"videoJobId": m.VideoJobID,
		"audioJobId": m.AudioJobID,

		// Quality Metrics
		"responseQualityScore": m.ResponseQualityScore,
	}
}

func main() {
	fmt.Println("============================================================")
	fmt.Println("  Revenium Fal.ai Middleware - Comprehensive E2E Test")
	fmt.Println("  All Metering Fields with Realistic Enterprise Values")
	fmt.Println("  Prompt Capture Enabled (inputMessages, outputResponse)")
	fmt.Println("============================================================")
	fmt.Println()

	// Verify required environment variables
	if os.Getenv("FAL_API_KEY") == "" {
		log.Fatal("FAL_API_KEY environment variable is required")
	}
	if os.Getenv("REVENIUM_METERING_API_KEY") == "" {
		log.Fatal("REVENIUM_METERING_API_KEY environment variable is required")
	}

	// Initialize the middleware with verbose logging to see the full payload
	os.Setenv("REVENIUM_LOG_LEVEL", "DEBUG")
	os.Setenv("REVENIUM_VERBOSE_STARTUP", "true")

	// Initialize with prompt capture enabled for analytics
	// This captures inputMessages, outputResponse, and promptsTruncated in metering payload
	if err := revenium.Initialize(
		revenium.WithCapturePrompts(true), // Enable prompt capture for analytics
	); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Run comprehensive image generation test
	fmt.Println("=== Test 1: Comprehensive Image Generation ===")
	fmt.Println()
	if err := comprehensiveImageTest(client); err != nil {
		log.Printf("Image test error: %v", err)
	}

	fmt.Println()
	fmt.Println("=== Test 2: Comprehensive Video Generation ===")
	fmt.Println()
	if err := comprehensiveVideoTest(client); err != nil {
		log.Printf("Video test error: %v", err)
	}

	// Allow time for async metering to complete
	fmt.Println()
	fmt.Println("Waiting for metering to complete...")
	time.Sleep(3 * time.Second)

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  Comprehensive E2E Test Complete")
	fmt.Println("  Check Revenium dashboard to verify all fields populated")
	fmt.Println("============================================================")
}

// comprehensiveImageTest runs an image generation test with ALL metadata fields populated
func comprehensiveImageTest(client *revenium.ReveniumFal) error {
	// Get realistic enterprise metadata
	metadata := GetRealisticEnterpriseMetadata()

	// Customize for image generation context
	metadata.TaskType = "hero-image-generation"
	metadata.TraceName = "spring-campaign-hero-image-pipeline"
	metadata.TaskID = "task-hero-img-20260122-143052-00001"

	// Update subscriber custom fields for this specific task
	if subscriber, ok := metadata.Subscriber["customFields"].(map[string]interface{}); ok {
		subscriber["assetType"] = "hero-image-1920x1080"
		subscriber["outputFormat"] = "png"
	}

	// Print the metadata being sent
	fmt.Println("Metadata being sent to Revenium:")
	fmt.Println("--------------------------------")
	contextMetadata := metadata.ToContextMetadata()
	prettyJSON, _ := json.MarshalIndent(contextMetadata, "", "  ")
	fmt.Println(string(prettyJSON))
	fmt.Println()

	// Create context with comprehensive metadata
	ctx := context.Background()
	ctx = revenium.WithUsageMetadata(ctx, contextMetadata)

	// Create image generation request
	request := &revenium.FalRequest{
		Prompt:              "A professional corporate hero image: modern glass skyscraper at golden hour, dramatic lighting, photorealistic, 8k quality, suitable for business website header",
		ImageSize:           "landscape_16_9", // Common hero image aspect ratio
		NumImages:           1,
		NumInferenceSteps:   28,    // High quality setting
		GuidanceScale:       7.5,   // Balanced creativity
		EnableSafetyChecker: true,
	}

	fmt.Println("Generating image with Flux model...")
	startTime := time.Now()

	// Use canonical model name with fal-ai/ prefix
	resp, err := client.GenerateImage(ctx, "fal-ai/flux/dev", request)
	if err != nil {
		return fmt.Errorf("image generation failed: %w", err)
	}

	elapsed := time.Since(startTime)

	// Display results
	fmt.Println()
	fmt.Println("Image Generation Results:")
	fmt.Println("-------------------------")
	fmt.Printf("Generated %d image(s) in %v\n", len(resp.Images), elapsed)
	for i, img := range resp.Images {
		fmt.Printf("  Image %d:\n", i+1)
		fmt.Printf("    URL: %s\n", img.URL)
		fmt.Printf("    Dimensions: %dx%d\n", img.Width, img.Height)
		if img.ContentType != "" {
			fmt.Printf("    Content-Type: %s\n", img.ContentType)
		}
	}
	if resp.Seed != 0 {
		fmt.Printf("  Seed: %d\n", resp.Seed)
	}
	if resp.TimeTaken > 0 {
		fmt.Printf("  Fal.ai processing time: %.2fs\n", resp.TimeTaken)
	}

	return nil
}

// comprehensiveVideoTest runs a video generation test with ALL metadata fields populated
func comprehensiveVideoTest(client *revenium.ReveniumFal) error {
	// Get realistic enterprise metadata
	metadata := GetRealisticEnterpriseMetadata()

	// Customize for video generation context
	metadata.TaskType = "promotional-video-generation"
	metadata.TraceName = "spring-campaign-promo-video-pipeline"
	metadata.TaskID = "task-promo-vid-20260122-143052-00002"
	metadata.VideoJobID = "vjob-promo-20260122-143052-main"

	// Update subscriber custom fields for video task
	if subscriber, ok := metadata.Subscriber["customFields"].(map[string]interface{}); ok {
		subscriber["assetType"] = "promotional-video-5s"
		subscriber["outputFormat"] = "mp4"
		subscriber["targetPlatform"] = "social-media"
	}

	// Print the metadata being sent
	fmt.Println("Metadata being sent to Revenium:")
	fmt.Println("--------------------------------")
	contextMetadata := metadata.ToContextMetadata()
	prettyJSON, _ := json.MarshalIndent(contextMetadata, "", "  ")
	fmt.Println(string(prettyJSON))
	fmt.Println()

	// Create context with comprehensive metadata
	ctx := context.Background()
	ctx = revenium.WithUsageMetadata(ctx, contextMetadata)

	// Create video generation request (using shortest duration for cost efficiency)
	request := &revenium.FalRequest{
		Prompt:   "A dynamic 5-second promotional clip: modern corporate office with professionals collaborating, smooth camera movement, cinematic lighting, professional quality",
		Duration: "5", // 5 seconds - shortest option for testing
	}

	fmt.Println("Generating video with Kling model (5 seconds)...")
	fmt.Println("Note: Video generation typically takes 2-5 minutes")
	startTime := time.Now()

	// Use canonical model name with fal-ai/ prefix
	resp, err := client.GenerateVideo(ctx, "fal-ai/kling-video/v1/standard/text-to-video", request)
	if err != nil {
		return fmt.Errorf("video generation failed: %w", err)
	}

	elapsed := time.Since(startTime)

	// Display results
	fmt.Println()
	fmt.Println("Video Generation Results:")
	fmt.Println("-------------------------")
	fmt.Printf("Generated video in %v\n", elapsed)
	fmt.Printf("  URL: %s\n", resp.Video.URL)
	if resp.Video.Duration > 0 {
		fmt.Printf("  Duration: %.2f seconds\n", resp.Video.Duration)
	}
	if resp.Video.Width > 0 && resp.Video.Height > 0 {
		fmt.Printf("  Dimensions: %dx%d\n", resp.Video.Width, resp.Video.Height)
	}
	if resp.Video.ContentType != "" {
		fmt.Printf("  Content-Type: %s\n", resp.Video.ContentType)
	}
	if resp.TimeTaken > 0 {
		fmt.Printf("  Fal.ai processing time: %.2fs\n", resp.TimeTaken)
	}

	return nil
}
