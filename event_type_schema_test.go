package hivehook

import (
	"context"
	"testing"
)

func TestEventTypeSchemaList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"eventTypeSchemas": map[string]any{
				"nodes":    []map[string]any{{"id": "ets-1", "eventType": "user.created", "description": "d"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	ss, _, err := client.EventTypeSchemas.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(ss) != 1 {
		t.Errorf("got %d", len(ss))
	}
}

func TestEventTypeSchemaGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"eventTypeSchema": map[string]any{"id": "ets-1", "eventType": "user.created"},
		}, nil
	})
	defer srv.Close()
	s, err := client.EventTypeSchemas.Get(context.Background(), "ets-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if s.EventType != "user.created" {
		t.Errorf("EventType = %q", s.EventType)
	}
}

func TestEventTypeSchemaCreate(t *testing.T) {
	t.Parallel()
	desc := "User created"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createEventTypeSchema": map[string]any{"id": "ets-new", "eventType": "user.created"},
		}, nil
	})
	defer srv.Close()
	s, err := client.EventTypeSchemas.Create(context.Background(), &CreateEventTypeSchemaInput{
		EventType:   "user.created",
		Description: &desc,
		Schema:      map[string]any{"type": "object"},
		Example:     map[string]any{"id": "1"},
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if s.EventType != "user.created" {
		t.Errorf("EventType = %q", s.EventType)
	}
}

func TestEventTypeSchemaUpdate(t *testing.T) {
	t.Parallel()
	et := "user.updated"
	desc := "renamed"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateEventTypeSchema": map[string]any{"id": "ets-1", "eventType": "user.updated"},
		}, nil
	})
	defer srv.Close()
	s, err := client.EventTypeSchemas.Update(context.Background(), "ets-1", &UpdateEventTypeSchemaInput{
		EventType:   &et,
		Description: &desc,
		Schema:      map[string]any{"type": "object"},
		Example:     map[string]any{"id": "1"},
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if s.EventType != "user.updated" {
		t.Errorf("EventType = %q", s.EventType)
	}
}

func TestEventTypeSchemaDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteEventTypeSchema": true}, nil
	})
	defer srv.Close()
	if err := client.EventTypeSchemas.Delete(context.Background(), "ets-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}
