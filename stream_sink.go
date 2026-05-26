package hivehook

import "context"

const streamSinkFragment = `
	id streamId name sinkType config batchSize flushInterval cursorSequence status lastFlushedAt createdAt
`

const listStreamSinksQuery = `query($streamId: UUID!, $status: SinkStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	streamSinks(streamId: $streamId, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + streamSinkFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getStreamSinkQuery = `query($id: UUID!) {
	streamSink(id: $id) {` + streamSinkFragment + `}
}`

const createStreamSinkMutation = `mutation($input: CreateStreamSinkInput!) {
	createStreamSink(input: $input) {` + streamSinkFragment + `}
}`

const updateStreamSinkMutation = `mutation($id: UUID!, $input: UpdateStreamSinkInput!) {
	updateStreamSink(id: $id, input: $input) {` + streamSinkFragment + `}
}`

const deleteStreamSinkMutation = `mutation($id: UUID!) {
	deleteStreamSink(id: $id)
}`

// StreamSinkService manages StreamSinks (S3, Postgres, webhook).
type StreamSinkService struct {
	gql *graphqlClient
}

// ListStreamSinksOptions filters list results.
type ListStreamSinksOptions struct {
	ListOptions
	Status SinkStatus
}

// CreateStreamSinkInput is the input for the matching Create method.
type CreateStreamSinkInput struct {
	StreamID      string         `json:"streamId"`
	Name          string         `json:"name"`
	SinkType      SinkType       `json:"sinkType"`
	Config        map[string]any `json:"config,omitempty"`
	BatchSize     *int           `json:"batchSize,omitempty"`
	FlushInterval *string        `json:"flushInterval,omitempty"`
}

// UpdateStreamSinkInput is the input for the matching Update method.
type UpdateStreamSinkInput struct {
	Config        map[string]any `json:"config,omitempty"`
	BatchSize     *int           `json:"batchSize,omitempty"`
	FlushInterval *string        `json:"flushInterval,omitempty"`
	Status        *SinkStatus    `json:"status,omitempty"`
}

// List returns StreamSinks matching opts and a PageInfo cursor.
func (s *StreamSinkService) List(ctx context.Context, streamID string, opts *ListStreamSinksOptions) ([]*StreamSink, *PageInfo, error) {
	vars := map[string]any{"streamId": streamID}
	if opts != nil {
		for k, v := range opts.toVars() {
			vars[k] = v
		}
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
	}
	var result struct {
		StreamSinks struct {
			Nodes    []*StreamSink `json:"nodes"`
			PageInfo PageInfo      `json:"pageInfo"`
		} `json:"streamSinks"`
	}
	if err := s.gql.do(ctx, listStreamSinksQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.StreamSinks.Nodes, &result.StreamSinks.PageInfo, nil
}

// Get fetches a single StreamSink by ID.
func (s *StreamSinkService) Get(ctx context.Context, id string) (*StreamSink, error) {
	var result struct {
		StreamSink *StreamSink `json:"streamSink"`
	}
	if err := s.gql.do(ctx, getStreamSinkQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.StreamSink, nil
}

// Create persists a new StreamSink.
func (s *StreamSinkService) Create(ctx context.Context, input *CreateStreamSinkInput) (*StreamSink, error) {
	m := map[string]any{
		"streamId": input.StreamID,
		"name":     input.Name,
		"sinkType": input.SinkType,
	}
	if input.Config != nil {
		m["config"] = input.Config
	}
	if input.BatchSize != nil {
		m["batchSize"] = *input.BatchSize
	}
	if input.FlushInterval != nil {
		m["flushInterval"] = *input.FlushInterval
	}
	var result struct {
		CreateStreamSink StreamSink `json:"createStreamSink"`
	}
	if err := s.gql.do(ctx, createStreamSinkMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateStreamSink, nil
}

// Update partially updates an existing StreamSink by ID.
func (s *StreamSinkService) Update(ctx context.Context, id string, input *UpdateStreamSinkInput) (*StreamSink, error) {
	m := map[string]any{}
	if input.Config != nil {
		m["config"] = input.Config
	}
	if input.BatchSize != nil {
		m["batchSize"] = *input.BatchSize
	}
	if input.FlushInterval != nil {
		m["flushInterval"] = *input.FlushInterval
	}
	if input.Status != nil {
		m["status"] = *input.Status
	}
	var result struct {
		UpdateStreamSink StreamSink `json:"updateStreamSink"`
	}
	if err := s.gql.do(ctx, updateStreamSinkMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateStreamSink, nil
}

// Delete removes a StreamSink by ID.
func (s *StreamSinkService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteStreamSinkMutation, map[string]any{"id": id}, nil)
}
