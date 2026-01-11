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

	// Initialize the middleware
	if err := revenium.Initialize(); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

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
		"organizationId": "org-fal-example",
		"productId":      "product-image-gen",
		"subscriber": map[string]interface{}{
			"id":    "user-123",
			"email": "user@example.com",
		},
		"taskType": "image-generation",
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create image generation request
	request := &revenium.FalRequest{
		Prompt:    "A serene landscape with mountains and a lake at sunset",
		ImageSize: "landscape_4_3",
		NumImages: 1,
		EnableSafetyChecker: true,
	}

	// Generate image
	resp, err := client.GenerateImage(ctx, "flux/dev", request)
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
		"organizationId": "org-fal-example",
		"productId":      "product-video-gen",
		"subscriber": map[string]interface{}{
			"id":    "user-123",
			"email": "user@example.com",
		},
		"taskType": "video-generation",
	}
	ctx = revenium.WithUsageMetadata(ctx, metadata)

	// Create video generation request - use shortest duration (5s) for faster testing
	request := &revenium.FalRequest{
		Prompt:   "A drone flying over a futuristic city at night",
		Duration: "5", // Shortest possible video (5 seconds)
	}

	// Generate video
	resp, err := client.GenerateVideo(ctx, "kling-video/v1/standard/text-to-video", request)
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
