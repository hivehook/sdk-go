package hivehook

import "context"

const sourceFragment = `
	id name slug providerType verifyConfig status rateLimitRps spikeProtection maxIngestRps createdAt
`

const listSourcesQuery = `query($status: SourceStatus, $providerType: String, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	sources(status: $status, providerType: $providerType, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + sourceFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getSourceQuery = `query($id: UUID!) {
	source(id: $id) {` + sourceFragment + `}
}`

const createSourceMutation = `mutation($input: CreateSourceInput!) {
	createSource(input: $input) {` + sourceFragment + `}
}`

const updateSourceMutation = `mutation($id: UUID!, $input: UpdateSourceInput!) {
	updateSource(id: $id, input: $input) {` + sourceFragment + `}
}`

const deleteSourceMutation = `mutation($id: UUID!) {
	deleteSource(id: $id)
}`

const rotateSourceSecretMutation = `mutation($id: UUID!) {
	rotateSourceSecret(id: $id) {` + sourceFragment + `}
}`

const clearSourceSecondarySecretMutation = `mutation($id: UUID!) {
	clearSourceSecondarySecret(id: $id) {` + sourceFragment + `}
}`

// SourceService manages inbound webhook Sources.
type SourceService struct {
	gql *graphqlClient
}

// ListSourcesOptions filters SourceService.List results.
type ListSourcesOptions struct {
	ListOptions
	Status       SourceStatus
	ProviderType string
}

// CreateSourceInput is the input for SourceService.Create.
type CreateSourceInput struct {
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	ProviderType    string         `json:"providerType"`
	VerifyConfig    map[string]any `json:"verifyConfig,omitempty"`
	RateLimitRps    *int           `json:"rateLimitRps,omitempty"`
	Status          *SourceStatus  `json:"status,omitempty"`
	SpikeProtection *bool          `json:"spikeProtection,omitempty"`
	MaxIngestRps    *int           `json:"maxIngestRps,omitempty"`
}

// UpdateSourceInput is the input for SourceService.Update.
type UpdateSourceInput struct {
	Name            *string        `json:"name,omitempty"`
	Slug            *string        `json:"slug,omitempty"`
	ProviderType    *string        `json:"providerType,omitempty"`
	VerifyConfig    map[string]any `json:"verifyConfig,omitempty"`
	Status          *SourceStatus  `json:"status,omitempty"`
	RateLimitRps    *int           `json:"rateLimitRps,omitempty"`
	SpikeProtection *bool          `json:"spikeProtection,omitempty"`
	MaxIngestRps    *int           `json:"maxIngestRps,omitempty"`
}

// List returns Sources matching opts and a PageInfo cursor.
func (s *SourceService) List(ctx context.Context, opts *ListSourcesOptions) ([]*Source, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
		if opts.ProviderType != "" {
			vars["providerType"] = opts.ProviderType
		}
	}
	var result struct {
		Sources struct {
			Nodes    []*Source `json:"nodes"`
			PageInfo PageInfo  `json:"pageInfo"`
		} `json:"sources"`
	}
	if err := s.gql.do(ctx, listSourcesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Sources.Nodes, &result.Sources.PageInfo, nil
}

// Get fetches a single Source by ID.
func (s *SourceService) Get(ctx context.Context, id string) (*Source, error) {
	var result struct {
		Source *Source `json:"source"`
	}
	if err := s.gql.do(ctx, getSourceQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Source, nil
}

// Create persists a new Source.
func (s *SourceService) Create(ctx context.Context, input *CreateSourceInput) (*Source, error) {
	m := map[string]any{
		"name":         input.Name,
		"slug":         input.Slug,
		"providerType": input.ProviderType,
	}
	if input.VerifyConfig != nil {
		m["verifyConfig"] = input.VerifyConfig
	}
	if input.RateLimitRps != nil {
		m["rateLimitRps"] = *input.RateLimitRps
	}
	if input.Status != nil {
		m["status"] = *input.Status
	}
	if input.SpikeProtection != nil {
		m["spikeProtection"] = *input.SpikeProtection
	}
	if input.MaxIngestRps != nil {
		m["maxIngestRps"] = *input.MaxIngestRps
	}
	var result struct {
		CreateSource Source `json:"createSource"`
	}
	if err := s.gql.do(ctx, createSourceMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateSource, nil
}

// Update partially updates an existing Source by ID.
func (s *SourceService) Update(ctx context.Context, id string, input *UpdateSourceInput) (*Source, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.Slug != nil {
		m["slug"] = *input.Slug
	}
	if input.ProviderType != nil {
		m["providerType"] = *input.ProviderType
	}
	if input.VerifyConfig != nil {
		m["verifyConfig"] = input.VerifyConfig
	}
	if input.Status != nil {
		m["status"] = *input.Status
	}
	if input.RateLimitRps != nil {
		m["rateLimitRps"] = *input.RateLimitRps
	}
	if input.SpikeProtection != nil {
		m["spikeProtection"] = *input.SpikeProtection
	}
	if input.MaxIngestRps != nil {
		m["maxIngestRps"] = *input.MaxIngestRps
	}
	var result struct {
		UpdateSource Source `json:"updateSource"`
	}
	if err := s.gql.do(ctx, updateSourceMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateSource, nil
}

// Delete removes a Source by ID.
func (s *SourceService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteSourceMutation, map[string]any{"id": id}, nil)
}

// RotateSecret generates a new signing secret for the Source and keeps the previous
// one as secondary until ClearSecondarySecret is called.
func (s *SourceService) RotateSecret(ctx context.Context, id string) (*Source, error) {
	var result struct {
		RotateSourceSecret Source `json:"rotateSourceSecret"`
	}
	if err := s.gql.do(ctx, rotateSourceSecretMutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return &result.RotateSourceSecret, nil
}

// ClearSecondarySecret drops the previous secret kept during rotation.
func (s *SourceService) ClearSecondarySecret(ctx context.Context, id string) (*Source, error) {
	var result struct {
		ClearSourceSecondarySecret Source `json:"clearSourceSecondarySecret"`
	}
	if err := s.gql.do(ctx, clearSourceSecondarySecretMutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return &result.ClearSourceSecondarySecret, nil
}
