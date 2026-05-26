package hivehook

import (
	"context"
	"time"
)

const apiKeyFragment = `
	id name keyPrefix scopes sourceIds createdAt expiresAt revokedAt lastUsedAt
`

const listAPIKeysQuery = `query($search: String, $limit: Int, $offset: Int) {
	apiKeys(search: $search, limit: $limit, offset: $offset) {
		nodes {` + apiKeyFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getAPIKeyQuery = `query($id: UUID!) {
	apiKey(id: $id) {` + apiKeyFragment + `}
}`

const createAPIKeyMutation = `mutation($input: CreateAPIKeyInput!) {
	createAPIKey(input: $input) {
		apiKey {` + apiKeyFragment + `}
		rawKey
	}
}`

const revokeAPIKeyMutation = `mutation($id: UUID!) {
	revokeAPIKey(id: $id)
}`

// APIKeyService manages API keys.
type APIKeyService struct {
	gql *graphqlClient
}

// CreateAPIKeyInput is the input for the matching Create method.
type CreateAPIKeyInput struct {
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes,omitempty"`
	SourceIDs []string   `json:"sourceIds,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// List returns APIKeys matching opts and a PageInfo cursor.
func (s *APIKeyService) List(ctx context.Context, opts *ListOptions) ([]*APIKey, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
	}
	var result struct {
		APIKeys struct {
			Nodes    []*APIKey `json:"nodes"`
			PageInfo PageInfo  `json:"pageInfo"`
		} `json:"apiKeys"`
	}
	if err := s.gql.do(ctx, listAPIKeysQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.APIKeys.Nodes, &result.APIKeys.PageInfo, nil
}

// Get fetches a single APIKey by ID.
func (s *APIKeyService) Get(ctx context.Context, id string) (*APIKey, error) {
	var result struct {
		APIKey *APIKey `json:"apiKey"`
	}
	if err := s.gql.do(ctx, getAPIKeyQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.APIKey, nil
}

// Create persists a new APIKey.
func (s *APIKeyService) Create(ctx context.Context, input *CreateAPIKeyInput) (*APIKeyWithSecret, error) {
	m := map[string]any{
		"name": input.Name,
	}
	if len(input.Scopes) > 0 {
		m["scopes"] = input.Scopes
	}
	if len(input.SourceIDs) > 0 {
		m["sourceIds"] = input.SourceIDs
	}
	if input.ExpiresAt != nil {
		m["expiresAt"] = input.ExpiresAt.Format(time.RFC3339)
	}
	var result struct {
		CreateAPIKey APIKeyWithSecret `json:"createAPIKey"`
	}
	if err := s.gql.do(ctx, createAPIKeyMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateAPIKey, nil
}

// Revoke revokes an APIKey by ID.
func (s *APIKeyService) Revoke(ctx context.Context, id string) error {
	return s.gql.do(ctx, revokeAPIKeyMutation, map[string]any{"id": id}, nil)
}
