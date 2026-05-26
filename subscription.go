package hivehook

import "context"

const subscriptionFragment = `
	id name sourceId destinationId filterConfig { eventTypes regex bodyMatch { path value operator } rules { path operator value rules { path operator value } } } transformConfig { envelope headers } enabled createdAt
`

const listSubscriptionsQuery = `query($sourceId: UUID, $destinationId: UUID, $enabled: Boolean, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	subscriptions(sourceId: $sourceId, destinationId: $destinationId, enabled: $enabled, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + subscriptionFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getSubscriptionQuery = `query($id: UUID!) {
	subscription(id: $id) {` + subscriptionFragment + `}
}`

const createSubscriptionMutation = `mutation($input: CreateSubscriptionInput!) {
	createSubscription(input: $input) {` + subscriptionFragment + `}
}`

const updateSubscriptionMutation = `mutation($id: UUID!, $input: UpdateSubscriptionInput!) {
	updateSubscription(id: $id, input: $input) {` + subscriptionFragment + `}
}`

const deleteSubscriptionMutation = `mutation($id: UUID!) {
	deleteSubscription(id: $id)
}`

// SubscriptionService manages Subscription routing rules.
type SubscriptionService struct {
	gql *graphqlClient
}

// ListSubscriptionsOptions filters list results.
type ListSubscriptionsOptions struct {
	ListOptions
	SourceID      string
	DestinationID string
	Enabled       *bool
}

// CreateSubscriptionInput is the input for the matching Create method.
type CreateSubscriptionInput struct {
	Name            string           `json:"name"`
	SourceID        string           `json:"sourceId"`
	DestinationID   string           `json:"destinationId"`
	FilterConfig    *FilterConfig    `json:"filterConfig,omitempty"`
	TransformConfig *TransformConfig `json:"transformConfig,omitempty"`
}

// UpdateSubscriptionInput is the input for the matching Update method.
type UpdateSubscriptionInput struct {
	Name            *string          `json:"name,omitempty"`
	SourceID        *string          `json:"sourceId,omitempty"`
	DestinationID   *string          `json:"destinationId,omitempty"`
	FilterConfig    *FilterConfig    `json:"filterConfig,omitempty"`
	TransformConfig *TransformConfig `json:"transformConfig,omitempty"`
	Enabled         *bool            `json:"enabled,omitempty"`
}

// List returns Subscriptions matching opts and a PageInfo cursor.
func (s *SubscriptionService) List(ctx context.Context, opts *ListSubscriptionsOptions) ([]*Subscription, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.SourceID != "" {
			vars["sourceId"] = opts.SourceID
		}
		if opts.DestinationID != "" {
			vars["destinationId"] = opts.DestinationID
		}
		if opts.Enabled != nil {
			vars["enabled"] = *opts.Enabled
		}
	}
	var result struct {
		Subscriptions struct {
			Nodes    []*Subscription `json:"nodes"`
			PageInfo PageInfo        `json:"pageInfo"`
		} `json:"subscriptions"`
	}
	if err := s.gql.do(ctx, listSubscriptionsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Subscriptions.Nodes, &result.Subscriptions.PageInfo, nil
}

// Get fetches a single Subscription by ID.
func (s *SubscriptionService) Get(ctx context.Context, id string) (*Subscription, error) {
	var result struct {
		Subscription *Subscription `json:"subscription"`
	}
	if err := s.gql.do(ctx, getSubscriptionQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Subscription, nil
}

// Create persists a new Subscription.
func (s *SubscriptionService) Create(ctx context.Context, input *CreateSubscriptionInput) (*Subscription, error) {
	m := map[string]any{
		"name":          input.Name,
		"sourceId":      input.SourceID,
		"destinationId": input.DestinationID,
	}
	if input.FilterConfig != nil {
		m["filterConfig"] = buildFilterConfigMap(input.FilterConfig)
	}
	if input.TransformConfig != nil {
		m["transformConfig"] = buildTransformConfigMap(input.TransformConfig)
	}
	var result struct {
		CreateSubscription Subscription `json:"createSubscription"`
	}
	if err := s.gql.do(ctx, createSubscriptionMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateSubscription, nil
}

// Update partially updates an existing Subscription by ID.
func (s *SubscriptionService) Update(ctx context.Context, id string, input *UpdateSubscriptionInput) (*Subscription, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.SourceID != nil {
		m["sourceId"] = *input.SourceID
	}
	if input.DestinationID != nil {
		m["destinationId"] = *input.DestinationID
	}
	if input.FilterConfig != nil {
		m["filterConfig"] = buildFilterConfigMap(input.FilterConfig)
	}
	if input.TransformConfig != nil {
		m["transformConfig"] = buildTransformConfigMap(input.TransformConfig)
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	var result struct {
		UpdateSubscription Subscription `json:"updateSubscription"`
	}
	if err := s.gql.do(ctx, updateSubscriptionMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateSubscription, nil
}

// Delete removes a Subscription by ID.
func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteSubscriptionMutation, map[string]any{"id": id}, nil)
}

func buildFilterConfigMap(fc *FilterConfig) map[string]any {
	m := map[string]any{}
	if len(fc.EventTypes) > 0 {
		m["eventTypes"] = fc.EventTypes
	}
	if len(fc.Regex) > 0 {
		m["regex"] = fc.Regex
	}
	if len(fc.BodyMatch) > 0 {
		rules := make([]map[string]any, len(fc.BodyMatch))
		for i, r := range fc.BodyMatch {
			rules[i] = map[string]any{
				"path":     r.Path,
				"value":    r.Value,
				"operator": r.Operator,
			}
		}
		m["bodyMatch"] = rules
	}
	if len(fc.Rules) > 0 {
		m["rules"] = buildFilterRules(fc.Rules)
	}
	return m
}

func buildFilterRules(rules []FilterRule) []map[string]any {
	out := make([]map[string]any, len(rules))
	for i, r := range rules {
		rm := map[string]any{
			"operator": r.Operator,
		}
		if r.Path != "" {
			rm["path"] = r.Path
		}
		if r.Value != nil {
			rm["value"] = r.Value
		}
		if len(r.Rules) > 0 {
			rm["rules"] = buildFilterRules(r.Rules)
		}
		out[i] = rm
	}
	return out
}

func buildTransformConfigMap(tc *TransformConfig) map[string]any {
	m := map[string]any{
		"envelope": tc.Envelope,
	}
	if tc.Headers != nil {
		m["headers"] = tc.Headers
	}
	return m
}
