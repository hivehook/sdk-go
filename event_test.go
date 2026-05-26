package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestEventList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["sourceId"] != "src-1" {
			t.Errorf("expected sourceId=src-1")
		}
		if req.Variables["eventType"] != "user.created" {
			t.Errorf("expected eventType=user.created")
		}
		return map[string]any{
			"events": map[string]any{
				"nodes": []map[string]any{
					{"id": "ev-1", "sourceId": "src-1", "eventType": "user.created", "status": "DELIVERED"},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	events, _, err := client.Events.List(context.Background(), &ListEventsOptions{
		SourceID:  "src-1",
		EventType: "user.created",
		Status:    EventStatusDelivered,
	})
	if err != nil {
		t.Fatalf("Events.List() error: %v", err)
	}
	if len(events) != 1 || events[0].EventType != "user.created" {
		t.Errorf("unexpected events: %+v", events)
	}
}

func TestEventGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "event(id:") {
			t.Errorf("expected event(id:...) query")
		}
		return map[string]any{
			"event": map[string]any{"id": "ev-1", "eventType": "user.created"},
		}, nil
	})
	defer srv.Close()
	ev, err := client.Events.Get(context.Background(), "ev-1")
	if err != nil {
		t.Fatalf("Events.Get() error: %v", err)
	}
	if ev.ID != "ev-1" {
		t.Errorf("ev.ID = %q", ev.ID)
	}
}
