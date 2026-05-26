package hivehook

import "context"

// EndpointService manages Endpoints owned by Applications.
type EndpointService struct {
	gql *graphqlClient
}

// ListEndpointsOptions filters list results.
type ListEndpointsOptions struct {
	ListOptions
	ApplicationID *string
	Status        EndpointStatus
}

// CreateEndpointInput is the input for the matching Create method.
type CreateEndpointInput struct {
	ApplicationID string
	URL           string
	Type          DestinationType
	TypeConfig    map[string]any
	FilterConfig  *FilterConfig
	RateLimitRps  *int
	TimeoutMs     *int
	RetryPolicy   *RetryPolicy
	Headers       map[string]any
	AuthType      AuthType
	OAuth2Config  *OAuth2Config
	MTLSCert      string
	MTLSKey       string
	DeliveryMode  DeliveryMode
	Ordered       *bool
	HealthConfig  *HealthConfig
	OutputFormat  string
}

// UpdateEndpointInput is the input for the matching Update method.
type UpdateEndpointInput struct {
	URL          *string
	FilterConfig *FilterConfig
	Status       EndpointStatus
	Type         *DestinationType
	TypeConfig   map[string]any
	RateLimitRps *int
	TimeoutMs    *int
	RetryPolicy  *RetryPolicy
	Headers      map[string]any
	AuthType     *AuthType
	OAuth2Config *OAuth2Config
	MTLSCert     *string
	MTLSKey      *string
	DeliveryMode *DeliveryMode
	Ordered      *bool
	HealthConfig *HealthConfig
	OutputFormat *string
}

const endpointFragment = `
	fragment EndpointFields on Endpoint {
		id applicationId url signingSecret status type typeConfig rateLimitRps timeoutMs
		retryPolicy { maxAttempts initialDelay maxDelay backoffFactor }
		filterConfig { eventTypes regex bodyMatch { path value operator } rules { path operator value rules { path operator value } } }
		headers authType oauth2Config { tokenUrl clientId clientSecret scopes audience } mtlsCert mtlsKey deliveryMode pollApiKeyPrefix pollApiKey ordered blockedDeliveryId healthScore disabledReason healthConfig { windowHours disableBelow } outputFormat createdAt
	}
`

const listEndpointsQuery = `
	query ListEndpoints($applicationId: UUID, $status: EndpointStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
		endpoints(applicationId: $applicationId, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
			nodes { ...EndpointFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + endpointFragment

const getEndpointQuery = `
	query GetEndpoint($id: UUID!) {
		endpoint(id: $id) { ...EndpointFields }
	}
` + endpointFragment

const createEndpointMutation = `
	mutation CreateEndpoint($input: CreateEndpointInput!) {
		createEndpoint(input: $input) { ...EndpointFields }
	}
` + endpointFragment

const updateEndpointMutation = `
	mutation UpdateEndpoint($id: UUID!, $input: UpdateEndpointInput!) {
		updateEndpoint(id: $id, input: $input) { ...EndpointFields }
	}
` + endpointFragment

const deleteEndpointMutation = `
	mutation DeleteEndpoint($id: UUID!) {
		deleteEndpoint(id: $id)
	}
`

const rotateEndpointSecretMutation = `
	mutation RotateEndpointSecret($id: UUID!) {
		rotateEndpointSecret(id: $id) { ...EndpointFields }
	}
