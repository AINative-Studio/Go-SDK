package ainative

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/login", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", req.Username)
		assert.Equal(t, "password123", req.Password)

		now := time.Now()
		response := TokenResponse{
			AccessToken:  "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "refresh_token_123",
			Scope:        "read write",
			IssuedAt:     now,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &LoginRequest{
		Username: "test@example.com",
		Password: "password123",
	}

	token, err := client.Auth.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "Bearer", token.TokenType)
	assert.Equal(t, 3600, token.ExpiresIn)
	assert.Equal(t, "refresh_token_123", token.RefreshToken)
	assert.NotZero(t, token.ExpiresAt)
}

func TestAuthService_RefreshToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/refresh", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req RefreshTokenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "refresh_token_123", req.RefreshToken)

		now := time.Now()
		response := TokenResponse{
			AccessToken: "new_access_token_456",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			IssuedAt:    now,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
	}

	token, err := client.Auth.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "new_access_token_456", token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)
}

func TestAuthService_GetUserInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/me", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		response := UserInfo{
			ID:           "user_123",
			Email:        "test@example.com",
			Name:         "Test User",
			Organization: "Test Org",
			Role:         "admin",
			IsActive:     true,
			CreatedAt:    time.Now(),
			LastLoginAt:  &[]time.Time{time.Now()}[0],
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	userInfo, err := client.Auth.GetUserInfo(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "user_123", userInfo.ID)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "admin", userInfo.Role)
	assert.True(t, userInfo.IsActive)
}

func TestAuthService_ListAPIKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/api-keys", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := struct {
			APIKeys []APIKeyInfo `json:"api_keys"`
		}{
			APIKeys: []APIKeyInfo{
				{
					ID:          "key_1",
					Name:        "Production Key",
					Prefix:      "ak_prod",
					IsActive:    true,
					CreatedAt:   time.Now(),
					UsageCount:  1500,
					Permissions: []string{"read", "write"},
				},
				{
					ID:          "key_2",
					Name:        "Development Key",
					Prefix:      "ak_dev",
					IsActive:    true,
					CreatedAt:   time.Now(),
					UsageCount:  250,
					Permissions: []string{"read"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	keys, err := client.Auth.ListAPIKeys(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, keys)
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, "Production Key", keys[0].Name)
	assert.Equal(t, "ak_prod", keys[0].Prefix)
	assert.Equal(t, int64(1500), keys[0].UsageCount)
}

func TestAuthService_CreateAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/api-keys", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateAPIKeyRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "New API Key", req.Name)
		assert.Equal(t, []string{"read", "write"}, req.Permissions)

		response := CreateAPIKeyResponse{
			ID:     "key_new",
			Name:   req.Name,
			Key:    "ak_new_1234567890abcdef1234567890abcdef12345678",
			Prefix: "ak_new",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	req := &CreateAPIKeyRequest{
		Name:        "New API Key",
		Permissions: []string{"read", "write"},
		ExpiresAt:   &expiresAt,
		Metadata: map[string]string{
			"purpose": "testing",
		},
	}

	key, err := client.Auth.CreateAPIKey(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "key_new", key.ID)
	assert.Equal(t, "New API Key", key.Name)
	assert.Equal(t, "ak_new", key.Prefix)
	assert.Contains(t, key.Key, "ak_new_")
}

func TestAuthService_RevokeAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/api-keys/key_123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.Auth.RevokeAPIKey(ctx, "key_123")

	assert.NoError(t, err)
}

func TestAuthService_GetAPIKeyInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/api-keys/key_123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := APIKeyInfo{
			ID:          "key_123",
			Name:        "Test Key",
			Prefix:      "ak_test",
			IsActive:    true,
			CreatedAt:   time.Now(),
			UsageCount:  42,
			Permissions: []string{"read"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	keyInfo, err := client.Auth.GetAPIKeyInfo(ctx, "key_123")

	assert.NoError(t, err)
	assert.NotNil(t, keyInfo)
	assert.Equal(t, "key_123", keyInfo.ID)
	assert.Equal(t, "Test Key", keyInfo.Name)
	assert.Equal(t, int64(42), keyInfo.UsageCount)
}

func TestAuthService_ValidateToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/validate", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req map[string]string
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "test_token_123", req["token"])

		response := UserInfo{
			ID:       "user_123",
			Email:    "test@example.com",
			Name:     "Test User",
			Role:     "user",
			IsActive: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	userInfo, err := client.Auth.ValidateToken(ctx, "test_token_123")

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "user_123", userInfo.ID)
	assert.Equal(t, "test@example.com", userInfo.Email)
}

func TestAuthService_Logout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/auth/logout", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "logged out"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.Auth.Logout(ctx)

	assert.NoError(t, err)
}

func TestAuthService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test Login with nil request
	_, err = client.Auth.Login(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Login with empty username
	_, err = client.Auth.Login(ctx, &LoginRequest{Username: "", Password: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")

	// Test Login with empty password
	_, err = client.Auth.Login(ctx, &LoginRequest{Username: "test", Password: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password is required")

	// Test RefreshToken with nil request
	_, err = client.Auth.RefreshToken(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test RefreshToken with empty token
	_, err = client.Auth.RefreshToken(ctx, &RefreshTokenRequest{RefreshToken: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token is required")

	// Test CreateAPIKey with nil request
	_, err = client.Auth.CreateAPIKey(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test CreateAPIKey with empty name
	_, err = client.Auth.CreateAPIKey(ctx, &CreateAPIKeyRequest{Name: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test RevokeAPIKey with empty key ID
	err = client.Auth.RevokeAPIKey(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key ID is required")

	// Test GetAPIKeyInfo with empty key ID
	_, err = client.Auth.GetAPIKeyInfo(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key ID is required")

	// Test ValidateToken with empty token
	_, err = client.Auth.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is required")
}

func TestGenerateAPIKey(t *testing.T) {
	// Test with default prefix
	key, err := GenerateAPIKey("")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.True(t, len(key) > 10)
	assert.Contains(t, key, "ak_")

	// Test with custom prefix
	key, err = GenerateAPIKey("custom")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.Contains(t, key, "custom_")

	// Test uniqueness
	key1, err := GenerateAPIKey("test")
	assert.NoError(t, err)
	key2, err := GenerateAPIKey("test")
	assert.NoError(t, err)
	assert.NotEqual(t, key1, key2)
}

func TestParseToken(t *testing.T) {
	// This is a test JWT token (not signed properly, just for parsing test)
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcl8xMjMiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImV4cCI6MTcwNjcyNDAwMH0.fake_signature"
	
	claims, err := ParseToken(testToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user_123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "user", claims.Role)
}

func TestIsTokenExpired(t *testing.T) {
	// Test with expired token (exp: 1706724000 = 2024-01-31)
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcl8xMjMiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImV4cCI6MTcwNjcyNDAwMH0.fake_signature"
	
	expired, err := IsTokenExpired(expiredToken)
	assert.NoError(t, err)
	assert.True(t, expired)

	// Test with invalid token
	_, err = IsTokenExpired("invalid_token")
	assert.Error(t, err)
}

func TestGetTokenExpirationTime(t *testing.T) {
	// Test with token that has expiration
	tokenWithExp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcl8xMjMiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImV4cCI6MTcwNjcyNDAwMH0.fake_signature"
	
	expTime, err := GetTokenExpirationTime(tokenWithExp)
	assert.NoError(t, err)
	assert.NotNil(t, expTime)
	assert.Equal(t, int64(1706724000), expTime.Unix())

	// Test with invalid token
	_, err = GetTokenExpirationTime("invalid_token")
	assert.Error(t, err)
}