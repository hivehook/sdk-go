package hivehook

import (
	"context"
	"testing"
)

func TestStreamList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["applicationId"] != "app-1" {
			t.Errorf("expected applicationId=app-1")
		}
		return map[string]any{
			"streams": map[string]any{
				"nodes":    []map[string]any{{"id": "stm-1", "name": "Events", "status": "ACTIVE"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	ss, _, err := client.Streams.List(context.Background(), &ListStreamsOptions{
		ApplicationID: "app-1",
		Status:        StreamStatusActive,
	})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(ss) != 1 {
		t.Errorf("got %d", len(ss))
	}
}

func TestStreamGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"stream": map[string]any{"id": "stm-1", "name": "Events"}}, nil
	})
	defer srv.Close()
	s, err := client.Streams.Get(context.Background(), "stm-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if s.ID != "stm-1" {
		t.Errorf("s.ID = %q", s.ID)
	}
}

func TestStreamCreate(t *testing.T) {
	t.Parallel()
	days := 30
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createStream": map[string]any{"id": "stm-new", "name": "New", "retentionDays": 30},
		}, nil
	})
	defer srv.Close()
	s, err := client.Streams.Create(context.Background(), &CreateStreamInput{
		ApplicationID: "app-1",
		Name:          "New",
		RetentionDays: &days,
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if s.RetentionDays != 30 {
		t.Errorf("RetentionDays = %d", s.RetentionDays)
	}
}

func TestStreamUpdate(t *testing.T) {
	t.Parallel()
	name := "Renamed"
	status := StreamStatusPaused
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateStream": map[string]any{"id": "stm-1", "name": "Renamed", "status": "PAUSED"},
		}, nil
	})
	defer srv.Close()
	s, err := client.Streams.Update(context.Background(), "stm-1", &UpdateStreamInput{
		Name:   &name,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if s.Name != "Renamed" {
		t.Errorf("s.Name = %q", s.Name)
	}
}

func TestStreamDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteStream": true}, nil
	})
	defer srv.Close()
	if err := client.Streams.Delete(context.Background(), "stm-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}
