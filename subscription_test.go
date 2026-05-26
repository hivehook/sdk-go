package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestSubscriptionList(t *testing.T) {
	t.Parallel()
	enabled := true
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if v, ok := req.Variables["enabled"]; !ok || v != true {
			t.Errorf("expected enabled=true var, got %v", req.Variables)
		}
		if req.Variables["sourceId"] != "src-1" {
			t.Errorf("expected sourceId=src-1")
		}
		return map[string]any{
			"subscriptions": map[string]any{
				"nodes": []map[string]any{
					{"id": "sub-1", "name": "S1", "sourceId": "src-1", "destinationId": "dst-1", "enabled": true},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	subs, _, err := client.Subscriptions.List(context.Background(), &ListSubscriptionsOptions{
		SourceID: "src-1",
		Enabled:  &enabled,
	})
	if err != nil {
		t.Fatalf("Subscriptions.List() error: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("got %d, want 1", len(subs))
	}
}

func TestSubscriptionGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"subscription": map[string]any{"id": "sub-1", "name": "S1"},
		}, nil
	})
	defer srv.Close()
	sub, err := client.Subscriptions.Get(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("Subscriptions.Get() error: %v", err)
	}
	if sub.ID != "sub-1" {
		t.Errorf("sub.ID = %q", sub.ID)
	}
}

func TestSubscriptionCreate(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createSubscription") {
			t.Errorf("expected createSubscription mutation")
		}
		return map[string]any{
			"createSubscription": map[string]any{
				"id": "sub-new", "name": "New Sub", "sourceId": "src-1", "destinationId": "dst-1",
			},
		}, nil
	})
	defer srv.Close()

	sub, err := client.Subscriptions.Create(context.Background(), &CreateSubscriptionInput{
		Name:          "New Sub",
		SourceID:      "src-1",
		DestinationID: "dst-1",
		FilterConfig: &FilterConfig{
			EventTypes: []string{"user.created"},
			BodyMatch:  []BodyMatchRule{{Path: "$.kind", Value: "x", Operator: "eq"}},
			Rules: []FilterRule{
				{Operator: "and", Rules: []FilterRule{{Path: "$.a", Operator: "eq", Value: "1"}}},
			},
		},
		TransformConfig: &TransformConfig{Envelope: true, Headers: map[string]any{"X": "Y"}},
	})
	if err != nil {
		t.Fatalf("Subscriptions.Create() error: %v", err)
	}
	if sub.Name != "New Sub" {
		t.Errorf("sub.Name = %q", sub.Name)
	}
}

func TestSubscriptionUpdate(t *testing.T) {
	t.Parallel()
	name := "Renamed"
	enabled := false
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateSubscription": map[string]any{"id": "sub-1", "name": "Renamed", "enabled": false},
		}, nil
	})
	defer srv.Close()
	sub, err := client.Subscriptions.Update(context.Background(), "sub-1", &UpdateSubscriptionInput{
		Name:    &name,
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("Subscriptions.Update() error: %v", err)
	}
	if sub.Name != "Renamed" {
		t.Errorf("sub.Name = %q", sub.Name)
	}
}

func TestSubscriptionDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteSubscription": true}, nil
	})
	defer srv.Close()
	if err := client.Subscriptions.Delete(context.Background(), "sub-1"); err != nil {
		t.Fatalf("Subscriptions.Delete() error: %v", err)
	}
}
