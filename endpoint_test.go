package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestEndpointList(t *testing.T) {
	t.Parallel()
	appID := "app-1"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["applicationId"] != "app-1" {
			t.Errorf("expected applicationId=app-1, got %v", req.Variables["applicationId"])
		}
		return map[string]any{
			"endpoints": map[string]any{
				"nodes":    []map[string]any{{"id": "ep-1", "applicationId": "app-1", "url": "https://e.test"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	eps, _, err := client.Endpoints.List(context.Background(), &ListEndpointsOptions{
		ApplicationID: &appID,
		Status:        EndpointStatusActive,
	})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(eps) != 1 {
		t.Errorf("got %d", len(eps))
	}
}

func TestEndpointGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"endpoint": map[string]any{"id": "ep-1", "url": "https://e.test"},
		}, nil
	})
	defer srv.Close()
	ep, err := client.Endpoints.Get(context.Background(), "ep-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if ep.ID != "ep-1" {
		t.Errorf("ep.ID = %q", ep.ID)
	}
}

func TestEndpointCreate(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createEndpoint") {
			t.Errorf("expected createEndpoint mutation")
		}
		return map[string]any{
			"createEndpoint": map[string]any{"id": "ep-new", "url": "https://x.test", "applicationId": "app-1"},
		}, nil
	})
	defer srv.Close()
	timeout := 5000
	ep, err := client.Endpoints.Create(context.Background(), &CreateEndpointInput{
		ApplicationID: "app-1",
		URL:           "https://x.test",
		TimeoutMs:     &timeout,
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if ep.ID != "ep-new" {
		t.Errorf("ep.ID = %q", ep.ID)
	}
}

func TestEndpointUpdate(t *testing.T) {
	t.Parallel()
	url := "https://updated.test"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateEndpoint": map[string]any{"id": "ep-1", "url": "https://updated.test"},
		}, nil
	})
	defer srv.Close()
	ep, err := client.Endpoints.Update(context.Background(), "ep-1", &UpdateEndpointInput{
		URL:    &url,
		Status: EndpointStatusInactive,
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if ep.URL != "https://updated.test" {
		t.Errorf("URL = %q", ep.URL)
	}
}

func TestEndpointDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteEndpoint": true}, nil
	})
	defer srv.Close()
	if err := client.Endpoints.Delete(context.Background(), "ep-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestEndpointRotateSecret(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"rotateEndpointSecret": map[string]any{"id": "ep-1", "signingSecret": "new"},
		}, nil
	})
	defer srv.Close()
	ep, err := client.Endpoints.RotateSecret(context.Background(), "ep-1")
	if err != nil {
		t.Fatalf("RotateSecret() error: %v", err)
	}
	if ep.SigningSecret != "new" {
		t.Errorf("SigningSecret = %q", ep.SigningSecret)
	}
}

func TestEndpointPollDeliveries(t *testing.T) {
	t.Parallel()
	limit := 5
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["endpointId"] != "ep-1" {
			t.Errorf("expected endpointId=ep-1")
		}
		return map[string]any{
			"pollOutboundDeliveries": map[string]any{
				"nodes":    []map[string]any{{"id": "od-1", "messageId": "msg-1"}},
				"pageInfo": map[string]any{"total": 1, "limit": 5, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	ds, _, err := client.Endpoints.PollDeliveries(context.Background(), "ep-1", nil, &limit)
	if err != nil {
		t.Fatalf("PollDeliveries() error: %v", err)
	}
	if len(ds) != 1 {
		t.Errorf("got %d", len(ds))
	}
}

func TestEndpointAckDeliveries(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"ackOutboundDeliveries": 2}, nil
	})
	defer srv.Close()
	n, err := client.Endpoints.AckDeliveries(context.Background(), "ep-1", []string{"a", "b"})
	if err != nil {
		t.Fatalf("AckDeliveries() error: %v", err)
	}
	if n != 2 {
		t.Errorf("n = %d", n)
	}
}

func TestEndpointRegeneratePollAPIKey(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"regenerateOutboundPollApiKey": map[string]any{"id": "ep-1", "pollApiKey": "k"},
		}, nil
	})
	defer srv.Close()
	ep, err := client.Endpoints.RegeneratePollAPIKey(context.Background(), "ep-1")
	if err != nil {
		t.Fatalf("RegeneratePollAPIKey() error: %v", err)
	}
	if ep.PollAPIKey != "k" {
		t.Errorf("PollAPIKey = %q", ep.PollAPIKey)
	}
}

func TestEndpointSkipOutboundDLQEntry(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"skipOutboundDlqEntry": true}, nil
	})
	defer srv.Close()
	if err := client.Endpoints.SkipOutboundDLQEntry(context.Background(), "dlq-1"); err != nil {
		t.Fatalf("SkipOutboundDLQEntry() error: %v", err)
	}
}
