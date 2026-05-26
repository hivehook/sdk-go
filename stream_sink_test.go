package hivehook

import (
	"context"
	"testing"
)

func TestStreamSinkList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["streamId"] != "stm-1" {
			t.Errorf("expected streamId=stm-1")
		}
		return map[string]any{
			"streamSinks": map[string]any{
				"nodes":    []map[string]any{{"id": "snk-1", "streamId": "stm-1", "name": "s1", "sinkType": "S3"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	ss, _, err := client.StreamSinks.List(context.Background(), "stm-1", &ListStreamSinksOptions{Status: SinkStatusActive})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(ss) != 1 {
		t.Errorf("got %d", len(ss))
	}
}

func TestStreamSinkGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"streamSink": map[string]any{"id": "snk-1", "name": "s1"}}, nil
	})
	defer srv.Close()
	s, err := client.StreamSinks.Get(context.Background(), "snk-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if s.ID != "snk-1" {
		t.Errorf("s.ID = %q", s.ID)
	}
}

func TestStreamSinkCreate(t *testing.T) {
	t.Parallel()
	batchSize := 100
	flush := "10s"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createStreamSink": map[string]any{"id": "snk-new", "name": "s1", "sinkType": "S3", "batchSize": 100, "flushInterval": "10s"},
		}, nil
	})
	defer srv.Close()
	s, err := client.StreamSinks.Create(context.Background(), &CreateStreamSinkInput{
		StreamID:      "stm-1",
		Name:          "s1",
		SinkType:      SinkTypeS3,
		Config:        map[string]any{"bucket": "my-bucket"},
		BatchSize:     &batchSize,
		FlushInterval: &flush,
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if s.BatchSize != 100 {
		t.Errorf("BatchSize = %d", s.BatchSize)
	}
}

func TestStreamSinkUpdate(t *testing.T) {
	t.Parallel()
	batchSize := 200
	status := SinkStatusPaused
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateStreamSink": map[string]any{"id": "snk-1", "name": "s1", "batchSize": 200, "status": "PAUSED"},
		}, nil
	})
	defer srv.Close()
	s, err := client.StreamSinks.Update(context.Background(), "snk-1", &UpdateStreamSinkInput{
		BatchSize: &batchSize,
		Status:    &status,
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if s.BatchSize != 200 {
		t.Errorf("BatchSize = %d", s.BatchSize)
	}
}

func TestStreamSinkDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteStreamSink": true}, nil
	})
	defer srv.Close()
	if err := client.StreamSinks.Delete(context.Background(), "snk-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}
