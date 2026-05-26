package hivehook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorIs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		err    error
		target error
	}{
		{"APIError", &APIError{Message: "x"}, &APIError{}},
		{"NotFoundError", &NotFoundError{Message: "x"}, &NotFoundError{}},
		{"ConflictError", &ConflictError{Message: "x"}, &ConflictError{}},
		{"AuthError", &AuthError{Message: "x"}, &AuthError{}},
		{"ValidationError", &ValidationError{Message: "x"}, &ValidationError{}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tc.err, tc.target) {
				t.Errorf("errors.Is(%T, %T) = false, want true", tc.err, tc.target)
			}
			// Wrapped error should still match.
			wrapped := fmt.Errorf("context: %w", tc.err)
			if !errors.Is(wrapped, tc.target) {
				t.Errorf("errors.Is(wrapped %T, %T) = false, want true", tc.err, tc.target)
			}
		})
	}
}

func TestErrorIsSentinels(t *testing.T) {
	t.Parallel()
	if !errors.Is(&NotFoundError{Message: "x"}, ErrNotFound) {
		t.Errorf("errors.Is(&NotFoundError{}, ErrNotFound) = false")
	}
	if !errors.Is(&AuthError{Message: "x"}, ErrAuth) {
		t.Errorf("errors.Is(&AuthError{}, ErrAuth) = false")
	}
	if !errors.Is(&ConflictError{Message: "x"}, ErrConflict) {
		t.Errorf("errors.Is(&ConflictError{}, ErrConflict) = false")
	}
	if !errors.Is(&ValidationError{Message: "x"}, ErrValidation) {
		t.Errorf("errors.Is(&ValidationError{}, ErrValidation) = false")
	}
	if !errors.Is(&APIError{Message: "x"}, ErrAPI) {
		t.Errorf("errors.Is(&APIError{}, ErrAPI) = false")
	}

	// Different error types should not match.
	if errors.Is(&NotFoundError{}, ErrAuth) {
		t.Errorf("NotFoundError should not match ErrAuth")
	}
}

// TestClassifyError verifies the graphql error classifier produces the right error type
// for each known extension code, and that errors.Is/As can be used to detect them.
func TestClassifyError(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		code    string
		sentinel error
	}{
		{"NOT_FOUND", "NOT_FOUND", ErrNotFound},
		{"CONFLICT", "CONFLICT", ErrConflict},
		{"VALIDATION", "VALIDATION", ErrValidation},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(graphqlResponse{
					Errors: []graphqlError{{
						Message:    "boom",
						Extensions: map[string]any{"code": tc.code},
					}},
				})
			}))
			defer srv.Close()
			client := New(WithBaseURL(srv.URL))
			_, err := client.Sources.Get(context.Background(), "src-1")
			if err == nil {
				t.Fatal("expected error")
			}
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(err, %T) = false, want true (err=%v)", tc.sentinel, err)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	t.Parallel()
	if (&APIError{StatusCode: 500, Message: "x"}).Error() == "" {
		t.Errorf("APIError empty message")
	}
	if (&APIError{Message: "x"}).Error() == "" {
		t.Errorf("APIError empty message no status")
	}
	if (&NotFoundError{Message: "x"}).Error() == "" {
		t.Errorf("NotFoundError empty")
	}
	if (&ConflictError{Message: "x"}).Error() == "" {
		t.Errorf("ConflictError empty")
	}
	if (&AuthError{Message: "x"}).Error() == "" {
		t.Errorf("AuthError empty")
	}
	if (&ValidationError{Message: "x"}).Error() == "" {
		t.Errorf("ValidationError empty")
	}
}