` + endpointFragment

// List returns Endpoints matching opts and a PageInfo cursor.
func (s *EndpointService) List(ctx context.Context, opts *ListEndpointsOptions) ([]*Endpoint, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
		if opts.ApplicationID != nil {
			vars = mergeVars(vars, map[string]any{"applicationId": *opts.ApplicationID})
		}
		if opts.Status != "" {
			vars = mergeVars(vars, map[string]any{"status": opts.Status})
		}
	}
	var result struct {
		Endpoints struct {
			Nodes    []*Endpoint `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"endpoints"`
	}
	if err := s.gql.do(ctx, listEndpointsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Endpoints.Nodes, &result.Endpoints.PageInfo, nil
}

// Get fetches a single Endpoint by ID.
func (s *EndpointService) Get(ctx context.Context, id string) (*Endpoint, error) {
	vars := map[string]any{"id": id}
	var result struct {
		Endpoint *Endpoint `json:"endpoint"`
	}
	if err := s.gql.do(ctx, getEndpointQuery, vars, &result); err != nil {
		return nil, err
	}
	return result.Endpoint, nil
}

// Create persists a new Endpoint.
func (s *EndpointService) Create(ctx context.Context, input *CreateEndpointInput) (*Endpoint, error) {
	inp := map[string]any{
		"applicationId": input.ApplicationID,
		"url":           input.URL,
	}
	if input.Type != "" {
		inp["type"] = input.Type
	}
	if input.TypeConfig != nil {
		inp["typeConfig"] = input.TypeConfig
	}
	if input.FilterConfig != nil {
		inp["filterConfig"] = input.FilterConfig
	}
	if input.RateLimitRps != nil {
		inp["rateLimitRps"] = *input.RateLimitRps
	}
	if input.TimeoutMs != nil {
		inp["timeoutMs"] = *input.TimeoutMs
	}
	if input.RetryPolicy != nil {
		inp["retryPolicy"] = input.RetryPolicy
	}
	if input.Headers != nil {
		inp["headers"] = input.Headers
	}
	if input.AuthType != "" {
		inp["authType"] = input.AuthType
	}
	if input.OAuth2Config != nil {
		inp["oauth2Config"] = map[string]any{
			"tokenUrl":     input.OAuth2Config.TokenURL,
			"clientId":     input.OAuth2Config.ClientID,
			"clientSecret": input.OAuth2Config.ClientSecret,
			"scopes":       input.OAuth2Config.Scopes,
			"audience":     input.OAuth2Config.Audience,
		}
	}
	if input.MTLSCert != "" {
		inp["mtlsCert"] = input.MTLSCert
	}
	if input.MTLSKey != "" {
		inp["mtlsKey"] = input.MTLSKey
	}
	if input.DeliveryMode != "" {
		inp["deliveryMode"] = input.DeliveryMode
	}
	if input.Ordered != nil {
		inp["ordered"] = *input.Ordered
	}
	if input.HealthConfig != nil {
		inp["healthConfig"] = map[string]any{
			"windowHours":  input.HealthConfig.WindowHours,
			"disableBelow": input.HealthConfig.DisableBelow,
		}
	}
	if input.OutputFormat != "" {
		inp["outputFormat"] = input.OutputFormat
	}
	vars := map[string]any{"input": inp}
	var result struct {
		CreateEndpoint *Endpoint `json:"createEndpoint"`
	}
	if err := s.gql.do(ctx, createEndpointMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.CreateEndpoint, nil
}

// Update partially updates an existing Endpoint by ID.
func (s *EndpointService) Update(ctx context.Context, id string, input *UpdateEndpointInput) (*Endpoint, error) {
	inp := map[string]any{}
	if input.URL != nil {
		inp["url"] = *input.URL
	}
	if input.Type != nil {
		inp["type"] = *input.Type
	}
	if input.TypeConfig != nil {
		inp["typeConfig"] = input.TypeConfig
	}
	if input.FilterConfig != nil {
		inp["filterConfig"] = input.FilterConfig
	}
	if input.Status != "" {
		inp["status"] = input.Status
	}
	if input.RateLimitRps != nil {
		inp["rateLimitRps"] = *input.RateLimitRps
	}
	if input.TimeoutMs != nil {
		inp["timeoutMs"] = *input.TimeoutMs
	}
	if input.RetryPolicy != nil {
		inp["retryPolicy"] = input.RetryPolicy
	}
	if input.Headers != nil {
		inp["headers"] = input.Headers
	}
	if input.AuthType != nil {
		inp["authType"] = *input.AuthType
	}
	if input.OAuth2Config != nil {
		inp["oauth2Config"] = map[string]any{
			"tokenUrl":     input.OAuth2Config.TokenURL,
			"clientId":     input.OAuth2Config.ClientID,
			"clientSecret": input.OAuth2Config.ClientSecret,
			"scopes":       input.OAuth2Config.Scopes,
			"audience":     input.OAuth2Config.Audience,
		}
	}
	if input.MTLSCert != nil {
		inp["mtlsCert"] = *input.MTLSCert
	}
	if input.MTLSKey != nil {
		inp["mtlsKey"] = *input.MTLSKey
	}
	if input.DeliveryMode != nil {
		inp["deliveryMode"] = *input.DeliveryMode
	}
	if input.Ordered != nil {
		inp["ordered"] = *input.Ordered
	}
	if input.HealthConfig != nil {
		inp["healthConfig"] = map[string]any{
			"windowHours":  input.HealthConfig.WindowHours,
			"disableBelow": input.HealthConfig.DisableBelow,
		}
	}
	if input.OutputFormat != nil {
		inp["outputFormat"] = *input.OutputFormat
	}
	vars := map[string]any{"id": id, "input": inp}
	var result struct {
		UpdateEndpoint *Endpoint `json:"updateEndpoint"`
	}
	if err := s.gql.do(ctx, updateEndpointMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.UpdateEndpoint, nil
}

// Delete removes a Endpoint by ID.
func (s *EndpointService) Delete(ctx context.Context, id string) error {
	vars := map[string]any{"id": id}
	return s.gql.do(ctx, deleteEndpointMutation, vars, nil)
}

// RotateSecret generates a new signing secret and keeps the prior one as secondary.
func (s *EndpointService) RotateSecret(ctx context.Context, id string) (*Endpoint, error) {
	vars := map[string]any{"id": id}
	var result struct {
		RotateEndpointSecret *Endpoint `json:"rotateEndpointSecret"`
	}
	if err := s.gql.do(ctx, rotateEndpointSecretMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.RotateEndpointSecret, nil
}

const pollOutboundDeliveriesQuery = `
	query PollOutboundDeliveries($endpointId: UUID!, $cursor: String, $limit: Int) {
		pollOutboundDeliveries(endpointId: $endpointId, cursor: $cursor, limit: $limit) {
			nodes { ...OutboundDeliveryFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + outboundDeliveryFragment

const ackOutboundDeliveriesMutation = `mutation($endpointId: UUID!, $deliveryIds: [UUID!]!) {
	ackOutboundDeliveries(endpointId: $endpointId, deliveryIds: $deliveryIds)
}`

const regenerateOutboundPollApiKeyMutation = `
	mutation RegenerateOutboundPollApiKey($endpointId: UUID!) {
		regenerateOutboundPollApiKey(endpointId: $endpointId) { ...EndpointFields }
	}
` + endpointFragment

// PollDeliveries pulls pending Deliveries for a poll-mode Destination (or Endpoint).
func (s *EndpointService) PollDeliveries(ctx context.Context, endpointID string, cursor *string, limit *int) ([]*OutboundDelivery, *PageInfo, error) {
	vars := map[string]any{"endpointId": endpointID}
	if cursor != nil {
		vars["cursor"] = *cursor
	}
	if limit != nil {
		vars["limit"] = *limit
	}
	var result struct {
		PollOutboundDeliveries struct {
			Nodes    []*OutboundDelivery `json:"nodes"`
			PageInfo PageInfo            `json:"pageInfo"`
		} `json:"pollOutboundDeliveries"`
	}
	if err := s.gql.do(ctx, pollOutboundDeliveriesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.PollOutboundDeliveries.Nodes, &result.PollOutboundDeliveries.PageInfo, nil
}

// AckDeliveries acknowledges (and removes) Deliveries pulled by PollDeliveries.
func (s *EndpointService) AckDeliveries(ctx context.Context, endpointID string, deliveryIDs []string) (int, error) {
	var result struct {
		AckOutboundDeliveries int `json:"ackOutboundDeliveries"`
	}
	if err := s.gql.do(ctx, ackOutboundDeliveriesMutation, map[string]any{"endpointId": endpointID, "deliveryIds": deliveryIDs}, &result); err != nil {
		return 0, err
	}
	return result.AckOutboundDeliveries, nil
}

// RegeneratePollAPIKey rotates the API key used to authenticate poll-mode clients.
func (s *EndpointService) RegeneratePollAPIKey(ctx context.Context, endpointID string) (*Endpoint, error) {
	var result struct {
		RegenerateOutboundPollApiKey *Endpoint `json:"regenerateOutboundPollApiKey"`
	}
	if err := s.gql.do(ctx, regenerateOutboundPollApiKeyMutation, map[string]any{"endpointId": endpointID}, &result); err != nil {
		return nil, err
	}
	return result.RegenerateOutboundPollApiKey, nil
}

const skipOutboundDLQEntryMutation = `
	mutation SkipOutboundDLQEntry($id: UUID!) {
		skipOutboundDlqEntry(id: $id)
	}
`

// SkipOutboundDLQEntry marks an outbound DLQ entry as skipped (no replay).
func (s *EndpointService) SkipOutboundDLQEntry(ctx context.Context, id string) error {
	return s.gql.do(ctx, skipOutboundDLQEntryMutation, map[string]any{"id": id}, nil)
}
