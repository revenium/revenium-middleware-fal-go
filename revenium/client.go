package revenium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FalClient handles communication with the Fal.ai API
type FalClient struct {
	config     *Config
	httpClient *http.Client
}

// NewFalClient creates a new Fal.ai client
func NewFalClient(config *Config) (*FalClient, error) {
	if config == nil {
		return nil, NewConfigError("config cannot be nil", nil)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &FalClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.RequestTimeout, // Configurable via FAL_REQUEST_TIMEOUT (default: 30 min)
		},
	}, nil
}

// GenerateImage generates images using a Fal.ai model
func (c *FalClient) GenerateImage(ctx context.Context, model string, request *FalRequest) (*FalImageResponse, error) {
	endpoint := fmt.Sprintf("%s/fal-ai/%s", c.config.FalBaseURL, model)

	// Marshal request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, NewProviderError("failed to marshal request", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewNetworkError("failed to create request", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.config.FalAPIKey))

	logRequest("POST", endpoint, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Key [REDACTED]",
	})

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("request failed", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewNetworkError("failed to read response", err)
	}

	logResponse(resp.StatusCode, string(body))

	// Check for errors
	if resp.StatusCode >= 400 {
		var falErr FalError
		if err := json.Unmarshal(body, &falErr); err == nil {
			falErr.Status = resp.StatusCode
			return nil, NewProviderError(fmt.Sprintf("Fal.ai API error: %s", falErr.Error()), &falErr)
		}
		return nil, NewProviderError(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)), nil)
	}

	// Parse response
	var imageResp FalImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		return nil, NewProviderError("failed to parse response", err)
	}

	return &imageResp, nil
}

// GenerateVideo generates a video using a Fal.ai model
func (c *FalClient) GenerateVideo(ctx context.Context, model string, request *FalRequest) (*FalVideoResponse, error) {
	endpoint := fmt.Sprintf("%s/fal-ai/%s", c.config.FalBaseURL, model)

	// Marshal request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, NewProviderError("failed to marshal request", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewNetworkError("failed to create request", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.config.FalAPIKey))

	logRequest("POST", endpoint, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Key [REDACTED]",
	})

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError("request failed", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewNetworkError("failed to read response", err)
	}

	logResponse(resp.StatusCode, string(body))

	// Check for errors
	if resp.StatusCode >= 400 {
		var falErr FalError
		if err := json.Unmarshal(body, &falErr); err == nil {
			falErr.Status = resp.StatusCode
			return nil, NewProviderError(fmt.Sprintf("Fal.ai API error: %s", falErr.Error()), &falErr)
		}
		return nil, NewProviderError(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)), nil)
	}

	// Parse response
	var videoResp FalVideoResponse
	if err := json.Unmarshal(body, &videoResp); err != nil {
		return nil, NewProviderError("failed to parse response", err)
	}

	return &videoResp, nil
}
