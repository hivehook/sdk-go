package hivehook

import (
	"context"
	"strings"
	"testing"
)

func TestAlertRuleList(t *testing.T) {
	t.Parallel()
	enabled := true
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if v, ok := req.Variables["enabled"]; !ok || v != true {
			t.Errorf("expected enabled=true")
		}
		return map[string]any{
			"alertRules": map[string]any{
				"nodes": []map[string]any{
					{"id": "ar-1", "name": "DLQ Spike", "conditionType": "dlq_size", "threshold": 10, "enabled": true},
				},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	rules, _, err := client.AlertRules.List(context.Background(), &ListAlertRulesOptions{Enabled: &enabled})
	if err != nil {
		t.Fatalf("AlertRules.List() error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("got %d", len(rules))
	}
}

func TestAlertRuleGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"alertRule": map[string]any{"id": "ar-1", "name": "Spike"},
		}, nil
	})
	defer srv.Close()
	r, err := client.AlertRules.Get(context.Background(), "ar-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if r.ID != "ar-1" {
		t.Errorf("r.ID = %q", r.ID)
	}
}

func TestAlertRuleCreate(t *testing.T) {
	t.Parallel()
	channel := AlertChannelSlack
	enabled := true
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if !strings.Contains(req.Query, "createAlertRule") {
			t.Errorf("expected createAlertRule mutation")
		}
		return map[string]any{
			"createAlertRule": map[string]any{
				"id": "ar-new", "name": "New Rule", "conditionType": "dlq_size", "threshold": 5, "enabled": true,
			},
		}, nil
	})
	defer srv.Close()
	r, err := client.AlertRules.Create(context.Background(), &CreateAlertRuleInput{
		Name:          "New Rule",
		ConditionType: "dlq_size",
		Threshold:     5,
		Channel:       &channel,
		SlackConfig:   &SlackAlertConfig{WebhookURL: "https://hooks", Channel: "#alerts"},
		EmailConfig:   &EmailAlertConfig{To: []string{"a@b.com"}, SubjectTemplate: "subj"},
		Enabled:       &enabled,
	})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if r.Name != "New Rule" {
		t.Errorf("r.Name = %q", r.Name)
	}
}

func TestAlertRuleUpdate(t *testing.T) {
	t.Parallel()
	name := "Renamed"
	threshold := 10
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateAlertRule": map[string]any{"id": "ar-1", "name": "Renamed", "threshold": 10},
		}, nil
	})
	defer srv.Close()
	r, err := client.AlertRules.Update(context.Background(), "ar-1", &UpdateAlertRuleInput{
		Name:      &name,
		Threshold: &threshold,
	})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if r.Name != "Renamed" {
		t.Errorf("r.Name = %q", r.Name)
	}
}

func TestAlertRuleDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteAlertRule": true}, nil
	})
	defer srv.Close()
	if err := client.AlertRules.Delete(context.Background(), "ar-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestAlertRuleTest(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"testAlertRule": true}, nil
	})
	defer srv.Close()
	if err := client.AlertRules.Test(context.Background(), "ar-1"); err != nil {
		t.Fatalf("Test() error: %v", err)
	}
}
