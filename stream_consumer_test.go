package hivehook

import (
	"context"
	"testing"
)

func TestStreamConsumerList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["streamId"] != "stm-1" {
			t.Errorf("expected streamId=stm-1")
		}
		return map[string]any{
			"streamConsumers": map[string]any{
				"nodes":    []map[string]any{{"id": "cns-1", "streamId": "stm-1", "name": "c1"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	cs, _, err := client.StreamConsumers.List(context.Background(), "stm-1", nil)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(cs) != 1 {
		t.Errorf("got %d", len(cs))
	}
}

func TestStreamConsumerGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"streamConsumer": map[string]any{"id": "cns-1", "name": "c1"}}, nil
	})
	defer srv.Close()
	c, err := client.StreamConsumers.Get(context.Background(), "cns-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if c.ID != "cns-1" {
		t.Errorf("c.ID = %q", c.ID)
	}
}

func TestStreamConsumerCreate(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createStreamConsumer": map[string]any{"id": "cns-new", "streamId": "stm-1", "name": "c1"},
		}, nil
	})
	defer srv.Close()
	c, err := client.StreamConsumers.Create(context.Background(), &CreateStreamConsumerInput{
		StreamID: "stm-1",
		Name:     "c1",
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if c.Name != "c1" {
		t.Errorf("c.Name = %q", c.Name)
	}
}

func TestStreamConsumerDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteStreamConsumer": true}, nil
	})
	defer srv.Close()
	if err := client.StreamConsumers.Delete(context.Background(), "cns-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestStreamConsumerAdvanceCursor(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if v, ok := req.Variables["sequence"]; !ok || v != float64(42) {
			t.Errorf("expected sequence=42, got %v (%T)", v, v)
		}
		return map[string]any{
			"advanceConsumerCursor": map[string]any{"id": "cns-1", "cursorSequence": 42},
		}, nil
	})
	defer srv.Close()
	c, err := client.StreamConsumers.AdvanceCursor(context.Background(), "cns-1", 42)
	if err != nil {
		t.Fatalf("AdvanceCursor() error: %v", err)
	}
	if c.CursorSequence != 42 {
		t.Errorf("CursorSequence = %d", c.CursorSequence)
	}
}
