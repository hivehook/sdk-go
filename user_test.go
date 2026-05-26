package hivehook

import (
	"context"
	"testing"
)

func TestUserList(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		if req.Variables["organizationId"] != "org-1" {
			t.Errorf("expected organizationId=org-1")
		}
		return map[string]any{
			"users": map[string]any{
				"nodes":    []map[string]any{{"id": "u-1", "email": "a@b.com", "name": "A"}},
				"pageInfo": map[string]any{"total": 1, "limit": 50, "offset": 0, "hasNextPage": false},
			},
		}, nil
	})
	defer srv.Close()
	us, _, err := client.Users.List(context.Background(), &ListUsersOptions{OrganizationID: "org-1"})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(us) != 1 {
		t.Errorf("got %d", len(us))
	}
}

func TestUserMe(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"me": map[string]any{"id": "u-1", "email": "me@a.com", "name": "Me"},
		}, nil
	})
	defer srv.Close()
	u, err := client.Users.Me(context.Background())
	if err != nil {
		t.Fatalf("Me() error: %v", err)
	}
	if u.Email != "me@a.com" {
		t.Errorf("u.Email = %q", u.Email)
	}
}

func TestUserInvite(t *testing.T) {
	t.Parallel()
	name := "Bob"
	role := "admin"
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"inviteUser": map[string]any{"id": "u-new", "email": "bob@a.com", "name": "Bob", "role": "admin"},
		}, nil
	})
	defer srv.Close()
	u, err := client.Users.Invite(context.Background(), "org-1", &InviteUserInput{
		Email: "bob@a.com",
		Name:  &name,
		Role:  &role,
	})
	if err != nil {
		t.Fatalf("Invite() error: %v", err)
	}
	if u.Role != "admin" {
		t.Errorf("u.Role = %q", u.Role)
	}
}

func TestUserRemove(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{"removeUser": true}, nil
	})
	defer srv.Close()
	if err := client.Users.Remove(context.Background(), "u-1"); err != nil {
		t.Fatalf("Remove() error: %v", err)
	}
}

func TestUserUpdateRole(t *testing.T) {
	t.Parallel()
	client, srv := newTestClient(func(req graphqlRequest) (any, error) {
		return map[string]any{
			"updateUserRole": map[string]any{"id": "u-1", "role": "viewer"},
		}, nil
	})
	defer srv.Close()
	u, err := client.Users.UpdateRole(context.Background(), "u-1", &UpdateUserRoleInput{Role: "viewer"})
	if err != nil {
		t.Fatalf("UpdateRole() error: %v", err)
	}
	if u.Role != "viewer" {
		t.Errorf("u.Role = %q", u.Role)
	}
}
