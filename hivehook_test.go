package hivehook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandler struct {
	handler func(req graphqlRequest) (any, error)
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/graphql" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var gqlReq graphqlRequest
	json.NewDecoder(r.Body).Decode(&gqlReq)

	data, err := m.handler(gqlReq)
	if err != nil {
		json.NewEncoder(w).Encode(graphqlResponse{
			Errors: []graphqlError{{Message: err.Error()}},
		})
		return
	}

	rawData, _ := json.Marshal(data)
	json.NewEncoder(w).Encode(graphqlResponse{Data: rawData})
}

func newTestClient(handler func(req graphqlRequest) (any, error)) (*Client, *httptest.Server) {
	srv := httptest.NewServer(&mockHandler{handler: handler})
	client := New(WithBaseURL(srv.URL), WithAPIKey("test-key"))
	return client, srv
}

func TestSourceList(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"sources": map[string]any{
				"nodes": []map[string]any{
					{"id": "src-1", "name": "GitHub", "slug": "github", "providerType": "github", "status": "ACTIVE"},
					{"id": "src-2", "name": "Stripe", "slug": "stripe", "providerType": "stripe", "status": "ACTIVE"},
				},
				"pageInfo": map[string]any{"total": 2, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	sources, pi, err := client.Sources.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Sources.List() error: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("got %d sources, want 2", len(sources))
	}
	if sources[0].Name != "GitHub" {
		t.Errorf("sources[0].Name = %q, want %q", sources[0].Name, "GitHub")
	}
	if pi.Total != 2 {
		t.Errorf("pageInfo.Total = %d, want 2", pi.Total)
	}
}

func TestSourceGet(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"source": map[string]any{
				"id": "src-1", "name": "GitHub", "slug": "github", "providerType": "github", "status": "ACTIVE",
			},
		}, nil
	})
	defer srv.Close()

	source, err := client.Sources.Get(context.Background(), "src-1")
	if err != nil {
		t.Fatalf("Sources.Get() error: %v", err)
	}
	if source.ID != "src-1" {
		t.Errorf("source.ID = %q, want %q", source.ID, "src-1")
	}
}

func TestSourceCreate(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createSource": map[string]any{
				"id": "src-new", "name": "New Source", "slug": "new-source", "providerType": "github", "status": "ACTIVE",
			},
		}, nil
	})
	defer srv.Close()

	source, err := client.Sources.Create(context.Background(), &CreateSourceInput{
		Name:         "New Source",
		Slug:         "new-source",
		ProviderType: "github",
	})
	if err != nil {
		t.Fatalf("Sources.Create() error: %v", err)
	}
	if source.Name != "New Source" {
		t.Errorf("source.Name = %q, want %q", source.Name, "New Source")
	}
}

func TestSourceDelete(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteSource": true}, nil
	})
	defer srv.Close()

	err := client.Sources.Delete(context.Background(), "src-1")
	if err != nil {
		t.Fatalf("Sources.Delete() error: %v", err)
	}
}

func TestApplicationCRUD(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createApplication": map[string]any{"id": "app-1", "name": "Test App", "uid": "test-app"},
		}, nil
	})
	defer srv.Close()

	app, err := client.Applications.Create(context.Background(), &CreateApplicationInput{Name: "Test App"})
	if err != nil {
		t.Fatalf("Applications.Create() error: %v", err)
	}
	if app.Name != "Test App" {
		t.Errorf("app.Name = %q, want %q", app.Name, "Test App")
	}
}

func TestMessageSend(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"sendMessage": map[string]any{
				"id": "msg-1", "applicationId": "app-1", "eventType": "user.created", "status": "PENDING",
			},
		}, nil
	})
	defer srv.Close()

	msg, err := client.Messages.Send(context.Background(), &SendMessageInput{
		ApplicationID: "app-1",
		EventType:     "user.created",
		Payload:       []byte(`{"user":"test"}`),
	})
	if err != nil {
		t.Fatalf("Messages.Send() error: %v", err)
	}
	if msg.EventType != "user.created" {
		t.Errorf("msg.EventType = %q, want %q", msg.EventType, "user.created")
	}
}

func TestStatusGet(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"status": map[string]any{
				"status": "ok", "version": "v0.1.0-beta", "uptime": 3600, "queueDepth": 0,
			},
		}, nil
	})
	defer srv.Close()

	status, err := client.Status.Get(context.Background())
	if err != nil {
		t.Fatalf("Status.Get() error: %v", err)
	}
	if status.Status != "ok" {
		t.Errorf("status.Status = %q, want %q", status.Status, "ok")
	}
	if status.Version != "v0.1.0-beta" {
		t.Errorf("status.Version = %q, want %q", status.Version, "v0.1.0-beta")
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("API error", func(t *testing.T) {
		client, srv := newTestClient(func(req graphqlRequest) (any, error) {
			return nil, &APIError{Message: "something went wrong"}
		})
		defer srv.Close()

		_, err := client.Sources.Get(context.Background(), "src-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		client := New(WithBaseURL(srv.URL))
		_, err := client.Sources.Get(context.Background(), "src-1")
		if err == nil {
			t.Fatal("expected auth error")
		}
		var authErr *AuthError
		if !errors.As(err, &authErr) {
			t.Errorf("expected *AuthError, got %T", err)
		}
		if !errors.Is(err, &AuthError{}) {
			t.Errorf("errors.Is(err, &AuthError{}) = false, want true")
		}
	})
}

func TestAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		rawData, _ := json.Marshal(map[string]any{"source": nil})
		json.NewEncoder(w).Encode(graphqlResponse{Data: rawData})
	}))
	defer srv.Close()

	client := New(WithBaseURL(srv.URL), WithAPIKey("htk_test123"))
	client.Sources.Get(context.Background(), "src-1")

	if gotAuth != "Bearer htk_test123" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer htk_test123")
	}
}

func TestDLQOperations(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"replayAllDLQ": map[string]any{"deliveries": 5},
		}, nil
	})
	defer srv.Close()

	result, err := client.DLQ.ReplayAll(context.Background())
	if err != nil {
		t.Fatalf("DLQ.ReplayAll() error: %v", err)
	}
	if result.Deliveries != 5 {
		t.Errorf("result.Deliveries = %d, want 5", result.Deliveries)
	}
}
