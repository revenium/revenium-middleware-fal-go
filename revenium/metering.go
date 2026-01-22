package revenium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

// MaxPromptLength is the maximum length for captured prompts.
// Prompts exceeding this length will be truncated with "...[TRUNCATED]" suffix.
const MaxPromptLength = 50000

// formatPromptAsInputMessages formats a single prompt string as JSON inputMessages
// for compatibility with the Revenium dashboard's unified prompt view.
//
// Output format: [{"role": "user", "content": "<prompt>"}]
//
// This format matches the LLM middleware pattern, enabling the Revenium dashboard
// to display prompts consistently across all AI providers (text, image, video).
//
// Returns:
//   - JSON string: The formatted inputMessages JSON
//   - bool: true if the prompt was truncated (exceeded MaxPromptLength)
func formatPromptAsInputMessages(prompt string) (string, bool) {
	if prompt == "" {
		return "", false
	}

	truncated := false
	if len(prompt) > MaxPromptLength {
		prompt = prompt[:MaxPromptLength] + "...[TRUNCATED]"
		truncated = true
	}

	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}

	jsonBytes, err := json.Marshal(messages)
	if err != nil {
		Warn("Failed to serialize prompt as inputMessages: %v", err)
		return "", truncated
	}

	return string(jsonBytes), truncated
}

// normalizeModelName ensures the model name has the required "fal-ai/" prefix
// for Revenium's model naming convention. Fal.ai endpoint IDs use the format
// "fal-ai/flux/dev" which matches the billing API's endpoint_id field.
//
// IMPORTANT: Users should pass the canonical model name with the "fal-ai/" prefix.
// This function provides backward compatibility but will log a warning when
// normalization is applied, as it indicates incorrect usage.
func normalizeModelName(model string) string {
	const falPrefix = "fal-ai/"
	if strings.HasPrefix(model, falPrefix) {
		return model
	}
	// Log warning when normalization is needed - indicates user passed incorrect format
	Warn("Model name '%s' is missing 'fal-ai/' prefix. Use canonical format 'fal-ai/%s' for clarity. Auto-normalizing to '%s%s'",
		model, model, falPrefix, model)
	return falPrefix + model
}

// buildImageMeteringPayload builds a metering payload for image generation
func buildImageMeteringPayload(
	model string,
	imageResp *FalImageResponse,
	metadata map[string]interface{},
	duration time.Duration,
	requestTime time.Time,
	capturePrompts bool,
	prompt string,
	outputURLs []string,
) *MeteringPayload {
	payload := &MeteringPayload{
		StopReason:       "END",
		CostType:         "AI",
		OperationType:    string(OperationTypeImage),
		Model:            normalizeModelName(model),
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

	// Add prompt capture fields when enabled (opt-in)
	if capturePrompts && prompt != "" {
		inputMessages, truncated := formatPromptAsInputMessages(prompt)
		if inputMessages != "" {
			payload.InputMessages = inputMessages
		}
		if truncated {
			payload.PromptsTruncated = true
		}
		// Output response contains the generated image URL(s)
		if len(outputURLs) > 0 {
			outputJSON, err := json.Marshal(outputURLs)
			if err == nil {
				payload.OutputResponse = string(outputJSON)
			}
		}
		Debug("Prompt capture enabled: captured %d chars, output %d URLs", len(prompt), len(outputURLs))
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
	requestedDuration string,
	capturePrompts bool,
	prompt string,
	outputURL string,
) *MeteringPayload {
	payload := &MeteringPayload{
		StopReason:       "END",
		CostType:         "AI",
		OperationType:    string(OperationTypeVideo),
		Model:            normalizeModelName(model),
		Provider:         "fal",
		ModelSource:      "FAL",
		TransactionID:    generateTransactionID(),
		RequestTime:      requestTime,
		ResponseTime:     requestTime.Add(duration),
		RequestDuration:  duration.Milliseconds(),
		MiddlewareSource: GetMiddlewareSource(),
	}

	// Parse requested duration from request (e.g., "5" or "10" seconds)
	// This is REQUIRED for PER_SECOND billing even if Fal.ai doesn't return actual duration
	var reqDurSeconds float64
	if requestedDuration != "" {
		// Trim whitespace to handle edge cases like " 5 " or "10\n"
		trimmed := strings.TrimSpace(requestedDuration)
		if parsed, err := strconv.ParseFloat(trimmed, 64); err == nil {
			reqDurSeconds = parsed
		} else {
			Warn("Failed to parse requestedDuration '%s': %v - video billing may fail with 422", requestedDuration, err)
		}
	}

	// Add video-specific billing fields (TOP LEVEL per API contract)
	// RequestedDurationSeconds = what user asked for (from request)
	// DurationSeconds = what was actually produced (from response, or fallback to requested)

	// Set RequestedDurationSeconds from user's request (semantic: what they asked for)
	if reqDurSeconds > 0 {
		payload.RequestedDurationSeconds = &reqDurSeconds
	}

	// Set DurationSeconds from actual response, or fallback to requested
	if videoResp != nil && videoResp.Video.Duration > 0 {
		payload.DurationSeconds = &videoResp.Video.Duration
		// If user didn't specify duration, use actual as fallback for requested
		if payload.RequestedDurationSeconds == nil {
			payload.RequestedDurationSeconds = &videoResp.Video.Duration
		}
	} else if reqDurSeconds > 0 {
		// Fal.ai didn't return duration - use requested as best estimate for actual
		payload.DurationSeconds = &reqDurSeconds
	}

	// Video dimensions go in attributes (metadata, not billing)
	if videoResp != nil {
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

	// Add prompt capture fields when enabled (opt-in)
	if capturePrompts && prompt != "" {
		inputMessages, truncated := formatPromptAsInputMessages(prompt)
		if inputMessages != "" {
			payload.InputMessages = inputMessages
		}
		if truncated {
			payload.PromptsTruncated = true
		}
		// Output response contains the generated video URL
		if outputURL != "" {
			payload.OutputResponse = outputURL
		}
		Debug("Prompt capture enabled: captured %d chars, output URL: %s", len(prompt), outputURL)
	}

	return payload
}
