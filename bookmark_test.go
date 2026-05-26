package hivehook

import (
	"context"
	"testing"
)

func TestBookmarkList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["eventId"] != "ev-1" {
			t.Errorf("expected eventId=ev-1")
		}
		return map[string]any{
			"bookmarks": map[string]any{
				"nodes":    []map[string]any{{"id": "bm-1", "eventId": "ev-1", "name": "Important"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	bms, _, err := client.Bookmarks.List(context.Background(), &ListBookmarksOptions{EventID: "ev-1"})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(bms) != 1 {
		t.Errorf("got %d", len(bms))
	}
}

func TestBookmarkGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"bookmark": map[string]any{"id": "bm-1", "name": "x"}}, nil
	})
	defer srv.Close()
	bm, err := client.Bookmarks.Get(context.Background(), "bm-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if bm.ID != "bm-1" {
		t.Errorf("bm.ID = %q", bm.ID)
	}
}

func TestBookmarkCreate(t *testing.T) {
	t.Parallel()
	name := "Important"
	notes := "Check this"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["eventId"] != "ev-1" || req.Variables["name"] != "Important" {
			t.Errorf("vars: %v", req.Variables)
		}
		return map[string]any{
			"createBookmark": map[string]any{"id": "bm-new", "eventId": "ev-1", "name": "Important", "notes": "Check this"},
		}, nil
	})
	defer srv.Close()
	bm, err := client.Bookmarks.Create(context.Background(), "ev-1", &name, &notes)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if bm.Name != "Important" {
		t.Errorf("bm.Name = %q", bm.Name)
	}
}

func TestBookmarkDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteBookmark": true}, nil
	})
	defer srv.Close()
	if err := client.Bookmarks.Delete(context.Background(), "bm-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}
