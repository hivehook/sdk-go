package hivehook

import (
	"context"
	"testing"
)

func TestTransformationList(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"transformations": map[string]any{
				"nodes": []map[string]any{
					{"id": "tf-1", "name": "Enrich", "code": "return event", "enabled": true, "failOpen": false, "timeoutMs": 1000},
					{"id": "tf-2", "name": "Filter", "code": "return null", "enabled": false, "failOpen": true, "timeoutMs": 500},
				},
				"pageInfo": map[string]any{"total": 2, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	transformations, pi, err := client.Transformations.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Transformations.List() error: %v", err)
	}
	if len(transformations) != 2 {
		t.Fatalf("got %d transformations, want 2", len(transformations))
	}
	if transformations[0].Name != "Enrich" {
		t.Errorf("transformations[0].Name = %q, want %q", transformations[0].Name, "Enrich")
	}
	if pi.Total != 2 {
		t.Errorf("pageInfo.Total = %d, want 2", pi.Total)
	}
}

func TestTransformationListWithOptions(t *testing.T) {
	enabled := true
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		vars := req.Variables
		if v, ok := vars["enabled"]; !ok || v != true {
			t.Errorf("expected enabled=true in variables, got %v", vars)
		}
		return map[string]any{
			"transformations": map[string]any{
				"nodes":    []map[string]any{{"id": "tf-1", "name": "Enrich", "code": "return event", "enabled": true}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()

	transformations, _, err := client.Transformations.List(context.Background(), &ListTransformationsOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("Transformations.List() error: %v", err)
	}
	if len(transformations) != 1 {
		t.Fatalf("got %d transformations, want 1", len(transformations))
	}
}

func TestTransformationGet(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"transformation": map[string]any{
				"id": "tf-1", "name": "Enrich", "code": "return event", "enabled": true, "failOpen": false, "timeoutMs": 1000,
			},
		}, nil
	})
	defer srv.Close()

	tf, err := client.Transformations.Get(context.Background(), "tf-1")
	if err != nil {
		t.Fatalf("Transformations.Get() error: %v", err)
	}
	if tf.ID != "tf-1" {
		t.Errorf("tf.ID = %q, want %q", tf.ID, "tf-1")
	}
	if tf.Name != "Enrich" {
		t.Errorf("tf.Name = %q, want %q", tf.Name, "Enrich")
	}
}

func TestTransformationCreate(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createTransformation": map[string]any{
				"id": "tf-new", "name": "New Transform", "code": "return event", "enabled": true, "failOpen": false, "timeoutMs": 1000,
			},
		}, nil
	})
	defer srv.Close()

	tf, err := client.Transformations.Create(context.Background(), &CreateTransformationInput{
		Name: "New Transform",
		Code: "return event",
	})
	if err != nil {
		t.Fatalf("Transformations.Create() error: %v", err)
	}
	if tf.Name != "New Transform" {
		t.Errorf("tf.Name = %q, want %q", tf.Name, "New Transform")
	}
}

func TestTransformationUpdate(t *testing.T) {
	name := "Updated Transform"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateTransformation": map[string]any{
				"id": "tf-1", "name": "Updated Transform", "code": "return event", "enabled": true, "failOpen": false, "timeoutMs": 1000,
			},
		}, nil
	})
	defer srv.Close()

	tf, err := client.Transformations.Update(context.Background(), "tf-1", &UpdateTransformationInput{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("Transformations.Update() error: %v", err)
	}
	if tf.Name != "Updated Transform" {
		t.Errorf("tf.Name = %q, want %q", tf.Name, "Updated Transform")
	}
}

func TestTransformationDelete(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteTransformation": true}, nil
	})
	defer srv.Close()

	err := client.Transformations.Delete(context.Background(), "tf-1")
	if err != nil {
		t.Fatalf("Transformations.Delete() error: %v", err)
	}
}

func TestTransformationTest(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"testTransformation": map[string]any{
				"success":    true,
				"output":     map[string]any{"user": "test", "enriched": true},
				"error":      "",
				"durationMs": 12,
			},
		}, nil
	})
	defer srv.Close()

	result, err := client.Transformations.Test(context.Background(), &TestTransformationInput{
		Code:      "event.enriched = true; return event",
		Payload:   map[string]any{"user": "test"},
		EventType: "user.created",
	})
	if err != nil {
		t.Fatalf("Transformations.Test() error: %v", err)
	}
	if !result.Success {
		t.Error("result.Success = false, want true")
	}
	if result.DurationMs != 12 {
		t.Errorf("result.DurationMs = %d, want 12", result.DurationMs)
	}
	if result.Output["enriched"] != true {
		t.Errorf("result.Output[enriched] = %v, want true", result.Output["enriched"])
	}
}

func TestTransformationTestError(t *testing.T) {
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"testTransformation": map[string]any{
				"success":    false,
				"error":      "ReferenceError: x is not defined",
				"durationMs": 3,
			},
		}, nil
	})
	defer srv.Close()

	result, err := client.Transformations.Test(context.Background(), &TestTransformationInput{
		Code:      "return x",
		Payload:   map[string]any{},
		EventType: "user.created",
	})
	if err != nil {
		t.Fatalf("Transformations.Test() error: %v", err)
	}
	if result.Success {
		t.Error("result.Success = true, want false")
	}
	if result.Error != "ReferenceError: x is not defined" {
		t.Errorf("result.Error = %q, want %q", result.Error, "ReferenceError: x is not defined")
	}
}
