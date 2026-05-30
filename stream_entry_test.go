package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestStreamEntries(t *testing.T) {
	t.Parallel()
	limit := 25
	after := 100
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "streamEntries") {
			t.Errorf("expected streamEntries query")
		}
		if v, ok := req.Variables["streamId"]; !ok || v != "stream-1" {
			t.Errorf("streamId = %v", v)
		}
		if v, ok := req.Variables["limit"]; !ok || v.(float64) != 25 {
			t.Errorf("limit = %v", v)
		}
		if v, ok := req.Variables["afterSequence"]; !ok || v.(float64) != 100 {
			t.Errorf("afterSequence = %v", v)
		}
		return map[string]any{
			"streamEntries": map[string]any{
				"nodes": []map[string]any{
					{"id": "se-1", "streamId": "stream-1", "sequence": 101, "eventType": "user.created", "payload": "eyJrIjoidiJ9"},
					{"id": "se-2", "streamId": "stream-1", "sequence": 102, "eventType": "user.updated", "payload": "eyJrIjoidiJ9"},
				},
				"pageInfo": map[string]any{"total": 2, "limit": 25, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	entries, _, err := client.Streams.Entries(context.Background(), "stream-1", &ListStreamEntriesOptions{
		AfterSequence: &after,
		Limit:         &limit,
	})
	if err != nil {
		t.Fatalf("Entries() error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries", len(entries))
	}
	if entries[0].Sequence != 101 {
		t.Errorf("entries[0].Sequence = %d", entries[0].Sequence)
	}
}
