package hivehook

import "context"

const eventTypeSchemaFragment = `
	id eventType description schema example createdAt updatedAt
`

const listEventTypeSchemasQuery = `query($search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	eventTypeSchemas(search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + eventTypeSchemaFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getEventTypeSchemaQuery = `query($id: UUID!) {
	eventTypeSchema(id: $id) {` + eventTypeSchemaFragment + `}
}`

const createEventTypeSchemaMutation = `mutation($input: CreateEventTypeSchemaInput!) {
	createEventTypeSchema(input: $input) {` + eventTypeSchemaFragment + `}
}`

const updateEventTypeSchemaMutation = `mutation($id: UUID!, $input: UpdateEventTypeSchemaInput!) {
	updateEventTypeSchema(id: $id, input: $input) {` + eventTypeSchemaFragment + `}
}`

const deleteEventTypeSchemaMutation = `mutation($id: UUID!) {
	deleteEventTypeSchema(id: $id)
}`

// EventTypeSchemaService manages event-type JSON schemas.
type EventTypeSchemaService struct {
	gql *graphqlClient
}

// CreateEventTypeSchemaInput is the input for the matching Create method.
type CreateEventTypeSchemaInput struct {
	EventType   string         `json:"eventType"`
	Description *string        `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
	Example     map[string]any `json:"example,omitempty"`
}

// UpdateEventTypeSchemaInput is the input for the matching Update method.
type UpdateEventTypeSchemaInput struct {
	EventType   *string        `json:"eventType,omitempty"`
	Description *string        `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
	Example     map[string]any `json:"example,omitempty"`
}

// List returns EventTypeSchemas matching opts and a PageInfo cursor.
func (s *EventTypeSchemaService) List(ctx context.Context, opts *ListOptions) ([]*EventTypeSchema, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
	}
	var result struct {
		EventTypeSchemas struct {
			Nodes    []*EventTypeSchema `json:"nodes"`
			PageInfo PageInfo           `json:"pageInfo"`
		} `json:"eventTypeSchemas"`
	}
	if err := s.gql.do(ctx, listEventTypeSchemasQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.EventTypeSchemas.Nodes, &result.EventTypeSchemas.PageInfo, nil
}

// Get fetches a single EventTypeSchema by ID.
func (s *EventTypeSchemaService) Get(ctx context.Context, id string) (*EventTypeSchema, error) {
	var result struct {
		EventTypeSchema *EventTypeSchema `json:"eventTypeSchema"`
	}
	if err := s.gql.do(ctx, getEventTypeSchemaQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.EventTypeSchema, nil
}

// Create persists a new EventTypeSchema.
func (s *EventTypeSchemaService) Create(ctx context.Context, input *CreateEventTypeSchemaInput) (*EventTypeSchema, error) {
	m := map[string]any{
		"eventType": input.EventType,
	}
	if input.Description != nil {
		m["description"] = *input.Description
	}
	if input.Schema != nil {
		m["schema"] = input.Schema
	}
	if input.Example != nil {
		m["example"] = input.Example
	}
	var result struct {
		CreateEventTypeSchema EventTypeSchema `json:"createEventTypeSchema"`
	}
	if err := s.gql.do(ctx, createEventTypeSchemaMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateEventTypeSchema, nil
}

// Update partially updates an existing EventTypeSchema by ID.
func (s *EventTypeSchemaService) Update(ctx context.Context, id string, input *UpdateEventTypeSchemaInput) (*EventTypeSchema, error) {
	m := map[string]any{}
	if input.EventType != nil {
		m["eventType"] = *input.EventType
	}
	if input.Description != nil {
		m["description"] = *input.Description
	}
	if input.Schema != nil {
		m["schema"] = input.Schema
	}
	if input.Example != nil {
		m["example"] = input.Example
	}
	var result struct {
		UpdateEventTypeSchema EventTypeSchema `json:"updateEventTypeSchema"`
	}
	if err := s.gql.do(ctx, updateEventTypeSchemaMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateEventTypeSchema, nil
}

// Delete removes a EventTypeSchema by ID.
func (s *EventTypeSchemaService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteEventTypeSchemaMutation, map[string]any{"id": id}, nil)
}
