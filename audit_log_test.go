package hivehook

import (
	"context"
	"testing"
	"time"
)

func TestAuditLogList(t *testing.T) {
	t.Parallel()
	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["actorType"] != "user" {
			t.Errorf("expected actorType=user")
		}
		if req.Variables["resourceType"] != "destination" {
			t.Errorf("expected resourceType=destination")
		}
		if req.Variables["action"] != "delete" {
			t.Errorf("expected action=delete")
		}
		return map[string]any{
			"auditLogs": map[string]any{
				"nodes": []map[string]any{
					{"id": "al-1", "actorType": "user", "action": "delete", "resourceType": "destination"},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	logs, _, err := client.AuditLogs.List(context.Background(), &ListAuditLogsOptions{
		ActorType:    "user",
		ResourceType: "destination",
		ResourceID:   "dst-1",
		Action:       "delete",
		Since:        &since,
		Until:        &until,
	})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("got %d", len(logs))
	}
}

func TestAuditLogGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"auditLog": map[string]any{"id": "al-1", "action": "create"},
		}, nil
	})
	defer srv.Close()
	l, err := client.AuditLogs.Get(context.Background(), "al-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if l.ID != "al-1" {
		t.Errorf("l.ID = %q", l.ID)
	}
}
