package ainative

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication operations
type AuthService struct {
	client *Client
}

// NewAuthService creates a new auth service
func NewAuthService(client *Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// TokenResponse represents an authentication token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// APIKeyInfo represents API key information
type APIKeyInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Prefix      string            `json:"prefix"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	LastUsedAt  *time.Time        `json:"last_used_at,omitempty"`
	UsageCount  int64             `json:"usage_count"`
	Permissions []string          `json:"permissions"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateAPIKeyRequest represents a request to create an API key
type CreateAPIKeyRequest struct {
	Name        string            `json:"name"`
	Permissions []string          `json:"permissions,omitempty"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateAPIKeyResponse represents a response from creating an API key
type CreateAPIKeyResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`    // Full key (only returned once)
	Prefix string `json:"prefix"` // Key prefix for identification
}

// UserInfo represents user information
type UserInfo struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Organization string    `json:"organization,omitempty"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// Login authenticates with username and password
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.Username == "" {
		return nil, NewValidationError("username", "username is required", req.Username)
	}

	if req.Password == "" {
		return nil, NewValidationError("password", "password is required", "")
	}

	var result TokenResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/auth/login", req, &result)
	if err != nil {
		return nil, err
	}

	// Calculate expiration time
	result.ExpiresAt = result.IssuedAt.Add(time.Duration(result.ExpiresIn) * time.Second)

	return &result, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.RefreshToken == "" {
		return nil, NewValidationError("refresh_token", "refresh token is required", req.RefreshToken)
	}

	var result TokenResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/auth/refresh", req, &result)
	if err != nil {
		return nil, err
	}

	// Calculate expiration time
	result.ExpiresAt = result.IssuedAt.Add(time.Duration(result.ExpiresIn) * time.Second)

	return &result, nil
}

// GetUserInfo retrieves information about the current user
func (s *AuthService) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	var result UserInfo

	err := s.client.makeRequest(ctx, "GET", "/api/v1/auth/me", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListAPIKeys lists API keys for the current user
func (s *AuthService) ListAPIKeys(ctx context.Context) ([]APIKeyInfo, error) {
	var result struct {
		APIKeys []APIKeyInfo `json:"api_keys"`
	}

	err := s.client.makeRequest(ctx, "GET", "/api/v1/auth/api-keys", nil, &result)
	if err != nil {
		return nil, err
	}

	return result.APIKeys, nil
}

// CreateAPIKey creates a new API key
func (s *AuthService) CreateAPIKey(ctx context.Context, req *CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.Name == "" {
		return nil, NewValidationError("name", "name is required", req.Name)
	}

	var result CreateAPIKeyResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/auth/api-keys", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RevokeAPIKey revokes an API key
func (s *AuthService) RevokeAPIKey(ctx context.Context, keyID string) error {
	if keyID == "" {
		return NewValidationError("key_id", "key ID is required", keyID)
	}

	path := fmt.Sprintf("/api/v1/auth/api-keys/%s", keyID)

	return s.client.makeRequest(ctx, "DELETE", path, nil, nil)
}

// GetAPIKeyInfo retrieves information about a specific API key
func (s *AuthService) GetAPIKeyInfo(ctx context.Context, keyID string) (*APIKeyInfo, error) {
	if keyID == "" {
		return nil, NewValidationError("key_id", "key ID is required", keyID)
	}

	var result APIKeyInfo

	path := fmt.Sprintf("/api/v1/auth/api-keys/%s", keyID)

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
	if token == "" {
		return nil, NewValidationError("token", "token is required", token)
	}

	req := map[string]string{
		"token": token,
	}

	var result UserInfo

	err := s.client.makeRequest(ctx, "POST", "/api/v1/auth/validate", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Logout logs out the current user (invalidates tokens)
func (s *AuthService) Logout(ctx context.Context) error {
	return s.client.makeRequest(ctx, "POST", "/api/v1/auth/logout", nil, nil)
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID       string   `json:"user_id"`
	Email        string   `json:"email"`
	Organization string   `json:"organization,omitempty"`
	Role         string   `json:"role"`
	Permissions  []string `json:"permissions"`
	jwt.RegisteredClaims
}

// ParseToken parses a JWT token without validation (for inspection only)
func ParseToken(tokenString string) (*TokenClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &TokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GenerateAPIKey generates a secure API key
func GenerateAPIKey(prefix string) (string, error) {
	if prefix == "" {
		prefix = "ak"
	}

	// Generate 32 random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as base64
	keyPart := base64.URLEncoding.EncodeToString(randomBytes)

	return fmt.Sprintf("%s_%s", prefix, keyPart), nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(tokenString string) (bool, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return true, err
	}

	if claims.ExpiresAt == nil {
		return false, nil // No expiration
	}

	return claims.ExpiresAt.Time.Before(time.Now()), nil
}

// GetTokenExpirationTime gets the expiration time of a token
func GetTokenExpirationTime(tokenString string) (*time.Time, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt == nil {
		return nil, nil // No expiration
	}

	return &claims.ExpiresAt.Time, nil
}