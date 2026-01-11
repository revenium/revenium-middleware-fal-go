package revenium

import "time"

// OperationType represents the type of AI operation
type OperationType string

const (
	OperationTypeImage OperationType = "IMAGE"
	OperationTypeVideo OperationType = "VIDEO"
)

// FalRequest represents a request to the Fal.ai API
type FalRequest struct {
	Prompt              string                 `json:"prompt"`
	ImageSize           string                 `json:"image_size,omitempty"`
	NumInferenceSteps   int                    `json:"num_inference_steps,omitempty"`
	GuidanceScale       float64                `json:"guidance_scale,omitempty"`
	NumImages           int                    `json:"num_images,omitempty"`
	Seed                *int                   `json:"seed,omitempty"`
	EnableSafetyChecker bool                   `json:"enable_safety_checker,omitempty"`
	Duration            string                 `json:"duration,omitempty"`    // Video duration: "5" or "10" seconds
	AspectRatio         string                 `json:"aspect_ratio,omitempty"` // Video aspect ratio: "16:9", "9:16", "1:1"
	AdditionalParams    map[string]interface{} `json:"-"`
}

// FalImageResponse represents the response from Fal.ai image generation
type FalImageResponse struct {
	Images      []FalImage `json:"images"`
	Seed        int        `json:"seed,omitempty"`
	TimeTaken   float64    `json:"timeTaken,omitempty"`
	HasNSFWContent []bool  `json:"has_nsfw_content,omitempty"`
	Prompt      string     `json:"prompt,omitempty"`
}

// FalImage represents a single generated image
type FalImage struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type,omitempty"`
}

// FalVideoResponse represents the response from Fal.ai video generation
type FalVideoResponse struct {
	Video       FalVideo `json:"video"`
	Prompt      string   `json:"prompt,omitempty"`
	TimeTaken   float64  `json:"timeTaken,omitempty"`
}

// FalVideo represents a generated video
type FalVideo struct {
	URL         string  `json:"url"`
	Duration    float64 `json:"duration,omitempty"`
	Width       int     `json:"width,omitempty"`
	Height      int     `json:"height,omitempty"`
	ContentType string  `json:"content_type,omitempty"`
}

// FalError represents an error response from Fal.ai
type FalError struct {
	ErrorText string `json:"error"`
	Message   string `json:"message,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// Error implements the error interface
func (e *FalError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.ErrorText
}

// MeteringPayload represents the payload sent to Revenium API
type MeteringPayload struct {
	// Required fields
	StopReason       string    `json:"stopReason"`
	CostType         string    `json:"costType"`
	OperationType    string    `json:"operationType"`
	Model            string    `json:"model"`
	Provider         string    `json:"provider"`
	ModelSource      string    `json:"modelSource"`
	TransactionID    string    `json:"transactionId"`
	RequestTime      time.Time `json:"requestTime"`
	ResponseTime     time.Time `json:"responseTime"`
	RequestDuration  int64     `json:"requestDuration"` // milliseconds
	MiddlewareSource string    `json:"middlewareSource"`

	// Image-specific billing fields (TOP LEVEL per API contract)
	ActualImageCount    *int `json:"actualImageCount,omitempty"`
	RequestedImageCount *int `json:"requestedImageCount,omitempty"`

	// Video-specific billing fields (TOP LEVEL per API contract)
	DurationSeconds          *float64 `json:"durationSeconds,omitempty"`
	RequestedDurationSeconds *float64 `json:"requestedDurationSeconds,omitempty"` // Required for PER_SECOND billing

	// Image/Video metadata (in attributes, not billing)
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Optional business context
	OrganizationID   string                 `json:"organizationId,omitempty"`
	ProductID        string                 `json:"productId,omitempty"`
	TaskType         string                 `json:"taskType,omitempty"`
	Agent            string                 `json:"agent,omitempty"`
	SubscriptionID   string                 `json:"subscriptionId,omitempty"`
	TraceID          string                 `json:"traceId,omitempty"`
	// Distributed tracing fields
	ParentTransactionID string `json:"parentTransactionId,omitempty"`
	TraceType           string `json:"traceType,omitempty"`
	TraceName           string `json:"traceName,omitempty"`
	Environment         string `json:"environment,omitempty"`
	Region              string `json:"region,omitempty"`
	RetryNumber         *int   `json:"retryNumber,omitempty"`
	CredentialAlias     string `json:"credentialAlias,omitempty"`
	Subscriber       map[string]interface{} `json:"subscriber,omitempty"`
	TaskID           string                 `json:"taskId,omitempty"`
	// Multimodal job identifiers
	VideoJobID       string                 `json:"videoJobId,omitempty"`
	AudioJobID       string                 `json:"audioJobId,omitempty"`

	// Quality metrics
	ResponseQualityScore *float64 `json:"responseQualityScore,omitempty"`

	// Cost overrides
	TotalCost        *float64 `json:"totalCost,omitempty"`
}
