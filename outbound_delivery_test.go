package hivehook

import (
	"context"
	"testing"
)

func TestOutboundDeliveryList(t *testing.T) {
	t.Parallel()
	msgID := "msg-1"
	epID := "ep-1"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["messageId"] != "msg-1" {
			t.Errorf("expected messageId=msg-1")
		}
		if req.Variables["endpointId"] != "ep-1" {
			t.Errorf("expected endpointId=ep-1")
		}
		return map[string]any{
			"outboundDeliveries": map[string]any{
				"nodes":    []map[string]any{{"id": "od-1", "messageId": "msg-1", "endpointId": "ep-1"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	ds, _, err := client.OutboundDeliveries.List(context.Background(), &ListOutboundDeliveriesOptions{
		MessageID:  &msgID,
		EndpointID: &epID,
		Status:     DeliveryStatusDelivered,
	})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(ds) != 1 {
		t.Errorf("got %d", len(ds))
	}
}

func TestOutboundDeliveryGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"outboundDelivery": map[string]any{
				"id": "od-1", "status": "DELIVERED",
				"deliveryAttempts": []map[string]any{
					{"id": "att-1", "attemptNumber": 1, "responseStatus": 200},
				},
			},
		}, nil
	})
	defer srv.Close()
	d, err := client.OutboundDeliveries.Get(context.Background(), "od-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if d.ID != "od-1" {
		t.Errorf("d.ID = %q", d.ID)
	}
	if len(d.DeliveryAttempts) != 1 {
		t.Errorf("got %d delivery attempts", len(d.DeliveryAttempts))
	}
}
