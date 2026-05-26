package hivehook

import "context"

const streamConsumerFragment = `
	id streamId name cursorSequence createdAt updatedAt
`

const listStreamConsumersQuery = `query($streamId: UUID!, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	streamConsumers(streamId: $streamId, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + streamConsumerFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getStreamConsumerQuery = `query($id: UUID!) {
	streamConsumer(id: $id) {` + streamConsumerFragment + `}
}`

const createStreamConsumerMutation = `mutation($input: CreateStreamConsumerInput!) {
	createStreamConsumer(input: $input) {` + streamConsumerFragment + `}
}`

const deleteStreamConsumerMutation = `mutation($id: UUID!) {
	deleteStreamConsumer(id: $id)
}`

const advanceConsumerCursorMutation = `mutation($id: UUID!, $sequence: Int!) {
	advanceConsumerCursor(id: $id, sequence: $sequence) {` + streamConsumerFragment + `}
}`

// StreamConsumerService manages StreamConsumer cursors.
type StreamConsumerService struct {
	gql *graphqlClient
}

// ListStreamConsumersOptions filters list results.
type ListStreamConsumersOptions struct {
	ListOptions
}

// CreateStreamConsumerInput is the input for the matching Create method.
type CreateStreamConsumerInput struct {
	StreamID string `json:"streamId"`
	Name     string `json:"name"`
}

// List returns StreamConsumers matching opts and a PageInfo cursor.
func (s *StreamConsumerService) List(ctx context.Context, streamID string, opts *ListStreamConsumersOptions) ([]*StreamConsumer, *PageInfo, error) {
	vars := map[string]any{"streamId": streamID}
	if opts != nil {
		for k, v := range opts.toVars() {
			vars[k] = v
		}
	}
	var result struct {
		StreamConsumers struct {
			Nodes    []*StreamConsumer `json:"nodes"`
			PageInfo PageInfo          `json:"pageInfo"`
		} `json:"streamConsumers"`
	}
	if err := s.gql.do(ctx, listStreamConsumersQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.StreamConsumers.Nodes, &result.StreamConsumers.PageInfo, nil
}

// Get fetches a single StreamConsumer by ID.
func (s *StreamConsumerService) Get(ctx context.Context, id string) (*StreamConsumer, error) {
	var result struct {
		StreamConsumer *StreamConsumer `json:"streamConsumer"`
	}
	if err := s.gql.do(ctx, getStreamConsumerQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.StreamConsumer, nil
}

// Create persists a new StreamConsumer.
func (s *StreamConsumerService) Create(ctx context.Context, input *CreateStreamConsumerInput) (*StreamConsumer, error) {
	m := map[string]any{
		"streamId": input.StreamID,
		"name":     input.Name,
	}
	var result struct {
		CreateStreamConsumer StreamConsumer `json:"createStreamConsumer"`
	}
	if err := s.gql.do(ctx, createStreamConsumerMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateStreamConsumer, nil
}

// Delete removes a StreamConsumer by ID.
func (s *StreamConsumerService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteStreamConsumerMutation, map[string]any{"id": id}, nil)
}

// AdvanceCursor moves a StreamConsumer's cursor to the given sequence.
func (s *StreamConsumerService) AdvanceCursor(ctx context.Context, id string, sequence int) (*StreamConsumer, error) {
	var result struct {
		AdvanceConsumerCursor StreamConsumer `json:"advanceConsumerCursor"`
	}
	if err := s.gql.do(ctx, advanceConsumerCursorMutation, map[string]any{"id": id, "sequence": sequence}, &result); err != nil {
		return nil, err
	}
	return &result.AdvanceConsumerCursor, nil
}
