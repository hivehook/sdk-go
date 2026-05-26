package hivehook

import "context"

const destinationFragment = `
	id name url signingSecret status type typeConfig timeoutMs rateLimitRps retryPolicy { maxAttempts initialDelay maxDelay backoffFactor } headers authType oauth2Config { tokenUrl clientId clientSecret scopes audience } mtlsCert mtlsKey deliveryMode pollApiKeyPrefix pollApiKey ordered blockedDeliveryId healthScore disabledReason healthConfig { windowHours disableBelow } outputFormat createdAt
`

const listDestinationsQuery = `query($status: DestinationStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	destinations(status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + destinationFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getDestinationQuery = `query($id: UUID!) {
	destination(id: $id) {` + destinationFragment + `}
}`

const createDestinationMutation = `mutation($input: CreateDestinationInput!) {
	createDestination(input: $input) {` + destinationFragment + `}
}`

const updateDestinationMutation = `mutation($id: UUID!, $input: UpdateDestinationInput!) {
	updateDestination(id: $id, input: $input) {` + destinationFragment + `}
}`

const deleteDestinationMutation = `mutation($id: UUID!) {
	deleteDestination(id: $id)
}`

const rotateDestinationSecretMutation = `mutation($id: UUID!) {
	rotateDestinationSecret(id: $id) {` + destinationFragment + `}
}`

// DestinationService manages outbound Destinations.
type DestinationService struct {
	gql *graphqlClient
}

// ListDestinationsOptions filters list results.
type ListDestinationsOptions struct {
	ListOptions
	Status DestinationStatus
}

// CreateDestinationInput is the input for the matching Create method.
type CreateDestinationInput struct {
	Name         string          `json:"name"`
	URL          string          `json:"url"`
	Type         DestinationType `json:"type,omitempty"`
	TypeConfig   map[string]any  `json:"typeConfig,omitempty"`
	TimeoutMs    *int            `json:"timeoutMs,omitempty"`
	RateLimitRps *int            `json:"rateLimitRps,omitempty"`
	RetryPolicy  *RetryPolicy    `json:"retryPolicy,omitempty"`
	Headers      map[string]any  `json:"headers,omitempty"`
	AuthType     AuthType        `json:"authType,omitempty"`
	OAuth2Config *OAuth2Config   `json:"oauth2Config,omitempty"`
	MTLSCert     string          `json:"mtlsCert,omitempty"`
	MTLSKey      string          `json:"mtlsKey,omitempty"`
	DeliveryMode DeliveryMode    `json:"deliveryMode,omitempty"`
	Ordered      *bool           `json:"ordered,omitempty"`
	HealthConfig *HealthConfig   `json:"healthConfig,omitempty"`
	OutputFormat string          `json:"outputFormat,omitempty"`
}

// UpdateDestinationInput is the input for the matching Update method.
type UpdateDestinationInput struct {
	Name         *string            `json:"name,omitempty"`
	URL          *string            `json:"url,omitempty"`
	Status       *DestinationStatus `json:"status,omitempty"`
	Type         *DestinationType   `json:"type,omitempty"`
	TypeConfig   map[string]any     `json:"typeConfig,omitempty"`
	TimeoutMs    *int               `json:"timeoutMs,omitempty"`
	RateLimitRps *int               `json:"rateLimitRps,omitempty"`
	RetryPolicy  *RetryPolicy       `json:"retryPolicy,omitempty"`
	Headers      map[string]any     `json:"headers,omitempty"`
	AuthType     *AuthType          `json:"authType,omitempty"`
	OAuth2Config *OAuth2Config      `json:"oauth2Config,omitempty"`
	MTLSCert     *string            `json:"mtlsCert,omitempty"`
	MTLSKey      *string            `json:"mtlsKey,omitempty"`
	DeliveryMode *DeliveryMode      `json:"deliveryMode,omitempty"`
	Ordered      *bool              `json:"ordered,omitempty"`
	HealthConfig *HealthConfig      `json:"healthConfig,omitempty"`
	OutputFormat *string            `json:"outputFormat,omitempty"`
}

// List returns Destinations matching opts and a PageInfo cursor.
func (s *DestinationService) List(ctx context.Context, opts *ListDestinationsOptions) ([]*Destination, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
	}
	var result struct {
		Destinations struct {
			Nodes    []*Destination `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"destinations"`
	}
	if err := s.gql.do(ctx, listDestinationsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Destinations.Nodes, &result.Destinations.PageInfo, nil
}

// Get fetches a single Destination by ID.
func (s *DestinationService) Get(ctx context.Context, id string) (*Destination, error) {
	var result struct {
		Destination *Destination `json:"destination"`
	}
	if err := s.gql.do(ctx, getDestinationQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Destination, nil
}

// Create persists a new Destination.
func (s *DestinationService) Create(ctx context.Context, input *CreateDestinationInput) (*Destination, error) {
	m := map[string]any{
		"name": input.Name,
		"url":  input.URL,
	}
	if input.Type != "" {
		m["type"] = input.Type
	}
	if input.TypeConfig != nil {
		m["typeConfig"] = input.TypeConfig
	}
	if input.TimeoutMs != nil {
		m["timeoutMs"] = *input.TimeoutMs
	}
	if input.RateLimitRps != nil {
		m["rateLimitRps"] = *input.RateLimitRps
	}
	if input.RetryPolicy != nil {
		m["retryPolicy"] = map[string]any{
			"maxAttempts":   input.RetryPolicy.MaxAttempts,
			"initialDelay":  input.RetryPolicy.InitialDelay,
			"maxDelay":      input.RetryPolicy.MaxDelay,
			"backoffFactor": input.RetryPolicy.BackoffFactor,
		}
	}
	if input.Headers != nil {
		m["headers"] = input.Headers
	}
	if input.AuthType != "" {
		m["authType"] = input.AuthType
	}
	if input.OAuth2Config != nil {
		m["oauth2Config"] = map[string]any{
			"tokenUrl":     input.OAuth2Config.TokenURL,
			"clientId":     input.OAuth2Config.ClientID,
			"clientSecret": input.OAuth2Config.ClientSecret,
			"scopes":       input.OAuth2Config.Scopes,
			"audience":     input.OAuth2Config.Audience,
		}
	}
	if input.MTLSCert != "" {
		m["mtlsCert"] = input.MTLSCert
	}
	if input.MTLSKey != "" {
		m["mtlsKey"] = input.MTLSKey
	}
	if input.DeliveryMode != "" {
		m["deliveryMode"] = input.DeliveryMode
	}
	if input.Ordered != nil {
		m["ordered"] = *input.Ordered
	}
	if input.HealthConfig != nil {
		m["healthConfig"] = map[string]any{
			"windowHours":  input.HealthConfig.WindowHours,
			"disableBelow": input.HealthConfig.DisableBelow,
		}
	}
	if input.OutputFormat != "" {
		m["outputFormat"] = input.OutputFormat
	}
	var result struct {
		CreateDestination Destination `json:"createDestination"`
	}
	if err := s.gql.do(ctx, createDestinationMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateDestination, nil
}

// Update partially updates an existing Destination by ID.
func (s *DestinationService) Update(ctx context.Context, id string, input *UpdateDestinationInput) (*Destination, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.URL != nil {
		m["url"] = *input.URL
	}
	if input.Status != nil {
		m["status"] = *input.Status
	}
	if input.Type != nil {
		m["type"] = *input.Type
	}
	if input.TypeConfig != nil {
		m["typeConfig"] = input.TypeConfig
	}
	if input.TimeoutMs != nil {
		m["timeoutMs"] = *input.TimeoutMs
	}
	if input.RateLimitRps != nil {
		m["rateLimitRps"] = *input.RateLimitRps
	}
	if input.RetryPolicy != nil {
		m["retryPolicy"] = map[string]any{
			"maxAttempts":   input.RetryPolicy.MaxAttempts,
			"initialDelay":  input.RetryPolicy.InitialDelay,
			"maxDelay":      input.RetryPolicy.MaxDelay,
			"backoffFactor": input.RetryPolicy.BackoffFactor,
		}
	}
	if input.Headers != nil {
		m["headers"] = input.Headers
	}
	if input.AuthType != nil {
		m["authType"] = *input.AuthType
	}
	if input.OAuth2Config != nil {
		m["oauth2Config"] = map[string]any{
			"tokenUrl":     input.OAuth2Config.TokenURL,
			"clientId":     input.OAuth2Config.ClientID,
			"clientSecret": input.OAuth2Config.ClientSecret,
			"scopes":       input.OAuth2Config.Scopes,
			"audience":     input.OAuth2Config.Audience,
		}
	}
	if input.MTLSCert != nil {
		m["mtlsCert"] = *input.MTLSCert
	}
	if input.MTLSKey != nil {
		m["mtlsKey"] = *input.MTLSKey
	}
	if input.DeliveryMode != nil {
		m["deliveryMode"] = *input.DeliveryMode
	}
	if input.Ordered != nil {
		m["ordered"] = *input.Ordered
	}
	if input.HealthConfig != nil {
		m["healthConfig"] = map[string]any{
			"windowHours":  input.HealthConfig.WindowHours,
			"disableBelow": input.HealthConfig.DisableBelow,
		}
	}
	if input.OutputFormat != nil {
		m["outputFormat"] = *input.OutputFormat
	}
	var result struct {
		UpdateDestination Destination `json:"updateDestination"`
	}
	if err := s.gql.do(ctx, updateDestinationMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateDestination, nil
}

// Delete removes a Destination by ID.
func (s *DestinationService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteDestinationMutation, map[string]any{"id": id}, nil)
}

// RotateSecret generates a new signing secret and keeps the prior one as secondary.
func (s *DestinationService) RotateSecret(ctx context.Context, id string) (*Destination, error) {
	var result struct {
		RotateDestinationSecret Destination `json:"rotateDestinationSecret"`
	}
	if err := s.gql.do(ctx, rotateDestinationSecretMutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return &result.RotateDestinationSecret, nil
}

const pollDeliveriesQuery = `query($destinationId: UUID!, $cursor: String, $limit: Int) {
	pollDeliveries(destinationId: $destinationId, cursor: $cursor, limit: $limit) {
		nodes {` + deliveryFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const ackDeliveriesMutation = `mutation($destinationId: UUID!, $deliveryIds: [UUID!]!) {
	ackDeliveries(destinationId: $destinationId, deliveryIds: $deliveryIds)
}`

const regeneratePollApiKeyMutation = `mutation($destinationId: UUID!) {
	regeneratePollApiKey(destinationId: $destinationId) {` + destinationFragment + `}
}`

// PollDeliveries pulls pending Deliveries for a poll-mode Destination (or Endpoint).
func (s *DestinationService) PollDeliveries(ctx context.Context, destinationID string, cursor *string, limit *int) ([]*Delivery, *PageInfo, error) {
	vars := map[string]any{"destinationId": destinationID}
	if cursor != nil {
		vars["cursor"] = *cursor
	}
	if limit != nil {
		vars["limit"] = *limit
	}
	var result struct {
		PollDeliveries struct {
			Nodes    []*Delivery `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"pollDeliveries"`
	}
	if err := s.gql.do(ctx, pollDeliveriesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.PollDeliveries.Nodes, &result.PollDeliveries.PageInfo, nil
}

// AckDeliveries acknowledges (and removes) Deliveries pulled by PollDeliveries.
func (s *DestinationService) AckDeliveries(ctx context.Context, destinationID string, deliveryIDs []string) (int, error) {
	var result struct {
		AckDeliveries int `json:"ackDeliveries"`
	}
	if err := s.gql.do(ctx, ackDeliveriesMutation, map[string]any{"destinationId": destinationID, "deliveryIds": deliveryIDs}, &result); err != nil {
		return 0, err
	}
	return result.AckDeliveries, nil
}

// RegeneratePollAPIKey rotates the API key used to authenticate poll-mode clients.
func (s *DestinationService) RegeneratePollAPIKey(ctx context.Context, destinationID string) (*Destination, error) {
	var result struct {
		RegeneratePollApiKey Destination `json:"regeneratePollApiKey"`
	}
	if err := s.gql.do(ctx, regeneratePollApiKeyMutation, map[string]any{"destinationId": destinationID}, &result); err != nil {
		return nil, err
	}
	return &result.RegeneratePollApiKey, nil
}

const skipDLQEntryMutation = `mutation($id: UUID!) {
	skipDLQEntry(id: $id)
}`

// SkipDLQEntry marks a DLQ entry as skipped (no replay).
func (s *DestinationService) SkipDLQEntry(ctx context.Context, id string) error {
	return s.gql.do(ctx, skipDLQEntryMutation, map[string]any{"id": id}, nil)
}
