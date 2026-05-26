package hivehook

import (
	"fmt"
	"time"
)

// Sentinel error values for use with errors.Is.
var (
	ErrAPI        error = &APIError{}
	ErrNotFound   error = &NotFoundError{}
	ErrConflict   error = &ConflictError{}
	ErrAuth       error = &AuthError{}
	ErrValidation error = &ValidationError{}
	ErrRateLimit  error = &RateLimitError{}
	ErrServer     error = &ServerError{}
)

// GraphQLError is a single error entry from a GraphQL response.
type GraphQLError struct {
	Message    string         `json:"message"`
	Code       string         `json:"code"`
	Path       []any          `json:"path"`
	Extensions map[string]any `json:"extensions"`
}

// APIError is returned for generic API failures and unclassified errors.
type APIError struct {
	StatusCode    int
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *APIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("hivehook: API error (HTTP %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("hivehook: API error: %s", e.Message)
}

// Is supports errors.Is comparison against the *APIError sentinel.
func (e *APIError) Is(target error) bool {
	_, ok := target.(*APIError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *APIError) Unwrap() error { return e.cause }

// NotFoundError is returned when the requested resource does not exist.
type NotFoundError struct {
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("hivehook: not found: %s", e.Message)
}

// Is supports errors.Is comparison against the *NotFoundError sentinel.
func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *NotFoundError) Unwrap() error { return e.cause }

// ConflictError is returned when the request conflicts with current state.
type ConflictError struct {
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *ConflictError) Error() string {
	return fmt.Sprintf("hivehook: conflict: %s", e.Message)
}

// Is supports errors.Is comparison against the *ConflictError sentinel.
func (e *ConflictError) Is(target error) bool {
	_, ok := target.(*ConflictError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *ConflictError) Unwrap() error { return e.cause }

// AuthError is returned for 401/403 responses and authentication failures.
type AuthError struct {
	StatusCode    int
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *AuthError) Error() string {
	return fmt.Sprintf("hivehook: auth error: %s", e.Message)
}

// Is supports errors.Is comparison against the *AuthError sentinel.
func (e *AuthError) Is(target error) bool {
	_, ok := target.(*AuthError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *AuthError) Unwrap() error { return e.cause }

// ValidationError is returned when the request payload fails validation.
type ValidationError struct {
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("hivehook: validation error: %s", e.Message)
}

// Is supports errors.Is comparison against the *ValidationError sentinel.
func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *ValidationError) Unwrap() error { return e.cause }

// RateLimitError is returned when the API returns HTTP 429. RetryAfter is parsed
// from the Retry-After response header when present (zero otherwise).
type RateLimitError struct {
	StatusCode    int
	Message       string
	RetryAfter    time.Duration
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("hivehook: rate limited (HTTP %d, retry after %s): %s", e.StatusCode, e.RetryAfter, e.Message)
	}
	return fmt.Sprintf("hivehook: rate limited (HTTP %d): %s", e.StatusCode, e.Message)
}

// Is supports errors.Is comparison against the *RateLimitError sentinel.
func (e *RateLimitError) Is(target error) bool {
	_, ok := target.(*RateLimitError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *RateLimitError) Unwrap() error { return e.cause }

// ServerError is returned for HTTP 5xx server errors.
type ServerError struct {
	StatusCode    int
	Message       string
	GraphQLErrors []GraphQLError
	cause         error
}

// Error returns the formatted error message.
func (e *ServerError) Error() string {
	return fmt.Sprintf("hivehook: server error (HTTP %d): %s", e.StatusCode, e.Message)
}

// Is supports errors.Is comparison against the *ServerError sentinel.
func (e *ServerError) Is(target error) bool {
	_, ok := target.(*ServerError)
	return ok
}

// Unwrap returns the underlying cause, if any.
func (e *ServerError) Unwrap() error { return e.cause }
