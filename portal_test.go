package hivehook

import (
	"context"
	"testing"
)

func TestPortalGenerateToken(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["applicationId"] != "app-1" {
			t.Errorf("expected applicationId=app-1")
		}
		return map[string]any{
			"generatePortalToken": map[string]any{
				"token":     "portal-token-xyz",
				"expiresAt": "2030-01-01T00:00:00Z",
			},
		}, nil
	})
	defer srv.Close()
	tok, err := client.Portal.GenerateToken(context.Background(), "app-1")
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}
	if tok.Token != "portal-token-xyz" {
		t.Errorf("Token = %q", tok.Token)
	}
}
