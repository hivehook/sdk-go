package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestDestinationList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "destinations(") {
			t.Errorf("expected destinations query, got: %s", req.Query)
		}
		return map[string]any{
			"destinations": map[string]any{
				"nodes": []map[string]any{
					{"id": "dst-1", "name": "Primary", "url": "https://a.test", "status": "ACTIVE"},
					{"id": "dst-2", "name": "Backup", "url": "https://b.test", "status": "INACTIVE"},
				},
				"pageInfo": map[string]any{"total": 2, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	dests, pi, err := client.Destinations.List(context.Background(), &ListDestinationsOptions{Status: DestinationStatusActive})
	if err != nil {
		t.Fatalf("Destinations.List() error: %v", err)
	}
	if len(dests) != 2 {
		t.Fatalf("got %d destinations, want 2", len(dests))
	}
	if dests[0].Name != "Primary" {
		t.Errorf("dests[0].Name = %q", dests[0].Name)
	}
	if pi.Total != 2 {
		t.Errorf("pageInfo.Total = %d", pi.Total)
	}
}

func TestDestinationGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["id"] != "dst-1" {
			t.Errorf("expected id=dst-1, got %v", req.Variables["id"])
		}
		return map[string]any{
			"destination": map[string]any{"id": "dst-1", "name": "Primary", "url": "https://a.test"},
		}, nil
	})
	defer srv.Close()

	dest, err := client.Destinations.Get(context.Background(), "dst-1")
	if err != nil {
		t.Fatalf("Destinations.Get() error: %v", err)
	}
	if dest.ID != "dst-1" {
		t.Errorf("dest.ID = %q", dest.ID)
	}
}

func TestDestinationCreate(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createDestination") {
			t.Errorf("expected createDestination mutation")
		}
		return map[string]any{
			"createDestination": map[string]any{"id": "dst-new", "name": "New", "url": "https://new.test"},
		}, nil
	})
	defer srv.Close()

	timeout := 5000
	dest, err := client.Destinations.Create(context.Background(), &CreateDestinationInput{
		Name:      "New",
		URL:       "https://new.test",
		Type:      DestinationTypeHTTP,
		TimeoutMs: &timeout,
	})
	if err != nil {
		t.Fatalf("Destinations.Create() error: %v", err)
	}
	if dest.Name != "New" {
		t.Errorf("dest.Name = %q", dest.Name)
	}
}

func TestDestinationUpdate(t *testing.T) {
	t.Parallel()
	name := "Renamed"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "updateDestination") {
			t.Errorf("expected updateDestination mutation")
		}
		return map[string]any{
			"updateDestination": map[string]any{"id": "dst-1", "name": "Renamed", "url": "https://a.test"},
		}, nil
	})
	defer srv.Close()

	dest, err := client.Destinations.Update(context.Background(), "dst-1", &UpdateDestinationInput{Name: &name})
	if err != nil {
		t.Fatalf("Destinations.Update() error: %v", err)
	}
	if dest.Name != "Renamed" {
		t.Errorf("dest.Name = %q", dest.Name)
	}
}

func TestDestinationDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteDestination": true}, nil
	})
	defer srv.Close()
	if err := client.Destinations.Delete(context.Background(), "dst-1"); err != nil {
		t.Fatalf("Destinations.Delete() error: %v", err)
	}
}

func TestDestinationRotateSecret(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"rotateDestinationSecret": map[string]any{"id": "dst-1", "signingSecret": "new-secret"},
		}, nil
	})
	defer srv.Close()
	dest, err := client.Destinations.RotateSecret(context.Background(), "dst-1")
	if err != nil {
		t.Fatalf("RotateSecret() error: %v", err)
	}
	if dest.SigningSecret != "new-secret" {
		t.Errorf("SigningSecret = %q", dest.SigningSecret)
	}
}

func TestDestinationPollDeliveries(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["destinationId"] != "dst-1" {
			t.Errorf("expected destinationId=dst-1")
		}
		return map[string]any{
			"pollDeliveries": map[string]any{
				"nodes":    []map[string]any{{"id": "del-1", "status": "PENDING"}},
				"pageInfo": map[string]any{"total": 1, "limit": 10, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	limit := 10
	deliveries, _, err := client.Destinations.PollDeliveries(context.Background(), "dst-1", nil, &limit)
	if err != nil {
		t.Fatalf("PollDeliveries() error: %v", err)
	}
	if len(deliveries) != 1 {
		t.Errorf("got %d deliveries", len(deliveries))
	}
}

func TestDestinationAckDeliveries(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"ackDeliveries": 3}, nil
	})
	defer srv.Close()
	n, err := client.Destinations.AckDeliveries(context.Background(), "dst-1", []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("AckDeliveries() error: %v", err)
	}
	if n != 3 {
		t.Errorf("n = %d, want 3", n)
	}
}

func TestDestinationRegeneratePollAPIKey(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"regeneratePollApiKey": map[string]any{"id": "dst-1", "pollApiKey": "new-key"},
		}, nil
	})
	defer srv.Close()
	dest, err := client.Destinations.RegeneratePollAPIKey(context.Background(), "dst-1")
	if err != nil {
		t.Fatalf("RegeneratePollAPIKey() error: %v", err)
	}
	if dest.PollAPIKey != "new-key" {
		t.Errorf("PollAPIKey = %q", dest.PollAPIKey)
	}
}

func TestDestinationSkipDLQEntry(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"skipDLQEntry": true}, nil
	})
	defer srv.Close()
	if err := client.Destinations.SkipDLQEntry(context.Background(), "dlq-1"); err != nil {
		t.Fatalf("SkipDLQEntry() error: %v", err)
	}
}
