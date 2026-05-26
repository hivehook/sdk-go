package hivehook

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestAPIKeyList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"apiKeys": map[string]any{
				"nodes": []map[string]any{
					{"id": "ak-1", "name": "key1", "keyPrefix": "htk_a"},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	keys, _, err := client.APIKeys.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("APIKeys.List() error: %v", err)
	}
	if len(keys) != 1 || keys[0].Name != "key1" {
		t.Errorf("unexpected: %+v", keys)
	}
}

func TestAPIKeyGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"apiKey": map[string]any{"id": "ak-1", "name": "key1"},
		}, nil
	})
	defer srv.Close()
	k, err := client.APIKeys.Get(context.Background(), "ak-1")
	if err != nil {
		t.Fatalf("APIKeys.Get() error: %v", err)
	}
	if k.ID != "ak-1" {
		t.Errorf("k.ID = %q", k.ID)
	}
}

func TestAPIKeyCreate(t *testing.T) {
	t.Parallel()
	exp := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createAPIKey") {
			t.Errorf("expected createAPIKey mutation")
		}
		return map[string]any{
			"createAPIKey": map[string]any{
				"apiKey": map[string]any{"id": "ak-new", "name": "new-key"},
				"rawKey": "htk_secret",
			},
		}, nil
	})
	defer srv.Close()
	out, err := client.APIKeys.Create(context.Background(), &CreateAPIKeyInput{
		Name:      "new-key",
		Scopes:    []string{"read"},
		SourceIDs: []string{"src-1"},
		ExpiresAt: &exp,
	})
	if err != nil {
		t.Fatalf("APIKeys.Create() error: %v", err)
	}
	if out.RawKey != "htk_secret" {
		t.Errorf("RawKey = %q", out.RawKey)
	}
}

func TestAPIKeyRevoke(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"revokeAPIKey": true}, nil
	})
	defer srv.Close()
	if err := client.APIKeys.Revoke(context.Background(), "ak-1"); err != nil {
		t.Fatalf("APIKeys.Revoke() error: %v", err)
	}
}
