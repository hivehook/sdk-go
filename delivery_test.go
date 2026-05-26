package hivehook

import (
	"context"
	"testing"
)

func TestDeliveryList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["status"] != "PENDING" {
			t.Errorf("expected status=PENDING, got %v", req.Variables["status"])
		}
		return map[string]any{
			"deliveries": map[string]any{
				"nodes": []map[string]any{
					{"id": "del-1", "eventId": "ev-1", "status": "PENDING", "attempts": 1, "maxAttempts": 3},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	deliveries, _, err := client.Deliveries.List(context.Background(), &ListDeliveriesOptions{
		Status: DeliveryStatusPending,
	})
	if err != nil {
		t.Fatalf("Deliveries.List() error: %v", err)
	}
	if len(deliveries) != 1 {
		t.Errorf("got %d", len(deliveries))
	}
}

func TestDeliveryGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"delivery": map[string]any{
				"id": "del-1", "status": "DELIVERED", "attempts": 1,
				"deliveryAttempts": []map[string]any{
					{"id": "att-1", "deliveryId": "del-1", "attemptNumber": 1, "responseStatus": 200},
				},
			},
		}, nil
	})
	defer srv.Close()
	del, err := client.Deliveries.Get(context.Background(), "del-1")
	if err != nil {
		t.Fatalf("Deliveries.Get() error: %v", err)
	}
	if del.ID != "del-1" {
		t.Errorf("del.ID = %q", del.ID)
	}
	if len(del.DeliveryAttempts) != 1 {
		t.Errorf("got %d delivery attempts", len(del.DeliveryAttempts))
	}
}
