package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/revenium/revenium-middleware-fal-go/revenium"
)

func main() {
	fmt.Println("=== Revenium Middleware - Fal.ai Basic Example ===")
	fmt.Println()

	// Initialize the middleware with optional prompt capture
	// Option 1: Basic initialization (uses environment variables)
	// if err := revenium.Initialize(); err != nil {
	//     log.Fatalf("Failed to initialize middleware: %v", err)
	// }

	// Option 2: Initialize with programmatic options
	// Enable prompt capture for analytics (opt-in, default: false)
	// When enabled, generation prompts are captured and sent with metering data
	// Fields added: inputMessages, outputResponse, promptsTruncated
	if err := revenium.Initialize(
		revenium.WithCapturePrompts(true), // Enable prompt capture for analytics
	); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	// Alternative: Set via environment variable
	// REVENIUM_CAPTURE_PROMPTS=true

	// Get the client
	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}

	// Example 1: Image generation with Flux
	fmt.Println("--- Example 1: Image Generation ---")
	if err := generateImageExample(client); err != nil {
		log.Printf("Image generation error: %v", err)
	}

	fmt.Println()

	// Example 2: Video generation with Kling
	fmt.Println("--- Example 2: Video Generation ---")
	if err := generateVideoExample(client); err != nil {
		log.Printf("Video generation error: %v", err)
	}

	// Wait for metering to complete
	time.Sleep(2 * time.Second)

	fmt.Println("\nExamples completed successfully!")
}

func generateImageExample(client *revenium.ReveniumFal) error {
	// Create context with custom metadata
	ctx := context.Background()
	metadata := map[string]interface{}{
		"organizationName": "my-company",
		"productName":      "image-gen-app",
		"totalCost":        0.05, // Optional: override cost when provider pricing unavailable
		"subscriber": map[string]interface{}{
			"id":    "user-123",
			"email": "user@example.com",
		},
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create image generation request
	// Note: When CapturePrompts is enabled, the Prompt field will be captured
	// and sent to Revenium as inputMessages in the format:
	// [{"role": "user", "content": "<prompt>"}]
	// The generated image URL(s) will be captured as outputResponse
	request := &revenium.FalRequest{
		Prompt:    "A serene landscape with mountains and a lake at sunset",
		ImageSize: "landscape_4_3",
		NumImages: 1,
		EnableSafetyChecker: true,
	}

	// Generate image using canonical Fal.ai endpoint ID format
	// IMPORTANT: Always use the full "fal-ai/" prefix (e.g., "fal-ai/flux/dev")
	// This matches Fal.ai's billing API endpoint_id for accurate metering correlation
	resp, err := client.GenerateImage(ctx, "fal-ai/flux/dev", request)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	// Display results
	fmt.Printf("Generated %d image(s)\n", len(resp.Images))
	for i, img := range resp.Images {
		fmt.Printf("  Image %d: %s (%dx%d)\n", i+1, img.URL, img.Width, img.Height)
	}
	if resp.TimeTaken > 0 {
		fmt.Printf("Time taken: %.2f seconds\n", resp.TimeTaken)
	}

	return nil
}

func generateVideoExample(client *revenium.ReveniumFal) error {
	// Create context with custom metadata
	ctx := context.Background()
	metadata := map[string]interface{}{
		"organizationName": "my-company",
		"productName":      "video-gen-app",
		"totalCost":        0.50, // Optional: override cost when provider pricing unavailable
		"subscriber": map[string]interface{}{
			"id":    "user-123",
			"email": "user@example.com",
		},
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create video generation request - use shortest duration (5s) for faster testing
	// Note: When CapturePrompts is enabled, the Prompt field will be captured
	// and sent to Revenium. The generated video URL will be captured as outputResponse
	// Prompts exceeding 50,000 characters will be truncated (promptsTruncated=true)
	request := &revenium.FalRequest{
		Prompt:   "A drone flying over a futuristic city at night",
		Duration: "5", // Shortest possible video (5 seconds)
	}

	// Generate video using canonical Fal.ai endpoint ID format
	// IMPORTANT: Always use the full "fal-ai/" prefix for billing correlation
	resp, err := client.GenerateVideo(ctx, "fal-ai/kling-video/v1/standard/text-to-video", request)
	if err != nil {
		return fmt.Errorf("failed to generate video: %w", err)
	}

	// Display results
	fmt.Printf("Generated video: %s\n", resp.Video.URL)
	if resp.Video.Duration > 0 {
		fmt.Printf("Duration: %.2f seconds\n", resp.Video.Duration)
	}
	if resp.Video.Width > 0 && resp.Video.Height > 0 {
		fmt.Printf("Resolution: %dx%d\n", resp.Video.Width, resp.Video.Height)
	}
	if resp.TimeTaken > 0 {
		fmt.Printf("Time taken: %.2f seconds\n", resp.TimeTaken)
	}

	return nil
}
