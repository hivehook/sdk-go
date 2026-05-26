package hivehook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// graphqlClient performs HTTP+JSON GraphQL requests against the Hivehook API.
type graphqlClient struct {
	baseURL    string
	apiKey     string
	http       *http.Client
	maxRetries int
	sleep      func(time.Duration) // for tests; defaults to time.Sleep
}

// graphqlRequest is the wire-format payload for a single GraphQL operation.
type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// graphqlResponse is the wire-format reply for a GraphQL operation.
type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphqlError  `json:"errors"`
}

// graphqlError mirrors a single error entry in a GraphQL response.
type graphqlError struct {
	Message    string         `json:"message"`
	Path       []any          `json:"path"`
	Extensions map[string]any `json:"extensions"`
}

// toPublic converts a wire error to the exported GraphQLError.
func (e graphqlError) toPublic() GraphQLError {
	code, _ := e.Extensions["code"].(string)
	return GraphQLError{
		Message:    e.Message,
		Code:       code,
		Path:       e.Path,
		Extensions: e.Extensions,
	}
}

// toPublicSlice converts a slice of wire errors to exported GraphQLErrors.
func toPublicErrors(errs []graphqlError) []GraphQLError {
	if len(errs) == 0 {
		return nil
	}
	out := make([]GraphQLError, len(errs))
	for i, e := range errs {
		out[i] = e.toPublic()
	}
	return out
}

// do executes a GraphQL operation, decoding the data into result. It honors
// the configured retry policy for transient HTTP failures (5xx, 429) and
// respects Retry-After.
func (c *graphqlClient) do(ctx context.Context, query string, variables map[string]any, result any) error {
	body, err := json.Marshal(graphqlRequest{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("hivehook: marshal request: %w", err)
	}

	sleep := c.sleep
	if sleep == nil {
		sleep = time.Sleep
	}

	var lastErr error
	attempts := c.maxRetries + 1
	for attempt := 0; attempt < attempts; attempt++ {
		err := c.doOnce(ctx, body, result)
		if err == nil {
			return nil
		}
		lastErr = err

		// Never retry permanent classification errors.
		if !isRetryable(err) {
			return err
		}

		// Last attempt: return error without sleeping.
		if attempt == attempts-1 {
			return err
		}

		// Compute backoff: prefer Retry-After (rate-limit), else exponential.
		delay := backoff(attempt, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		sleep(delay)
	}
	return lastErr
}

// doOnce performs a single HTTP round-trip and decodes the response.
func (c *graphqlClient) doOnce(ctx context.Context, body []byte, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hivehook: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("hivehook: send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("hivehook: read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &AuthError{StatusCode: resp.StatusCode, Message: "unauthorized"}
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
			RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}
	if resp.StatusCode >= 500 {
		return &ServerError{StatusCode: resp.StatusCode, Message: string(respBody)}
	}
	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode, Message: string(respBody)}
	}

	var gqlResp graphqlResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return fmt.Errorf("hivehook: unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return classifyError(gqlResp.Errors)
	}

	if result != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("hivehook: unmarshal data: %w", err)
		}
	}

	return nil
}

// classifyError maps a GraphQL errors array onto the appropriate typed SDK error.
// All entries are preserved on the returned error's GraphQLErrors field.
func classifyError(errs []graphqlError) error {
	if len(errs) == 0 {
		return nil
	}
	first := errs[0]
	msg := first.Message
	all := toPublicErrors(errs)
	if code, ok := first.Extensions["code"].(string); ok {
		switch code {
		case "NOT_FOUND":
			return &NotFoundError{Message: msg, GraphQLErrors: all}
		case "CONFLICT":
			return &ConflictError{Message: msg, GraphQLErrors: all}
		case "VALIDATION":
			return &ValidationError{Message: msg, GraphQLErrors: all}
		case "UNAUTHENTICATED", "UNAUTHORIZED", "FORBIDDEN":
			return &AuthError{Message: msg, GraphQLErrors: all}
		}
	}
	return &APIError{Message: msg, GraphQLErrors: all}
}

// isRetryable reports whether an error from doOnce should be retried.
func isRetryable(err error) bool {
	switch err.(type) {
	case *RateLimitError, *ServerError:
		return true
	}
	return false
}

// backoff computes the delay before the next retry attempt. It prefers the
// server-provided Retry-After value when present.
func backoff(attempt int, err error) time.Duration {
	if rl, ok := err.(*RateLimitError); ok && rl.RetryAfter > 0 {
		return rl.RetryAfter
	}
	// Exponential backoff: 100ms, 200ms, 400ms... capped at 5s.
	d := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// parseRetryAfter interprets the Retry-After header as either seconds or an
// HTTP date. Returns zero if the header is empty or unparseable.
func parseRetryAfter(h string) time.Duration {
	h = strings.TrimSpace(h)
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}
