package hivehook

import (
	"context"
	"reflect"
	"testing"
)

func TestSourceIterate(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		after, _ := req.Variables["after"].(string)
		if after == "" {
			return map[string]any{"sources": map[string]any{
				"nodes":    []map[string]any{{"id": "s1"}, {"id": "s2"}},
				"pageInfo": map[string]any{"endCursor": "c1", "hasNextPage": true},
			}}, nil
		}
		return map[string]any{"sources": map[string]any{
			"nodes":    []map[string]any{{"id": "s3"}},
			"pageInfo": map[string]any{"hasNextPage": false},
		}}, nil
	})
	defer srv.Close()

	var ids []string
	err := client.Sources.Iterate(context.Background(), nil, func(s *Source) bool {
		ids = append(ids, s.ID)
		return true
	})
	if err != nil {
		t.Fatalf("Iterate() error: %v", err)
	}
	if want := []string{"s1", "s2", "s3"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("ids = %v, want %v", ids, want)
	}
}

func TestIterateEarlyStop(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"sources": map[string]any{
			"nodes":    []map[string]any{{"id": "a"}, {"id": "b"}, {"id": "c"}},
			"pageInfo": map[string]any{"endCursor": "next", "hasNextPage": true},
		}}, nil
	})
	defer srv.Close()

	var ids []string
	err := client.Sources.Iterate(context.Background(), nil, func(s *Source) bool {
		ids = append(ids, s.ID)
		return s.ID != "b"
	})
	if err != nil {
		t.Fatalf("Iterate() error: %v", err)
	}
	if want := []string{"a", "b"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("ids = %v, want %v (should stop when yield returns false)", ids, want)
	}
}
