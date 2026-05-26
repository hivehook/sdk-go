package hivehook

import "context"

const userFragment = `
	id organizationId email name role lastLoginAt createdAt updatedAt
`

const listUsersQuery = `query($organizationId: UUID, $search: String, $limit: Int, $offset: Int) {
	users(organizationId: $organizationId, search: $search, limit: $limit, offset: $offset) {
		nodes {` + userFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const meQuery = `query {
	me {` + userFragment + `}
}`

const inviteUserMutation = `mutation($organizationId: UUID!, $input: InviteUserInput!) {
	inviteUser(organizationId: $organizationId, input: $input) {` + userFragment + `}
}`

const removeUserMutation = `mutation($id: UUID!) {
	removeUser(id: $id)
}`

const updateUserRoleMutation = `mutation($id: UUID!, $input: UpdateUserRoleInput!) {
	updateUserRole(id: $id, input: $input) {` + userFragment + `}
}`

// UserService manages User accounts.
type UserService struct {
	gql *graphqlClient
}

// ListUsersOptions filters list results.
type ListUsersOptions struct {
	ListOptions
	OrganizationID string
}

// InviteUserInput is the input for the matching service method.
type InviteUserInput struct {
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
	Role  *string `json:"role,omitempty"`
}

// UpdateUserRoleInput is the input for the matching Update method.
type UpdateUserRoleInput struct {
	Role string `json:"role"`
}

// List returns Users matching opts and a PageInfo cursor.
func (s *UserService) List(ctx context.Context, opts *ListUsersOptions) ([]*User, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.OrganizationID != "" {
			vars["organizationId"] = opts.OrganizationID
		}
	}
	var result struct {
		Users struct {
			Nodes    []*User  `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"users"`
	}
	if err := s.gql.do(ctx, listUsersQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Users.Nodes, &result.Users.PageInfo, nil
}

// Me returns the current authenticated User.
func (s *UserService) Me(ctx context.Context) (*User, error) {
	var result struct {
		Me *User `json:"me"`
	}
	if err := s.gql.do(ctx, meQuery, nil, &result); err != nil {
		return nil, err
	}
	return result.Me, nil
}

// Invite invites a new User to an Organization.
func (s *UserService) Invite(ctx context.Context, organizationID string, input *InviteUserInput) (*User, error) {
	inp := map[string]any{
		"email": input.Email,
	}
	if input.Name != nil {
		inp["name"] = *input.Name
	}
	if input.Role != nil {
		inp["role"] = *input.Role
	}
	vars := map[string]any{"organizationId": organizationID, "input": inp}
	var result struct {
		InviteUser *User `json:"inviteUser"`
	}
	if err := s.gql.do(ctx, inviteUserMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.InviteUser, nil
}

// Remove removes a User by ID.
func (s *UserService) Remove(ctx context.Context, id string) error {
	return s.gql.do(ctx, removeUserMutation, map[string]any{"id": id}, nil)
}

// UpdateRole changes a User's role within their Organization.
func (s *UserService) UpdateRole(ctx context.Context, id string, input *UpdateUserRoleInput) (*User, error) {
	inp := map[string]any{
		"role": input.Role,
	}
	vars := map[string]any{"id": id, "input": inp}
	var result struct {
		UpdateUserRole *User `json:"updateUserRole"`
	}
	if err := s.gql.do(ctx, updateUserRoleMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.UpdateUserRole, nil
}
