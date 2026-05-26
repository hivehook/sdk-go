package hivehook

import "context"

// ApplicationService manages outbound webhook Applications.
type ApplicationService struct {
	gql *graphqlClient
}

// CreateApplicationInput is the input for the matching Create method.
type CreateApplicationInput struct {
	Name string
}

// UpdateApplicationInput is the input for the matching Update method.
type UpdateApplicationInput struct {
	Name *string
}

const applicationFragment = `
	fragment ApplicationFields on Application {
		id name uid createdAt
	}
`

const listApplicationsQuery = `
	query ListApplications($search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
		applications(search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
			nodes { ...ApplicationFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + applicationFragment

const getApplicationQuery = `
	query GetApplication($id: UUID!) {
		application(id: $id) { ...ApplicationFields }
	}
` + applicationFragment

const createApplicationMutation = `
	mutation CreateApplication($input: CreateApplicationInput!) {
		createApplication(input: $input) { ...ApplicationFields }
	}
` + applicationFragment

const updateApplicationMutation = `
	mutation UpdateApplication($id: UUID!, $input: UpdateApplicationInput!) {
		updateApplication(id: $id, input: $input) { ...ApplicationFields }
	}
` + applicationFragment

const deleteApplicationMutation = `
	mutation DeleteApplication($id: UUID!) {
		deleteApplication(id: $id)
	}
`

// List returns Applications matching opts and a PageInfo cursor.
func (s *ApplicationService) List(ctx context.Context, opts *ListOptions) ([]*Application, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
	}
	var result struct {
		Applications struct {
			Nodes    []*Application `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"applications"`
	}
	if err := s.gql.do(ctx, listApplicationsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Applications.Nodes, &result.Applications.PageInfo, nil
}

// Get fetches a single Application by ID.
func (s *ApplicationService) Get(ctx context.Context, id string) (*Application, error) {
	vars := map[string]any{"id": id}
	var result struct {
		Application *Application `json:"application"`
	}
	if err := s.gql.do(ctx, getApplicationQuery, vars, &result); err != nil {
		return nil, err
	}
	return result.Application, nil
}

// Create persists a new Application.
func (s *ApplicationService) Create(ctx context.Context, input *CreateApplicationInput) (*Application, error) {
	inp := map[string]any{
		"name": input.Name,
	}
	vars := map[string]any{"input": inp}
	var result struct {
		CreateApplication *Application `json:"createApplication"`
	}
	if err := s.gql.do(ctx, createApplicationMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.CreateApplication, nil
}

// Update partially updates an existing Application by ID.
func (s *ApplicationService) Update(ctx context.Context, id string, input *UpdateApplicationInput) (*Application, error) {
	inp := map[string]any{}
	if input.Name != nil {
		inp["name"] = *input.Name
	}
	vars := map[string]any{"id": id, "input": inp}
	var result struct {
		UpdateApplication *Application `json:"updateApplication"`
	}
	if err := s.gql.do(ctx, updateApplicationMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.UpdateApplication, nil
}

// Delete removes a Application by ID.
func (s *ApplicationService) Delete(ctx context.Context, id string) error {
	vars := map[string]any{"id": id}
	return s.gql.do(ctx, deleteApplicationMutation, vars, nil)
}
