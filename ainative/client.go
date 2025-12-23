package ainative

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

const (
	// DefaultBaseURL is the default AINative API base URL
	DefaultBaseURL = "https://api.ainative.studio"
	
	// DefaultTimeout is the default HTTP timeout
	DefaultTimeout = 30 * time.Second
	
	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3
	
	// DefaultRateLimit is the default rate limit (requests per second)
	DefaultRateLimit = 100
	
	// UserAgent is the SDK user agent string
	UserAgent = "AINative-Go-SDK/1.0.0"
)

// Client represents the main AINative API client
type Client struct {
	// HTTP client
	httpClient *resty.Client
	
	// Configuration
	config *Config
	
	// Rate limiter
	rateLimiter *rate.Limiter
	
	// OpenTelemetry tracer
	tracer trace.Tracer
	
	// API service clients
	ZeroDB              *ZeroDBService
	AgentSwarm          *AgentSwarmService
	AgentOrchestration  *AgentOrchestrationService
	AgentCoordination   *AgentCoordinationService
	AgentLearning       *AgentLearningService
	AgentState          *AgentStateService
	Auth                *AuthService
}

// Config holds the configuration for the AINative client
type Config struct {
	// Required: API key for authentication
	APIKey string
	
	// Optional: API secret for enhanced security
	APISecret string
	
	// Optional: Base URL for the API (defaults to production)
	BaseURL string
	
	// Optional: Organization ID for multi-tenant setups
	OrganizationID string
	
	// Optional: Project ID for project-scoped operations
	ProjectID string
	
	// Optional: Custom HTTP client
	HTTPClient *http.Client
	
	// Optional: Request timeout (defaults to 30s)
	Timeout time.Duration
	
	// Optional: Retry configuration
	RetryConfig *RetryConfig
	
	// Optional: Rate limit (requests per second)
	RateLimit int
	
	// Optional: OpenTelemetry tracer
	Tracer trace.Tracer
	
	// Optional: Debug mode
	Debug bool
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	// Maximum number of retry attempts
	MaxRetries int
	
	// Initial delay between retries
	InitialDelay time.Duration
	
	// Maximum delay between retries
	MaxDelay time.Duration
	
	// Backoff multiplier
	BackoffMultiplier float64
	
	// Add jitter to prevent thundering herd
	Jitter bool
}

// NewClient creates a new AINative API client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	
	// Set defaults
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	
	if config.RateLimit == 0 {
		config.RateLimit = DefaultRateLimit
	}
	
	// Validate base URL
	_, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	
	// Create HTTP client
	var baseHTTPClient *http.Client
	if config.HTTPClient != nil {
		baseHTTPClient = config.HTTPClient
	} else {
		// Set default HTTP client with connection pooling
		baseHTTPClient = &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	}
	
	httpClient := resty.NewWithClient(baseHTTPClient)
	
	// Configure client
	httpClient.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")
	
	// Set authentication header
	if config.APISecret != "" {
		// Use API key + secret for enhanced authentication
		httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s:%s", config.APIKey, config.APISecret))
	} else {
		httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
	}
	
	// Set organization header if provided
	if config.OrganizationID != "" {
		httpClient.SetHeader("X-Organization-ID", config.OrganizationID)
	}
	
	// Configure retry
	retryConfig := config.RetryConfig
	if retryConfig == nil {
		retryConfig = &RetryConfig{
			MaxRetries:        DefaultMaxRetries,
			InitialDelay:      100 * time.Millisecond,
			MaxDelay:          10 * time.Second,
			BackoffMultiplier: 2.0,
			Jitter:           true,
		}
	}
	
	httpClient.
		SetRetryCount(retryConfig.MaxRetries).
		SetRetryWaitTime(retryConfig.InitialDelay).
		SetRetryMaxWaitTime(retryConfig.MaxDelay).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors
			if err != nil {
				return true
			}
			
			// Retry on 5xx errors and 429 (rate limit)
			return r.StatusCode() >= 500 || r.StatusCode() == 429
		})
	
	// Enable debug mode if requested
	if config.Debug {
		httpClient.SetDebug(true)
	}
	
	// Create rate limiter
	rateLimiter := rate.NewLimiter(rate.Limit(config.RateLimit), config.RateLimit)
	
	// Set up tracer
	tracer := config.Tracer
	if tracer == nil {
		tracer = otel.Tracer("ainative-go-sdk")
	}
	
	// Create client
	client := &Client{
		httpClient:  httpClient,
		config:      config,
		rateLimiter: rateLimiter,
		tracer:      tracer,
	}
	
	// Initialize service clients
	client.ZeroDB = NewZeroDBService(client)
	client.AgentSwarm = NewAgentSwarmService(client)
	client.AgentOrchestration = NewAgentOrchestrationService(client)
	client.AgentCoordination = NewAgentCoordinationService(client)
	client.AgentLearning = NewAgentLearningService(client)
	client.AgentState = NewAgentStateService(client)
	client.Auth = NewAuthService(client)

	return client, nil
}

// makeRequest performs an HTTP request with rate limiting and tracing
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Start tracing span
	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("ainative.%s %s", method, path))
	defer span.End()
	
	// Apply rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait failed: %w", err)
	}
	
	// Create request
	req := c.httpClient.R().SetContext(ctx)
	
	// Set body if provided
	if body != nil {
		req.SetBody(body)
	}
	
	// Set result if provided
	if result != nil {
		req.SetResult(result)
	}
	
	// Set error response handler
	var apiError APIError
	req.SetError(&apiError)
	
	// Make request
	var resp *resty.Response
	var err error
	
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "DELETE":
		resp, err = req.Delete(path)
	case "PATCH":
		resp, err = req.Patch(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}
	
	// Handle network errors
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	
	// Handle API errors
	if resp.StatusCode() >= 400 {
		if apiError.Message != "" {
			apiError.StatusCode = resp.StatusCode()
			return &apiError
		}
		
		// Fallback error
		return &APIError{
			StatusCode: resp.StatusCode(),
			Message:    fmt.Sprintf("API request failed with status %d", resp.StatusCode()),
			RequestID:  resp.Header().Get("X-Request-ID"),
		}
	}
	
	return nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// SetProjectID sets the default project ID for operations
func (c *Client) SetProjectID(projectID string) {
	c.config.ProjectID = projectID
	if projectID != "" {
		c.httpClient.SetHeader("X-Project-ID", projectID)
	} else {
		c.httpClient.SetHeader("X-Project-ID", "")
	}
}

// SetOrganizationID sets the default organization ID for operations  
func (c *Client) SetOrganizationID(orgID string) {
	c.config.OrganizationID = orgID
	if orgID != "" {
		c.httpClient.SetHeader("X-Organization-ID", orgID)
	} else {
		c.httpClient.SetHeader("X-Organization-ID", "")
	}
}

// Health checks the API health
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var result HealthResponse
	
	err := c.makeRequest(ctx, "GET", "/health", nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// HealthResponse represents the API health status
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}