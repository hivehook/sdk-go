package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestMetaEventConfigList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "metaEventConfigs") {
			t.Errorf("expected metaEventConfigs query")
		}
		return map[string]any{
			"metaEventConfigs": map[string]any{
				"nodes": []map[string]any{
					{"id": "me-1", "name": "DLQ alerts", "url": "https://hooks.example/alerts", "eventTypes": []string{"delivery.failed"}, "enabled": true},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	configs, page, err := client.MetaEventConfigs.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("got %d configs", len(configs))
	}
	if configs[0].Name != "DLQ alerts" {
		t.Errorf("Name = %q", configs[0].Name)
	}
	if page.Total != 1 {
		t.Errorf("page.Total = %d", page.Total)
	}
}

func TestMetaEventConfigGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if v, ok := req.Variables["id"]; !ok || v != "me-7" {
			t.Errorf("expected id=me-7, got %v", v)
		}
		return map[string]any{
			"metaEventConfig": map[string]any{"id": "me-7", "name": "x", "url": "https://x", "enabled": true},
		}, nil
	})
	defer srv.Close()
	c, err := client.MetaEventConfigs.Get(context.Background(), "me-7")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if c.ID != "me-7" {
		t.Errorf("ID = %q", c.ID)
	}
}

func TestMetaEventConfigCreate(t *testing.T) {
	t.Parallel()
	enabled := true
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createMetaEventConfig") {
			t.Errorf("expected createMetaEventConfig mutation")
		}
		return map[string]any{
			"createMetaEventConfig": map[string]any{
				"id": "me-new", "name": "new", "url": "https://new", "eventTypes": []string{"source.created"}, "enabled": true,
			},
		}, nil
	})
	defer srv.Close()
	c, err := client.MetaEventConfigs.Create(context.Background(), &CreateMetaEventConfigInput{
		Name:       "new",
		URL:        "https://new",
		EventTypes: []string{"source.created"},
		Enabled:    &enabled,
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if c.ID != "me-new" {
		t.Errorf("ID = %q", c.ID)
	}
}

func TestMetaEventConfigUpdate(t *testing.T) {
	t.Parallel()
	name := "renamed"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "updateMetaEventConfig") {
			t.Errorf("expected updateMetaEventConfig mutation")
		}
		return map[string]any{
			"updateMetaEventConfig": map[string]any{"id": "me-1", "name": "renamed"},
		}, nil
	})
	defer srv.Close()
	c, err := client.MetaEventConfigs.Update(context.Background(), "me-1", &UpdateMetaEventConfigInput{Name: &name})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if c.Name != "renamed" {
		t.Errorf("Name = %q", c.Name)
	}
}

func TestMetaEventConfigDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "deleteMetaEventConfig") {
			t.Errorf("expected deleteMetaEventConfig mutation")
		}
		return map[string]any{"deleteMetaEventConfig": true}, nil
	})
	defer srv.Close()
	if err := client.MetaEventConfigs.Delete(context.Background(), "me-1"); err != nil {
		t.Errorf("Delete() error: %v", err)
	}
}
