package hivehook

import "context"

const organizationFragment = `
	id name slug ssoEnabled ssoProvider retentionEvents retentionMessages
	otlpConfig { endpoint headers insecure sampleRate }
	createdAt updatedAt
`

const listOrganizationsQuery = `query($search: String, $limit: Int, $offset: Int) {
	organizations(search: $search, limit: $limit, offset: $offset) {
		nodes {` + organizationFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getOrganizationQuery = `query($id: UUID!) {
	organization(id: $id) {` + organizationFragment + `}
}`

const createOrganizationMutation = `mutation($input: CreateOrganizationInput!) {
	createOrganization(input: $input) {` + organizationFragment + `}
}`

const updateOrganizationMutation = `mutation($id: UUID!, $input: UpdateOrganizationInput!) {
	updateOrganization(id: $id, input: $input) {` + organizationFragment + `}
}`

const deleteOrganizationMutation = `mutation($id: UUID!) {
	deleteOrganization(id: $id)
}`

const configureSSOmutation = `mutation($organizationId: UUID!, $input: SSOConfigInput!) {
	configureSSO(organizationId: $organizationId, input: $input) {` + organizationFragment + `}
}`

const disableSSOmutation = `mutation($organizationId: UUID!) {
	disableSSO(organizationId: $organizationId) {` + organizationFragment + `}
}`

const updateOrganizationRetentionMutation = `mutation($organizationId: UUID!, $input: RetentionInput!) {
	updateOrganizationRetention(organizationId: $organizationId, input: $input) {` + organizationFragment + `}
}`

const deleteOrganizationDataMutation = `mutation($organizationId: UUID!) {
	deleteOrganizationData(organizationId: $organizationId)
}`

const exportOrganizationDataQuery = `mutation($organizationId: UUID!) {
	exportOrganizationData(organizationId: $organizationId)
}`

const configureOTLPmutation = `mutation($organizationId: UUID!, $input: OTLPConfigInput!) {
	configureOTLP(organizationId: $organizationId, input: $input) {` + organizationFragment + `}
}`

const disableOTLPmutation = `mutation($organizationId: UUID!) {
	disableOTLP(organizationId: $organizationId) {` + organizationFragment + `}
}`

// OrganizationService manages Organizations.
type OrganizationService struct {
	gql *graphqlClient
}

// CreateOrganizationInput is the input for the matching Create method.
type CreateOrganizationInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UpdateOrganizationInput is the input for the matching Update method.
type UpdateOrganizationInput struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// SSOConfigInput is the input for the matching service method.
type SSOConfigInput struct {
	Provider       string  `json:"provider"`
	IDPMetadataURL *string `json:"idpMetadataUrl,omitempty"`
	EntityID       *string `json:"entityId,omitempty"`
	ACSBaseURL     *string `json:"acsBaseUrl,omitempty"`
	Issuer         *string `json:"issuer,omitempty"`
	ClientID       *string `json:"clientId,omitempty"`
	ClientSecret   *string `json:"clientSecret,omitempty"`
	RedirectURL    *string `json:"redirectUrl,omitempty"`
}

// RetentionInput is the input for the matching service method.
type RetentionInput struct {
	RetentionEvents   int `json:"retentionEvents"`
	RetentionMessages int `json:"retentionMessages"`
}

// OTLPConfigInput is the input for the matching service method.
type OTLPConfigInput struct {
	Endpoint   string         `json:"endpoint"`
	Headers    map[string]any `json:"headers,omitempty"`
	Insecure   *bool          `json:"insecure,omitempty"`
	SampleRate *float64       `json:"sampleRate,omitempty"`
}

// List returns Organizations matching opts and a PageInfo cursor.
func (s *OrganizationService) List(ctx context.Context, opts *ListOptions) ([]*Organization, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
	}
	var result struct {
		Organizations struct {
			Nodes    []*Organization `json:"nodes"`
			PageInfo PageInfo        `json:"pageInfo"`
		} `json:"organizations"`
	}
	if err := s.gql.do(ctx, listOrganizationsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Organizations.Nodes, &result.Organizations.PageInfo, nil
}

