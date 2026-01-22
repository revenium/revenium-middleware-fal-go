// Package main demonstrates Scenario B: Different values for hard-coding detection.
//
// This example uses COMPLETELY DIFFERENT values from comprehensive/main.go
// to verify no values are accidentally hard-coded in the middleware.
//
// Compare the metering payloads from both scenarios - they should differ
// in every user-settable field including:
// - All metadata fields
// - Prompt content (for prompt capture validation)
// - inputMessages, outputResponse, promptsTruncated fields
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

// ScenarioBMetadata - INTENTIONALLY DIFFERENT from Scenario A
// Every user-settable field has a DIFFERENT value to detect hard-coding
func GetScenarioBMetadata() map[string]interface{} {
	return map[string]interface{}{
		// Business Context - ALL DIFFERENT FROM SCENARIO A
		"organizationId": "org-startup-xyz-staging",           // Different: was "org-acme-corporation-prod"
		"productId":      "prod-mvp-image-api-beta",           // Different: was "prod-creative-suite-enterprise-v3"
		"subscriptionId": "sub-freemium-trial-2026-jan",       // Different: was "sub-enterprise-annual-2026-q1"
		"taskType":       "prototype-asset-creation",          // Different: was "creative-asset-generation"
		"agent":          "dev-worker-local-macbook-01",       // Different: was "media-worker-us-west-2-node-07"

		// Distributed Tracing - ALL DIFFERENT
		"traceId":             "trace-bbb222ccc333-ddd444-eee555", // Different UUID pattern
		"parentTransactionId": "txn-local-debug-session-999",     // Different format
		"traceType":           "local-debug",                     // Different: was "distributed"
		"traceName":           "dev-iteration-rapid-prototype",   // Different pipeline name
		"taskId":              "task-debug-20260122-999999-zzz",  // Different task ID

		// Environment & Infrastructure - ALL DIFFERENT
		"environment":     "development",               // Different: was "production"
		"region":          "eu-west-1",                 // Different: was "us-west-2"
		"credentialAlias": "fal-dev-key-secondary",     // Different: was "fal-production-key-primary"
		"retryNumber":     2,                           // Different: was 0

		// Subscriber Information - COMPLETELY DIFFERENT structure
		"subscriber": map[string]interface{}{
			"id":          "usr-dev-tester-bob-12345",        // Different user
			"email":       "bob.developer@startup-xyz.io",    // Different email
			"name":        "Bob Developer",                   // Different name
			"plan":        "freemium-trial",                  // Different: was "enterprise-unlimited"
			"department":  "engineering",                     // Different: was "marketing"
			"costCenter":  "CC-ENG-RND-2026",                 // Different cost center
			"teamId":      "team-backend-api",                // Different team
			"accountTier": "free",                            // Different: was "platinum"
			"billingCode": "BC-FREE-TRIAL-001",               // Different billing code
			"projectCode": "PROJ-MVP-VALIDATION",             // Different project
			"customFields": map[string]interface{}{
				"experimentId": "EXP-API-PERF-TEST-001",
				"debugMode":    true,
				"testCategory": "performance-benchmark",
			},
		},

		// Multimodal Job Identifiers - DIFFERENT
		"videoJobId": "vjob-debug-test-99999-audio",       // Different
		"audioJobId": "ajob-debug-test-99999-narration",   // Different

		// Quality Metrics - DIFFERENT
		"responseQualityScore": 0.72, // Different: was 0.95
	}
}

func main() {
	fmt.Println("============================================================")
	fmt.Println("  Fal.ai Middleware - SCENARIO B: Different Values Test")
	fmt.Println("  Purpose: Detect hard-coded values by using different data")
	fmt.Println("  Prompt Capture Enabled (inputMessages, outputResponse)")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Compare these values against Scenario A (comprehensive/main.go)")
	fmt.Println("Every field should be DIFFERENT in the metering payload.")
	fmt.Println()

	// Verify required environment variables
	if os.Getenv("FAL_API_KEY") == "" {
		log.Fatal("FAL_API_KEY environment variable is required")
	}
	if os.Getenv("REVENIUM_METERING_API_KEY") == "" {
		log.Fatal("REVENIUM_METERING_API_KEY environment variable is required")
	}

	// Initialize with verbose logging and prompt capture enabled
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

	// Get Scenario B metadata (DIFFERENT from Scenario A)
	metadata := GetScenarioBMetadata()

	// Print the metadata being sent
	fmt.Println("=== SCENARIO B METADATA (should differ from A) ===")
	fmt.Println()
	prettyJSON, _ := json.MarshalIndent(metadata, "", "  ")
	fmt.Println(string(prettyJSON))
	fmt.Println()

	// Create context with Scenario B metadata
	ctx := context.Background()
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Run image generation with DIFFERENT prompt too
	// Note: When CapturePrompts is enabled, this prompt will be captured as inputMessages
	// Compare with Scenario A's prompt to verify prompt capture reflects actual content
	request := &revenium.FalRequest{
		Prompt:              "Abstract geometric art: vibrant neon shapes floating in dark space, minimalist composition, digital art style, 4k resolution",
		ImageSize:           "square", // Different: was "landscape_16_9"
		NumImages:           1,
		NumInferenceSteps:   20,   // Different: was 28
		GuidanceScale:       5.0,  // Different: was 7.5
		EnableSafetyChecker: true,
	}

	fmt.Println("Generating image with Flux model (Scenario B)...")
	startTime := time.Now()

	resp, err := client.GenerateImage(ctx, "fal-ai/flux/dev", request)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	elapsed := time.Since(startTime)

	// Display results
	fmt.Println()
	fmt.Println("=== SCENARIO B RESULTS ===")
	fmt.Printf("Generated %d image(s) in %v\n", len(resp.Images), elapsed)
	for i, img := range resp.Images {
		fmt.Printf("  Image %d: %s (%dx%d)\n", i+1, img.URL, img.Width, img.Height)
	}

	// Allow time for async metering
	fmt.Println()
	fmt.Println("Waiting for metering to complete...")
	time.Sleep(3 * time.Second)

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  SCENARIO B Complete")
	fmt.Println()
	fmt.Println("  VALIDATION STEPS:")
	fmt.Println("  1. Run comprehensive/main.go (Scenario A)")
	fmt.Println("  2. Run comprehensive-b/main.go (this, Scenario B)")
	fmt.Println("  3. Compare DEBUG logs for metering payloads")
	fmt.Println("  4. ALL user-settable fields should be DIFFERENT:")
	fmt.Println("     - All metadata fields (organizationId, traceId, etc.)")
	fmt.Println("     - inputMessages should contain DIFFERENT prompts")
	fmt.Println("     - outputResponse should have different image URLs")
	fmt.Println("============================================================")
}
