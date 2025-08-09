package ainative

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
		errorMsg  string
	}{
		{
			name:      "nil config",
			config:    nil,
			wantError: true,
			errorMsg:  "config cannot be nil",
		},
		{
			name: "empty API key",
			config: &Config{
				APIKey: "",
			},
			wantError: true,
			errorMsg:  "API key is required",
		},
		{
			name: "invalid base URL",
			config: &Config{
				APIKey:  "test-key",
				BaseURL: "://invalid-url-scheme",
			},
			wantError: true,
			errorMsg:  "invalid base URL",
		},
		{
			name: "valid minimal config",
			config: &Config{
				APIKey: "test-key",
			},
			wantError: false,
		},
		{
			name: "valid full config",
			config: &Config{
				APIKey:         "test-key",
				APISecret:      "test-secret",
				BaseURL:        "https://api.example.com",
				OrganizationID: "org-123",
				ProjectID:      "proj-456",
				Timeout:        15 * time.Second,
				RateLimit:      50,
				Debug:          true,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.ZeroDB)
				assert.NotNil(t, client.AgentSwarm)
				assert.NotNil(t, client.Auth)
			}
		})
	}
}

func TestClientDefaults(t *testing.T) {
	config := &Config{
		APIKey: "test-key",
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	// Check defaults are set
	assert.Equal(t, DefaultBaseURL, client.config.BaseURL)
	assert.Equal(t, DefaultTimeout, client.config.Timeout)
	assert.Equal(t, DefaultRateLimit, client.config.RateLimit)
}

func TestClientWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	config := &Config{
		APIKey:     "test-key",
		HTTPClient: customClient,
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestClientSetProjectID(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	// Set project ID
	client.SetProjectID("proj-123")
	assert.Equal(t, "proj-123", client.config.ProjectID)

	// Clear project ID
	client.SetProjectID("")
	assert.Equal(t, "", client.config.ProjectID)
}

func TestClientSetOrganizationID(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	// Set organization ID
	client.SetOrganizationID("org-123")
	assert.Equal(t, "org-123", client.config.OrganizationID)

	// Clear organization ID
	client.SetOrganizationID("")
	assert.Equal(t, "", client.config.OrganizationID)
}

func TestClientMakeRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Contains(t, r.Header.Get("User-Agent"), "AINative-Go-SDK")
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		// Return test response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "data": {"id": "123"}}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	// Test GET request
	ctx := context.Background()
	var result map[string]interface{}
	err = client.makeRequest(ctx, "GET", "/test", nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "success", result["message"])
	assert.Equal(t, map[string]interface{}{"id": "123"}, result["data"])
}

func TestClientMakeRequestWithError(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request", "code": "INVALID_REQUEST"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	var result map[string]interface{}
	err = client.makeRequest(ctx, "GET", "/test", nil, &result)

	assert.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Equal(t, "Invalid request", apiErr.Message)
	assert.Equal(t, "INVALID_REQUEST", apiErr.Code)
	assert.Equal(t, "req-123", apiErr.RequestID)
}

func TestClientHealth(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "healthy",
			"version": "1.0.0",
			"timestamp": "2025-01-01T00:00:00Z",
			"services": {
				"database": "healthy",
				"cache": "healthy"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	health, err := client.Health(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "1.0.0", health.Version)
	assert.Equal(t, "healthy", health.Services["database"])
	assert.Equal(t, "healthy", health.Services["cache"])
}

func TestClientRetryConfiguration(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			// Return 500 error first two times
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Return success on third attempt
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:        3,
			InitialDelay:      10 * time.Millisecond,
			MaxDelay:          100 * time.Millisecond,
			BackoffMultiplier: 2.0,
			Jitter:           false,
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	var result map[string]interface{}
	err = client.makeRequest(ctx, "GET", "/test", nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, 3, attempts) // Should have retried twice and succeeded on third attempt
}

func TestClientRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:    "test-key",
		BaseURL:   server.URL,
		RateLimit: 2, // Very low rate limit for testing
	})
	require.NoError(t, err)

	ctx := context.Background()
	var result map[string]interface{}

	// First request should be immediate
	start := time.Now()
	err = client.makeRequest(ctx, "GET", "/test1", nil, &result)
	duration1 := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration1, 100*time.Millisecond)

	// Second request should be rate limited
	start = time.Now()
	err = client.makeRequest(ctx, "GET", "/test2", nil, &result)
	duration2 := time.Since(start)

	assert.NoError(t, err)
	// Should be delayed due to rate limiting
	assert.Greater(t, duration2, 100*time.Millisecond)
}

func TestClientUnsupportedMethod(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.makeRequest(ctx, "UNSUPPORTED", "/test", nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported HTTP method")
}