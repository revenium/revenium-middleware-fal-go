package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/revenium/revenium-middleware-fal-go/revenium"
)

// This example demonstrates trace visualization features for distributed tracing,
// retry tracking, and custom trace categorization with Fal.ai image/video generation.
//
// Features demonstrated:
// 1. Basic tracing with environment and region
// 2. Distributed tracing with parent-child relationships
// 3. Retry tracking for failed operations
// 4. Custom trace categorization and naming

func main() {
	fmt.Println("=== Revenium Middleware - Fal.ai Tracing Example ===")
	fmt.Println()

	// Initialize the middleware
	if err := revenium.Initialize(); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	// Get the client
	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Example 1: Basic tracing with environment and region
	fmt.Println("--- Example 1: Basic Tracing ---")
	if err := basicTracingExample(client); err != nil {
		log.Printf("Basic tracing error: %v", err)
	}

	fmt.Println()

	// Example 2: Distributed tracing (parent-child relationship)
	fmt.Println("--- Example 2: Distributed Tracing ---")
	if err := distributedTracingExample(client); err != nil {
		log.Printf("Distributed tracing error: %v", err)
	}

	fmt.Println()

	// Example 3: Retry tracking
	fmt.Println("--- Example 3: Retry Tracking ---")
	if err := retryTrackingExample(client); err != nil {
		log.Printf("Retry tracking error: %v", err)
	}

	// Wait for metering to complete
	time.Sleep(2 * time.Second)

	fmt.Println("\nTracing examples completed successfully!")
}

func basicTracingExample(client *revenium.ReveniumFal) error {
	// Create context with trace visualization metadata
	ctx := context.Background()
	metadata := map[string]interface{}{
		// Business context
		"organizationName": "my-company",
		"productName":      "media-service",
		"taskType":         "image-generation",
		"subscriber": map[string]interface{}{
			"id":    "user-123",
			"email": "user@example.com",
		},
		// Trace visualization fields
		"traceId":         fmt.Sprintf("trace-%d", time.Now().UnixMilli()),
		"environment":     "production",
		"region":          "us-east-1",
		"traceType":       "media-pipeline",
		"traceName":       "Product Image Generation",
		"credentialAlias": "fal-prod-key",
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Generate image
	request := &revenium.FalRequest{
		Prompt:    "A professional product photo of a smartphone on a white background",
		ImageSize: "square",
		NumImages: 1,
	}

	resp, err := client.GenerateImage(ctx, "fal-ai/flux/dev", request)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	fmt.Printf("Generated %d image(s) with tracing\n", len(resp.Images))
	fmt.Printf("  Trace ID: trace-%d\n", time.Now().UnixMilli())
	fmt.Printf("  Environment: production\n")
	fmt.Printf("  Region: us-east-1\n")

	return nil
}

func distributedTracingExample(client *revenium.ReveniumFal) error {
	// Parent transaction ID for linking related operations
	parentTxnID := fmt.Sprintf("parent-%d", time.Now().UnixMilli())

	// Step 1: Generate thumbnail (parent operation)
	fmt.Println("Step 1: Generating thumbnail (parent)...")
	ctx1 := context.Background()
	metadata1 := map[string]interface{}{
		"organizationName": "my-company",
		"productName":      "media-service",
		"taskType":         "thumbnail-generation",
		"traceId":          parentTxnID,
		"environment":      "production",
		"traceType":        "media-workflow",
		"traceName":        "E-commerce Image Pipeline",
	}
	ctx1 = revenium.WithUsageMetadata(ctx1, metadata1)

	request1 := &revenium.FalRequest{
		Prompt:    "A small thumbnail of a red sneaker",
		ImageSize: "square",
		NumImages: 1,
	}

	resp1, err := client.GenerateImage(ctx1, "fal-ai/flux/dev", request1)
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}
	fmt.Printf("  Parent transaction: %s\n", parentTxnID)
	fmt.Printf("  Generated thumbnail: %s\n", resp1.Images[0].URL[:50]+"...")

	// Step 2: Generate hero image (child operation linked to parent)
	fmt.Println("Step 2: Generating hero image (child)...")
	ctx2 := context.Background()
	metadata2 := map[string]interface{}{
		"organizationName":    "my-company",
		"productName":         "media-service",
		"taskType":            "hero-image-generation",
		"traceId":             fmt.Sprintf("child-%d", time.Now().UnixMilli()),
		"parentTransactionId": parentTxnID, // Links to parent
		"environment":         "production",
		"traceType":           "media-workflow",
		"traceName":           "E-commerce Image Pipeline",
	}
	ctx2 = revenium.WithUsageMetadata(ctx2, metadata2)

	request2 := &revenium.FalRequest{
		Prompt:    "A large hero image of a red sneaker with dramatic lighting",
		ImageSize: "landscape_16_9",
		NumImages: 1,
	}

	resp2, err := client.GenerateImage(ctx2, "fal-ai/flux/dev", request2)
	if err != nil {
		return fmt.Errorf("failed to generate hero image: %w", err)
	}
	fmt.Printf("  Child linked to parent: %s\n", parentTxnID)
	fmt.Printf("  Generated hero image: %s\n", resp2.Images[0].URL[:50]+"...")

	return nil
}

func retryTrackingExample(client *revenium.ReveniumFal) error {
	maxRetries := 3
	traceID := fmt.Sprintf("retry-trace-%d", time.Now().UnixMilli())

	for attempt := 0; attempt < maxRetries; attempt++ {
		fmt.Printf("Attempt %d/%d (retryNumber=%d)\n", attempt+1, maxRetries, attempt)

		ctx := context.Background()
		metadata := map[string]interface{}{
			"organizationName": "my-company",
			"productName":      "media-service",
			"taskType":         "image-generation-with-retry",
			"traceId":          traceID,
			"retryNumber":      attempt, // 0 for first attempt, 1+ for retries
			"environment":      "production",
			"traceName":        fmt.Sprintf("Image Generation Attempt %d", attempt+1),
		}
		ctx = revenium.WithUsageMetadata(ctx, metadata)

		request := &revenium.FalRequest{
			Prompt:    "A beautiful sunset over mountains",
			ImageSize: "landscape_4_3",
			NumImages: 1,
		}

		resp, err := client.GenerateImage(ctx, "fal-ai/flux/dev", request)
		if err != nil {
			fmt.Printf("  Error on attempt %d: %v\n", attempt+1, err)
			if attempt < maxRetries-1 {
				fmt.Println("  Retrying...")
				time.Sleep(1 * time.Second)
				continue
			}
			return fmt.Errorf("max retries reached: %w", err)
		}

		fmt.Printf("  Success on attempt %d\n", attempt+1)
		fmt.Printf("  Generated image: %s\n", resp.Images[0].URL[:50]+"...")
		return nil
	}

	return nil
}
