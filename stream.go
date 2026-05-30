package hivehook

import "context"

const streamFragment = `
	id applicationId name status retentionDays createdAt
`

const listStreamsQuery = `query($applicationId: UUID, $status: StreamStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	streams(applicationId: $applicationId, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + streamFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getStreamQuery = `query($id: UUID!) {
	stream(id: $id) {` + streamFragment + `}
}`

const createStreamMutation = `mutation($input: CreateStreamInput!) {
	createStream(input: $input) {` + streamFragment + `}
}`

const updateStreamMutation = `mutation($id: UUID!, $input: UpdateStreamInput!) {
	updateStream(id: $id, input: $input) {` + streamFragment + `}
}`

const deleteStreamMutation = `mutation($id: UUID!) {
	deleteStream(id: $id)
}`

// StreamService manages event Streams.
type StreamService struct {
	gql *graphqlClient
}

// ListStreamsOptions filters list results.
type ListStreamsOptions struct {
	ListOptions
	ApplicationID string
	Status        StreamStatus
}

// CreateStreamInput is the input for the matching Create method.
type CreateStreamInput struct {
	ApplicationID string `json:"applicationId"`
	Name          string `json:"name"`
	RetentionDays *int   `json:"retentionDays,omitempty"`
}

// UpdateStreamInput is the input for the matching Update method.
type UpdateStreamInput struct {
	Name          *string       `json:"name,omitempty"`
	RetentionDays *int          `json:"retentionDays,omitempty"`
	Status        *StreamStatus `json:"status,omitempty"`
}

// List returns Streams matching opts and a PageInfo cursor.
func (s *StreamService) List(ctx context.Context, opts *ListStreamsOptions) ([]*Stream, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.ApplicationID != "" {
			vars["applicationId"] = opts.ApplicationID
		}
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
	}
	var result struct {
		Streams struct {
			Nodes    []*Stream `json:"nodes"`
			PageInfo PageInfo  `json:"pageInfo"`
		} `json:"streams"`
	}
	if err := s.gql.do(ctx, listStreamsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Streams.Nodes, &result.Streams.PageInfo, nil
}

// Get fetches a single Stream by ID.
func (s *StreamService) Get(ctx context.Context, id string) (*Stream, error) {
	var result struct {
		Stream *Stream `json:"stream"`
	}
	if err := s.gql.do(ctx, getStreamQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Stream, nil
}

// Create persists a new Stream.
func (s *StreamService) Create(ctx context.Context, input *CreateStreamInput) (*Stream, error) {
	m := map[string]any{
		"applicationId": input.ApplicationID,
		"name":          input.Name,
	}
	if input.RetentionDays != nil {
		m["retentionDays"] = *input.RetentionDays
	}
	var result struct {
		CreateStream Stream `json:"createStream"`
	}
	if err := s.gql.do(ctx, createStreamMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateStream, nil
}

// Update partially updates an existing Stream by ID.
func (s *StreamService) Update(ctx context.Context, id string, input *UpdateStreamInput) (*Stream, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.RetentionDays != nil {
		m["retentionDays"] = *input.RetentionDays
	}
	if input.Status != nil {
		m["status"] = *input.Status
	}
	var result struct {
		UpdateStream Stream `json:"updateStream"`
	}
	if err := s.gql.do(ctx, updateStreamMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateStream, nil
}

// Delete removes a Stream by ID.
func (s *StreamService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteStreamMutation, map[string]any{"id": id}, nil)
}

const streamEntryFragment = `
	id streamId sequence messageId eventType payload createdAt
`

const listStreamEntriesQuery = `query($streamId: UUID!, $afterSequence: Int, $limit: Int) {
	streamEntries(streamId: $streamId, afterSequence: $afterSequence, limit: $limit) {
		nodes {` + streamEntryFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

// ListStreamEntriesOptions filters StreamService.Entries results.
type ListStreamEntriesOptions struct {
	AfterSequence *int
	Limit         *int
}

// Entries returns persisted entries from a Stream, ordered by sequence.
func (s *StreamService) Entries(ctx context.Context, streamID string, opts *ListStreamEntriesOptions) ([]*StreamEntry, *PageInfo, error) {
	vars := map[string]any{"streamId": streamID}
	if opts != nil {
		if opts.AfterSequence != nil {
			vars["afterSequence"] = *opts.AfterSequence
		}
		if opts.Limit != nil {
			vars["limit"] = *opts.Limit
		}
	}
	var result struct {
		StreamEntries struct {
			Nodes    []*StreamEntry `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"streamEntries"`
	}
	if err := s.gql.do(ctx, listStreamEntriesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.StreamEntries.Nodes, &result.StreamEntries.PageInfo, nil
}
