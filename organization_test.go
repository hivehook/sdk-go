package hivehook

import (
	"context"
	"testing"
)

func TestOrganizationList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"organizations": map[string]any{
				"nodes":    []map[string]any{{"id": "org-1", "name": "Acme", "slug": "acme"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	os, _, err := client.Organizations.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(os) != 1 {
		t.Errorf("got %d", len(os))
	}
}

func TestOrganizationGet(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"organization": map[string]any{"id": "org-1", "name": "Acme"}}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.Get(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if o.ID != "org-1" {
		t.Errorf("o.ID = %q", o.ID)
	}
}

func TestOrganizationCreate(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"createOrganization": map[string]any{"id": "org-new", "name": "Acme", "slug": "acme"},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.Create(context.Background(), &CreateOrganizationInput{Name: "Acme", Slug: "acme"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if o.Slug != "acme" {
		t.Errorf("o.Slug = %q", o.Slug)
	}
}

func TestOrganizationUpdate(t *testing.T) {
	t.Parallel()
	name := "Renamed"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateOrganization": map[string]any{"id": "org-1", "name": "Renamed"},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.Update(context.Background(), "org-1", &UpdateOrganizationInput{Name: &name})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if o.Name != "Renamed" {
		t.Errorf("o.Name = %q", o.Name)
	}
}

func TestOrganizationDelete(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteOrganization": true}, nil
	})
	defer srv.Close()
	if err := client.Organizations.Delete(context.Background(), "org-1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
}

func TestOrganizationConfigureSSO(t *testing.T) {
	t.Parallel()
	idp := "https://idp.example.com/metadata"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"configureSSO": map[string]any{"id": "org-1", "ssoEnabled": true, "ssoProvider": "saml"},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.ConfigureSSO(context.Background(), "org-1", &SSOConfigInput{
		Provider:       "saml",
		IDPMetadataURL: &idp,
	})
	if err != nil {
		t.Fatalf("ConfigureSSO() error: %v", err)
	}
	if !o.SSOEnabled {
		t.Errorf("SSOEnabled = false")
	}
}

func TestOrganizationDisableSSO(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"disableSSO": map[string]any{"id": "org-1", "ssoEnabled": false},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.DisableSSO(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("DisableSSO() error: %v", err)
	}
	if o.SSOEnabled {
		t.Errorf("SSOEnabled = true")
	}
}

func TestOrganizationUpdateRetention(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateOrganizationRetention": map[string]any{
				"id": "org-1", "retentionEvents": 60, "retentionMessages": 90,
			},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.UpdateRetention(context.Background(), "org-1", &RetentionInput{
		RetentionEvents:   60,
		RetentionMessages: 90,
	})
	if err != nil {
		t.Fatalf("UpdateRetention() error: %v", err)
	}
	if o.RetentionEvents != 60 {
		t.Errorf("RetentionEvents = %d", o.RetentionEvents)
	}
}

func TestOrganizationDeleteData(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"deleteOrganizationData": true}, nil
	})
	defer srv.Close()
	if err := client.Organizations.DeleteData(context.Background(), "org-1"); err != nil {
		t.Fatalf("DeleteData() error: %v", err)
	}
}

func TestOrganizationExportData(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"exportOrganizationData": map[string]any{"foo": "bar"}}, nil
	})
	defer srv.Close()
	data, err := client.Organizations.ExportData(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("ExportData() error: %v", err)
	}
	if data["foo"] != "bar" {
		t.Errorf("data[foo] = %v", data["foo"])
	}
}

func TestOrganizationConfigureOTLP(t *testing.T) {
	t.Parallel()
	insecure := true
	rate := 0.5
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"configureOTLP": map[string]any{
				"id": "org-1",
				"otlpConfig": map[string]any{
					"endpoint": "https://otlp.example.com", "insecure": true, "sampleRate": 0.5,
				},
			},
		}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.ConfigureOTLP(context.Background(), "org-1", &OTLPConfigInput{
		Endpoint:   "https://otlp.example.com",
		Headers:    map[string]any{"X-Api-Key": "abc"},
		Insecure:   &insecure,
		SampleRate: &rate,
	})
	if err != nil {
		t.Fatalf("ConfigureOTLP() error: %v", err)
	}
	if o.OTLPConfig == nil || o.OTLPConfig.Endpoint != "https://otlp.example.com" {
		t.Errorf("OTLPConfig = %+v", o.OTLPConfig)
	}
}

func TestOrganizationDisableOTLP(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"disableOTLP": map[string]any{"id": "org-1"}}, nil
	})
	defer srv.Close()
	o, err := client.Organizations.DisableOTLP(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("DisableOTLP() error: %v", err)
	}
	if o.ID != "org-1" {
		t.Errorf("o.ID = %q", o.ID)
	}
}
