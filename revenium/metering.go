package revenium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Package-level HTTP client with connection pooling for metering requests.
// This prevents creating a new client for each metering call, avoiding
// file descriptor exhaustion and TCP handshake overhead under high load.
var meteringHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true, // JSON is already small
	},
}

// MeteringClient handles communication with the Revenium metering API
type MeteringClient struct {
	config *Config
}

// NewMeteringClient creates a new metering client
func NewMeteringClient(config *Config) (*MeteringClient, error) {
	if config == nil {
		return nil, NewConfigError("config cannot be nil", nil)
	}

	return &MeteringClient{
		config: config,
	}, nil
}

// SendImageMetering sends image generation metering data to Revenium
func (mc *MeteringClient) SendImageMetering(payload *MeteringPayload) error {
	url := fmt.Sprintf("%s/meter/v2/ai/images", mc.config.ReveniumBaseURL)
	return mc.sendMetering(url, payload)
}

// SendVideoMetering sends video generation metering data to Revenium
func (mc *MeteringClient) SendVideoMetering(payload *MeteringPayload) error {
	url := fmt.Sprintf("%s/meter/v2/ai/video", mc.config.ReveniumBaseURL)
	return mc.sendMetering(url, payload)
}

// sendMetering sends metering data to the specified endpoint with retry logic
func (mc *MeteringClient) sendMetering(url string, payload *MeteringPayload) error {
	const maxRetries = 3
	const initialBackoff = 100 * time.Millisecond

	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		err := mc.sendMeteringRequest(url, payload)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on validation errors
		if IsValidationError(err) {
			return err
		}
	}

	return NewMeteringError(
		fmt.Sprintf("metering failed after %d retries", maxRetries),
		lastErr,
	)
}

// sendMeteringRequest sends a single metering request
func (mc *MeteringClient) sendMeteringRequest(url string, payload *MeteringPayload) error {
	// Marshal payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return NewMeteringError("failed to marshal metering payload", err)
	}

	logMeteringPayload(payload)
	Debug("Sending metering data to %s", url)

	// Create request with background context for fire-and-forget
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewNetworkError("failed to create metering request", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-key", mc.config.ReveniumAPIKey)
	req.Header.Set("User-Agent", "revenium-middleware-fal-go/1.0")

	// Send request using pooled client (avoids creating new client per instance)
	resp, err := meteringHTTPClient.Do(req)
	if err != nil {
		return NewNetworkError("metering request failed", err)
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	logResponse(resp.StatusCode, string(body))

	// Check status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return NewValidationError(
				fmt.Sprintf("metering API returned %d: %s", resp.StatusCode, string(body)),
				nil,
			)
		}
		return NewMeteringError(
			fmt.Sprintf("metering API error: %d", resp.StatusCode),
			fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		)
	}

	Info("Metering data sent successfully")
	return nil
}

// generateTransactionID generates a unique transaction ID
func generateTransactionID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano()%1000)
}

// buildImageMeteringPayload builds a metering payload for image generation
func buildImageMeteringPayload(
	model string,
	imageResp *FalImageResponse,
	metadata map[string]interface{},
	duration time.Duration,
	requestTime time.Time,
) *MeteringPayload {
	payload := &MeteringPayload{
		StopReason:       "END",
		CostType:         "AI",
		OperationType:    string(OperationTypeImage),
		Model:            model,
		Provider:         "fal",
		ModelSource:      "FAL",
		TransactionID:    generateTransactionID(),
		RequestTime:      requestTime,
		ResponseTime:     requestTime.Add(duration),
		RequestDuration:  duration.Milliseconds(),
		MiddlewareSource: GetMiddlewareSource(),
	}

	// Add image-specific billing fields (TOP LEVEL per API contract)
	if imageResp != nil {
		imageCount := len(imageResp.Images)
		payload.ActualImageCount = &imageCount
		payload.RequestedImageCount = &imageCount // Same as actual for Fal

		// Image dimensions go in attributes (metadata, not billing)
		if len(imageResp.Images) > 0 {
			payload.Attributes = map[string]interface{}{
				"width":  imageResp.Images[0].Width,
				"height": imageResp.Images[0].Height,
			}
		}
	}

	// Add metadata fields
	if metadata != nil {
		if orgID, ok := metadata["organizationId"].(string); ok {
			payload.OrganizationID = orgID
		}
		if productID, ok := metadata["productId"].(string); ok {
			payload.ProductID = productID
		}
		if taskType, ok := metadata["taskType"].(string); ok {
			payload.TaskType = taskType
		}
		if agent, ok := metadata["agent"].(string); ok {
			payload.Agent = agent
		}
		if subscriptionID, ok := metadata["subscriptionId"].(string); ok {
			payload.SubscriptionID = subscriptionID
		}
		if traceID, ok := metadata["traceId"].(string); ok {
			payload.TraceID = traceID
		}
		// Distributed tracing fields
		if parentTransactionID, ok := metadata["parentTransactionId"].(string); ok {
			payload.ParentTransactionID = parentTransactionID
		}
		if traceType, ok := metadata["traceType"].(string); ok {
			payload.TraceType = traceType
		}
		if traceName, ok := metadata["traceName"].(string); ok {
			payload.TraceName = traceName
		}
		if environment, ok := metadata["environment"].(string); ok {
			payload.Environment = environment
		}
		if region, ok := metadata["region"].(string); ok {
			payload.Region = region
		}
		if retryNumber, ok := metadata["retryNumber"].(int); ok {
			payload.RetryNumber = &retryNumber
		}
		if credentialAlias, ok := metadata["credentialAlias"].(string); ok {
			payload.CredentialAlias = credentialAlias
		}
		if subscriber, ok := metadata["subscriber"].(map[string]interface{}); ok {
			payload.Subscriber = subscriber
		}
		if taskID, ok := metadata["taskId"].(string); ok {
			payload.TaskID = taskID
		}
		if videoJobID, ok := metadata["videoJobId"].(string); ok {
			payload.VideoJobID = videoJobID
		}
		if audioJobID, ok := metadata["audioJobId"].(string); ok {
			payload.AudioJobID = audioJobID
		}
		if responseQualityScore, ok := metadata["responseQualityScore"].(float64); ok {
			payload.ResponseQualityScore = &responseQualityScore
		}
	}

	return payload
}

