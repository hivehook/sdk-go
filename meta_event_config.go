package hivehook

import "context"

const metaEventConfigFragment = `
	id name url signingSecret eventTypes enabled createdAt
`

const listMetaEventConfigsQuery = `query($search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	metaEventConfigs(search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + metaEventConfigFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getMetaEventConfigQuery = `query($id: UUID!) {
	metaEventConfig(id: $id) {` + metaEventConfigFragment + `}
}`

const createMetaEventConfigMutation = `mutation($input: CreateMetaEventConfigInput!) {
	createMetaEventConfig(input: $input) {` + metaEventConfigFragment + `}
}`

const updateMetaEventConfigMutation = `mutation($id: UUID!, $input: UpdateMetaEventConfigInput!) {
	updateMetaEventConfig(id: $id, input: $input) {` + metaEventConfigFragment + `}
}`

const deleteMetaEventConfigMutation = `mutation($id: UUID!) {
	deleteMetaEventConfig(id: $id)
}`

// MetaEventConfigService manages meta-event webhook configurations.
type MetaEventConfigService struct {
	gql *graphqlClient
}

// CreateMetaEventConfigInput is the input for MetaEventConfigService.Create.
type CreateMetaEventConfigInput struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	EventTypes []string `json:"eventTypes"`
	Enabled    *bool    `json:"enabled,omitempty"`
}

// UpdateMetaEventConfigInput is the input for MetaEventConfigService.Update.
type UpdateMetaEventConfigInput struct {
	Name       *string  `json:"name,omitempty"`
	URL        *string  `json:"url,omitempty"`
	EventTypes []string `json:"eventTypes,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
}

// List returns meta-event configs matching opts and a PageInfo cursor.
func (s *MetaEventConfigService) List(ctx context.Context, opts *ListOptions) ([]*MetaEventConfig, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
	}
	var result struct {
		MetaEventConfigs struct {
			Nodes    []*MetaEventConfig `json:"nodes"`
			PageInfo PageInfo           `json:"pageInfo"`
		} `json:"metaEventConfigs"`
	}
	if err := s.gql.do(ctx, listMetaEventConfigsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.MetaEventConfigs.Nodes, &result.MetaEventConfigs.PageInfo, nil
}

// Get fetches a single MetaEventConfig by ID.
func (s *MetaEventConfigService) Get(ctx context.Context, id string) (*MetaEventConfig, error) {
	var result struct {
		MetaEventConfig *MetaEventConfig `json:"metaEventConfig"`
	}
	if err := s.gql.do(ctx, getMetaEventConfigQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.MetaEventConfig, nil
}

// Create persists a new MetaEventConfig.
func (s *MetaEventConfigService) Create(ctx context.Context, input *CreateMetaEventConfigInput) (*MetaEventConfig, error) {
	m := map[string]any{
		"name":       input.Name,
		"url":        input.URL,
		"eventTypes": input.EventTypes,
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	var result struct {
		CreateMetaEventConfig MetaEventConfig `json:"createMetaEventConfig"`
	}
	if err := s.gql.do(ctx, createMetaEventConfigMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateMetaEventConfig, nil
}

// Update partially updates an existing MetaEventConfig by ID.
func (s *MetaEventConfigService) Update(ctx context.Context, id string, input *UpdateMetaEventConfigInput) (*MetaEventConfig, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.URL != nil {
		m["url"] = *input.URL
	}
	if input.EventTypes != nil {
		m["eventTypes"] = input.EventTypes
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	var result struct {
		UpdateMetaEventConfig MetaEventConfig `json:"updateMetaEventConfig"`
	}
	if err := s.gql.do(ctx, updateMetaEventConfigMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateMetaEventConfig, nil
}

// Delete removes a MetaEventConfig by ID.
func (s *MetaEventConfigService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteMetaEventConfigMutation, map[string]any{"id": id}, nil)
}
