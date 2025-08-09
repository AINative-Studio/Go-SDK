package ainative

import (
	"fmt"
	"net/http"
)

// APIError represents an API error response
type APIError struct {
	// HTTP status code
	StatusCode int `json:"status_code"`
	
	// Error message
	Message string `json:"message"`
	
	// Error code for programmatic handling
	Code string `json:"code,omitempty"`
	
	// Detailed error information
	Details map[string]interface{} `json:"details,omitempty"`
	
	// Request ID for debugging
	RequestID string `json:"request_id,omitempty"`
	
	// Timestamp of the error
	Timestamp string `json:"timestamp,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("AINative API error [%d:%s]: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("AINative API error [%d]: %s", e.StatusCode, e.Message)
}

// IsAuthenticationError returns true if the error is an authentication error
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// IsRateLimitError returns true if the error is a rate limit error
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if the error is a server error (5xx)
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsClientError returns true if the error is a client error (4xx)
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsRetryable returns true if the error is retryable
func (e *APIError) IsRetryable() bool {
	// Retry on server errors and rate limits
	return e.IsServerError() || e.IsRateLimitError()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("network error: %s (caused by: %s)", e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error for field '%s': %s", e.Field, e.Message)
}

// Common error codes
const (
	ErrorCodeInvalidRequest     = "INVALID_REQUEST"
	ErrorCodeUnauthorized      = "UNAUTHORIZED"
	ErrorCodeForbidden         = "FORBIDDEN"
	ErrorCodeNotFound          = "NOT_FOUND"
	ErrorCodeRateLimit         = "RATE_LIMIT"
	ErrorCodeInternalError     = "INTERNAL_ERROR"
	ErrorCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// NewAPIError creates a new API error
func NewAPIError(statusCode int, message, code string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(message string, cause error) *NetworkError {
	return &NetworkError{
		Message: message,
		Cause:   cause,
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(field, message string) *ConfigError {
	return &ConfigError{
		Field:   field,
		Message: message,
	}
}