// Package ainative provides the official Go SDK for AINative Studio APIs
package ainative

import (
	"context"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the main AINative client
type Client struct {
	config  *Config
	resty   *resty.Client
	ZeroDB  *ZeroDBService
	AgentSwarm *AgentSwarmService
}

// Config holds the client configuration
type Config struct {
	APIKey    string
	APISecret string
	BaseURL   string
	OrgID     string
	Timeout   time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

// NewClient creates a new AINative client with default configuration
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(&Config{
		APIKey:     apiKey,
		BaseURL:    "https://api.ainative.studio",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: time.Second,
	})
}

// NewClientFromEnv creates a client from environment variables
func NewClientFromEnv() *Client {
	return NewClient(os.Getenv("AINATIVE_API_KEY"))
}

// NewClientWithConfig creates a client with custom configuration
func NewClientWithConfig(config *Config) *Client {
	client := resty.New().
		SetBaseURL(config.BaseURL + "/api/v1").
		SetTimeout(config.Timeout).
		SetRetryCount(config.MaxRetries).
		SetRetryWaitTime(config.RetryDelay).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+config.APIKey)

	if config.OrgID != "" {
		client.SetHeader("X-Organization-ID", config.OrgID)
	}

	c := &Client{
		config: config,
		resty:  client,
	}

	// Initialize services
	c.ZeroDB = &ZeroDBService{client: c}
	c.AgentSwarm = &AgentSwarmService{client: c}

	return c
}

// HealthCheck checks if the API is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.resty.R().
		SetContext(ctx).
		Get("/health")
	return err
}