// Get fetches a single Organization by ID.
func (s *OrganizationService) Get(ctx context.Context, id string) (*Organization, error) {
	var result struct {
		Organization *Organization `json:"organization"`
	}
	if err := s.gql.do(ctx, getOrganizationQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Organization, nil
}

// Create persists a new Organization.
func (s *OrganizationService) Create(ctx context.Context, input *CreateOrganizationInput) (*Organization, error) {
	inp := map[string]any{
		"name": input.Name,
		"slug": input.Slug,
	}
	vars := map[string]any{"input": inp}
	var result struct {
		CreateOrganization *Organization `json:"createOrganization"`
	}
	if err := s.gql.do(ctx, createOrganizationMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.CreateOrganization, nil
}

// Update partially updates an existing Organization by ID.
func (s *OrganizationService) Update(ctx context.Context, id string, input *UpdateOrganizationInput) (*Organization, error) {
	inp := map[string]any{}
	if input.Name != nil {
		inp["name"] = *input.Name
	}
	if input.Slug != nil {
		inp["slug"] = *input.Slug
	}
	vars := map[string]any{"id": id, "input": inp}
	var result struct {
		UpdateOrganization *Organization `json:"updateOrganization"`
	}
	if err := s.gql.do(ctx, updateOrganizationMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.UpdateOrganization, nil
}

// Delete removes a Organization by ID.
func (s *OrganizationService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteOrganizationMutation, map[string]any{"id": id}, nil)
}

// ConfigureSSO enables and configures SSO for an Organization.
func (s *OrganizationService) ConfigureSSO(ctx context.Context, organizationID string, input *SSOConfigInput) (*Organization, error) {
	inp := map[string]any{
		"provider": input.Provider,
	}
	if input.IDPMetadataURL != nil {
		inp["idpMetadataUrl"] = *input.IDPMetadataURL
	}
	if input.EntityID != nil {
		inp["entityId"] = *input.EntityID
	}
	if input.ACSBaseURL != nil {
		inp["acsBaseUrl"] = *input.ACSBaseURL
	}
	if input.Issuer != nil {
		inp["issuer"] = *input.Issuer
	}
	if input.ClientID != nil {
		inp["clientId"] = *input.ClientID
	}
	if input.ClientSecret != nil {
		inp["clientSecret"] = *input.ClientSecret
	}
	if input.RedirectURL != nil {
		inp["redirectUrl"] = *input.RedirectURL
	}
	vars := map[string]any{"organizationId": organizationID, "input": inp}
	var result struct {
		ConfigureSSO *Organization `json:"configureSSO"`
	}
	if err := s.gql.do(ctx, configureSSOmutation, vars, &result); err != nil {
		return nil, err
	}
	return result.ConfigureSSO, nil
}

// DisableSSO disables SSO for an Organization.
func (s *OrganizationService) DisableSSO(ctx context.Context, organizationID string) (*Organization, error) {
	vars := map[string]any{"organizationId": organizationID}
	var result struct {
		DisableSSO *Organization `json:"disableSSO"`
	}
	if err := s.gql.do(ctx, disableSSOmutation, vars, &result); err != nil {
		return nil, err
	}
	return result.DisableSSO, nil
}

// UpdateRetention updates event/message retention windows for an Organization.
func (s *OrganizationService) UpdateRetention(ctx context.Context, organizationID string, input *RetentionInput) (*Organization, error) {
	inp := map[string]any{
		"retentionEvents":   input.RetentionEvents,
		"retentionMessages": input.RetentionMessages,
	}
	vars := map[string]any{"organizationId": organizationID, "input": inp}
	var result struct {
		UpdateOrganizationRetention *Organization `json:"updateOrganizationRetention"`
	}
	if err := s.gql.do(ctx, updateOrganizationRetentionMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.UpdateOrganizationRetention, nil
}

// DeleteData permanently deletes all data for an Organization.
func (s *OrganizationService) DeleteData(ctx context.Context, organizationID string) error {
	return s.gql.do(ctx, deleteOrganizationDataMutation, map[string]any{"organizationId": organizationID}, nil)
}

// ExportData exports all data for an Organization as a JSON document.
func (s *OrganizationService) ExportData(ctx context.Context, organizationID string) (map[string]any, error) {
	vars := map[string]any{"organizationId": organizationID}
	var result struct {
		ExportOrganizationData map[string]any `json:"exportOrganizationData"`
	}
	if err := s.gql.do(ctx, exportOrganizationDataQuery, vars, &result); err != nil {
		return nil, err
	}
	return result.ExportOrganizationData, nil
}

// ConfigureOTLP enables and configures the OTLP exporter for an Organization.
func (s *OrganizationService) ConfigureOTLP(ctx context.Context, organizationID string, input *OTLPConfigInput) (*Organization, error) {
	inp := map[string]any{
		"endpoint": input.Endpoint,
	}
	if input.Headers != nil {
		inp["headers"] = input.Headers
	}
	if input.Insecure != nil {
		inp["insecure"] = *input.Insecure
	}
	if input.SampleRate != nil {
		inp["sampleRate"] = *input.SampleRate
	}
	vars := map[string]any{"organizationId": organizationID, "input": inp}
	var result struct {
		ConfigureOTLP *Organization `json:"configureOTLP"`
	}
	if err := s.gql.do(ctx, configureOTLPmutation, vars, &result); err != nil {
		return nil, err
	}
	return result.ConfigureOTLP, nil
}

// DisableOTLP disables the OTLP exporter for an Organization.
func (s *OrganizationService) DisableOTLP(ctx context.Context, organizationID string) (*Organization, error) {
	vars := map[string]any{"organizationId": organizationID}
	var result struct {
		DisableOTLP *Organization `json:"disableOTLP"`
	}
	if err := s.gql.do(ctx, disableOTLPmutation, vars, &result); err != nil {
		return nil, err
	}
	return result.DisableOTLP, nil
}