// buildVideoMeteringPayload builds a metering payload for video generation
func buildVideoMeteringPayload(
	model string,
	videoResp *FalVideoResponse,
	metadata map[string]interface{},
	duration time.Duration,
	requestTime time.Time,
) *MeteringPayload {
	payload := &MeteringPayload{
		StopReason:       "END",
		CostType:         "AI",
		OperationType:    string(OperationTypeVideo),
		Model:            model,
		Provider:         "fal",
		ModelSource:      "FAL",
		TransactionID:    generateTransactionID(),
		RequestTime:      requestTime,
		ResponseTime:     requestTime.Add(duration),
		RequestDuration:  duration.Milliseconds(),
		MiddlewareSource: GetMiddlewareSource(),
	}

	// Add video-specific billing fields (TOP LEVEL per API contract)
	if videoResp != nil && videoResp.Video.Duration > 0 {
		payload.DurationSeconds = &videoResp.Video.Duration
		// RequestedDurationSeconds is required for PER_SECOND billing
		// For Fal.ai, we use the actual duration as the requested duration
		payload.RequestedDurationSeconds = &videoResp.Video.Duration

		// Video dimensions go in attributes (metadata, not billing)
		attrs := make(map[string]interface{})
		if videoResp.Video.Width > 0 {
			attrs["width"] = videoResp.Video.Width
		}
		if videoResp.Video.Height > 0 {
			attrs["height"] = videoResp.Video.Height
		}
		if len(attrs) > 0 {
			payload.Attributes = attrs
		}
	}

	// Add metadata fields
	if metadata != nil {
		if orgID, ok := metadata["organizationId"].(string); ok {
			payload.OrganizationID = orgID
		}
		if productID, ok := metadata["productId"].(string); ok {
			payload.ProductID = productID
		}
		if taskType, ok := metadata["taskType"].(string); ok {
			payload.TaskType = taskType
		}
		if agent, ok := metadata["agent"].(string); ok {
			payload.Agent = agent
		}
		if subscriptionID, ok := metadata["subscriptionId"].(string); ok {
			payload.SubscriptionID = subscriptionID
		}
		if traceID, ok := metadata["traceId"].(string); ok {
			payload.TraceID = traceID
		}
		// Distributed tracing fields
		if parentTransactionID, ok := metadata["parentTransactionId"].(string); ok {
			payload.ParentTransactionID = parentTransactionID
		}
		if traceType, ok := metadata["traceType"].(string); ok {
			payload.TraceType = traceType
		}
		if traceName, ok := metadata["traceName"].(string); ok {
			payload.TraceName = traceName
		}
		if environment, ok := metadata["environment"].(string); ok {
			payload.Environment = environment
		}
		if region, ok := metadata["region"].(string); ok {
			payload.Region = region
		}
		if retryNumber, ok := metadata["retryNumber"].(int); ok {
			payload.RetryNumber = &retryNumber
		}
		if credentialAlias, ok := metadata["credentialAlias"].(string); ok {
			payload.CredentialAlias = credentialAlias
		}
		if subscriber, ok := metadata["subscriber"].(map[string]interface{}); ok {
			payload.Subscriber = subscriber
		}
		if taskID, ok := metadata["taskId"].(string); ok {
			payload.TaskID = taskID
		}
		if videoJobID, ok := metadata["videoJobId"].(string); ok {
			payload.VideoJobID = videoJobID
		}
		if audioJobID, ok := metadata["audioJobId"].(string); ok {
			payload.AudioJobID = audioJobID
		}
		if responseQualityScore, ok := metadata["responseQualityScore"].(float64); ok {
			payload.ResponseQualityScore = &responseQualityScore
		}
	}

	return payload
}
