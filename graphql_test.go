package hivehook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimitErrorDispatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("slow down"))
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(0))
	_, err := client.Sources.Get(context.Background(), "src-1")
	if err == nil {
		t.Fatal("expected error")
	}

	var rl *RateLimitError
	if !errors.As(err, &rl) {
		t.Fatalf("expected *RateLimitError, got %T: %v", err, err)
	}
	if rl.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want %d", rl.StatusCode, http.StatusTooManyRequests)
	}
	if rl.RetryAfter != 7*time.Second {
		t.Errorf("RetryAfter = %s, want 7s", rl.RetryAfter)
	}
	if !errors.Is(err, ErrRateLimit) {
		t.Errorf("errors.Is(err, ErrRateLimit) = false")
	}
}

func TestServerErrorDispatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(0))
	_, err := client.Sources.Get(context.Background(), "src-1")
	if err == nil {
		t.Fatal("expected error")
	}
	var se *ServerError
	if !errors.As(err, &se) {
		t.Fatalf("expected *ServerError, got %T: %v", err, err)
	}
	if se.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", se.StatusCode)
	}
	if !errors.Is(err, ErrServer) {
		t.Errorf("errors.Is(err, ErrServer) = false")
	}
}

func TestRetryOn5xxThenSucceed(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		raw, _ := json.Marshal(map[string]any{"source": map[string]any{"id": "src-1", "name": "ok"}})
		json.NewEncoder(w).Encode(graphqlResponse{Data: raw})
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(3))
	// Override sleep to no-op via direct field access (test-only).
	// The exported API doesn't expose it; we mutate the underlying graphql client.
	client.Sources.gql.sleep = func(time.Duration) {}

	src, err := client.Sources.Get(context.Background(), "src-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src == nil || src.Name != "ok" {
		t.Errorf("got source = %+v, want Name=ok", src)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("calls = %d, want 3", got)
	}
}

func TestRetryHonorsRetryAfter(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		raw, _ := json.Marshal(map[string]any{"source": map[string]any{"id": "src-1"}})
		json.NewEncoder(w).Encode(graphqlResponse{Data: raw})
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(2))
	var observed time.Duration
	client.Sources.gql.sleep = func(d time.Duration) { observed = d }

	_, err := client.Sources.Get(context.Background(), "src-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if observed != 1*time.Second {
		t.Errorf("backoff sleep = %s, want 1s (Retry-After)", observed)
	}
}

func TestRetryDoesNotRetryAuth(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(5))
	client.Sources.gql.sleep = func(time.Duration) {}

	_, err := client.Sources.Get(context.Background(), "src-1")
	if err == nil {
		t.Fatal("expected auth error")
	}
	if !errors.Is(err, ErrAuth) {
		t.Errorf("expected ErrAuth, got %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("calls = %d, want 1 (no retry on auth)", got)
	}
}

func TestRetryExhausted(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithRetry(2))
	client.Sources.gql.sleep = func(time.Duration) {}

	_, err := client.Sources.Get(context.Background(), "src-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrServer) {
		t.Errorf("expected ErrServer, got %v", err)
	}
	// 1 initial + 2 retries = 3
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("calls = %d, want 3", got)
	}
}

func TestGraphQLErrorsPreserved(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(graphqlResponse{
			Errors: []graphqlError{
				{
					Message:    "first",
					Path:       []any{"source"},
					Extensions: map[string]any{"code": "VALIDATION", "field": "slug"},
				},
				{
					Message:    "second",
					Extensions: map[string]any{"code": "ANOTHER"},
				},
			},
		})
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL))
	_, err := client.Sources.Get(context.Background(), "src-1")
	if err == nil {
		t.Fatal("expected error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.GraphQLErrors) != 2 {
		t.Fatalf("GraphQLErrors len = %d, want 2", len(ve.GraphQLErrors))
	}
	if ve.GraphQLErrors[0].Code != "VALIDATION" {
		t.Errorf("[0].Code = %q, want VALIDATION", ve.GraphQLErrors[0].Code)
	}
	if ve.GraphQLErrors[0].Extensions["field"] != "slug" {
		t.Errorf("[0].Extensions[field] = %v, want slug", ve.GraphQLErrors[0].Extensions["field"])
	}
	if ve.GraphQLErrors[1].Message != "second" {
		t.Errorf("[1].Message = %q, want second", ve.GraphQLErrors[1].Message)
	}
}

func TestParseRetryAfter(t *testing.T) {
	if got := parseRetryAfter("5"); got != 5*time.Second {
		t.Errorf("parseRetryAfter(5) = %s, want 5s", got)
	}
	if got := parseRetryAfter(""); got != 0 {
		t.Errorf("parseRetryAfter(empty) = %s, want 0", got)
	}
	if got := parseRetryAfter("garbage"); got != 0 {
		t.Errorf("parseRetryAfter(garbage) = %s, want 0", got)
	}
	if got := parseRetryAfter("-1"); got != 0 {
		t.Errorf("parseRetryAfter(-1) = %s, want 0", got)
	}
	// HTTP date in the future.
	future := time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)
	if got := parseRetryAfter(future); got <= 0 || got > 3*time.Second {
		t.Errorf("parseRetryAfter(future date) = %s, want ~2s", got)
	}
}

func TestWithTimeout(t *testing.T) {
	// Verify WithTimeout sets the http client timeout.
	c := New(WithTimeout(123 * time.Millisecond))
	if c.http.Timeout != 123*time.Millisecond {
		t.Errorf("client timeout = %s, want 123ms", c.http.Timeout)
	}
}

func TestWithRetryClamped(t *testing.T) {
	// Verify negative values are clamped to 0.
	c := New(WithRetry(-5))
	if c.maxRetries != 0 {
		t.Errorf("maxRetries = %d, want 0", c.maxRetries)
	}
}

// silence unused import warnings if test bodies change later.
var _ = strconv.Itoa
