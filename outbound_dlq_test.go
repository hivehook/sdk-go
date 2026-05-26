package hivehook

import (
	"context"
	"testing"
)

func TestOutboundDLQList(t *testing.T) {
	t.Parallel()
	msgID := "msg-1"
	replayed := false
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["messageId"] != "msg-1" {
			t.Errorf("expected messageId=msg-1")
		}
		if v, ok := req.Variables["replayed"]; !ok || v != false {
			t.Errorf("expected replayed=false, got %v", v)
		}
		return map[string]any{
			"outboundDlqEntries": map[string]any{
				"nodes":    []map[string]any{{"id": "odlq-1", "messageId": "msg-1", "lastError": "boom"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	es, _, err := client.OutboundDLQ.List(context.Background(), &ListOutboundDLQOptions{
		MessageID: &msgID,
		Replayed:  &replayed,
	})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(es) != 1 {
		t.Errorf("got %d", len(es))
	}
}

func TestOutboundDLQReplay(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"replayOutboundDlqEntry": true}, nil
	})
	defer srv.Close()
	if err := client.OutboundDLQ.Replay(context.Background(), "odlq-1"); err != nil {
		t.Fatalf("Replay() error: %v", err)
	}
}

func TestOutboundDLQReplayAll(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"replayAllOutboundDlq": map[string]any{"deliveries": 7}}, nil
	})
	defer srv.Close()
	r, err := client.OutboundDLQ.ReplayAll(context.Background())
	if err != nil {
		t.Fatalf("ReplayAll() error: %v", err)
	}
	if r.Deliveries != 7 {
		t.Errorf("Deliveries = %d", r.Deliveries)
	}
}

func TestOutboundDLQPurge(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["olderThan"] != "24h" {
			t.Errorf("expected olderThan=24h, got %v", req.Variables["olderThan"])
		}
		return map[string]any{"purgeOutboundDlq": map[string]any{"purged": 42}}, nil
	})
	defer srv.Close()
	r, err := client.OutboundDLQ.Purge(context.Background(), "24h")
	if err != nil {
		t.Fatalf("Purge() error: %v", err)
	}
	if r.Purged != 42 {
		t.Errorf("Purged = %d", r.Purged)
	}
}
