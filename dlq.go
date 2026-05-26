package hivehook

import "context"

const dlqFragment = `
	id deliveryId eventId lastError replayedAt createdAt
`

const listDLQQuery = `query($eventId: UUID, $replayed: Boolean, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	dlqEntries(eventId: $eventId, replayed: $replayed, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + dlqFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const replayDLQMutation = `mutation($id: UUID!) {
	replayDLQEntry(id: $id)
}`

const replayAllDLQMutation = `mutation {
	replayAllDLQ { deliveries }
}`

const purgeDLQMutation = `mutation($olderThan: String) {
	purgeDLQ(olderThan: $olderThan) { purged }
}`

// DLQService manages the inbound dead-letter queue.
type DLQService struct {
	gql *graphqlClient
}

// ListDLQOptions filters list results.
type ListDLQOptions struct {
	ListOptions
	EventID  string
	Replayed *bool
}

// List returns DLQ entries matching opts and a PageInfo cursor.
func (s *DLQService) List(ctx context.Context, opts *ListDLQOptions) ([]*DLQEntry, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.EventID != "" {
			vars["eventId"] = opts.EventID
		}
		if opts.Replayed != nil {
			vars["replayed"] = *opts.Replayed
		}
	}
	var result struct {
		DLQEntries struct {
			Nodes    []*DLQEntry `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"dlqEntries"`
	}
	if err := s.gql.do(ctx, listDLQQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.DLQEntries.Nodes, &result.DLQEntries.PageInfo, nil
}

// Replay re-queues a single DLQ entry by ID for delivery.
func (s *DLQService) Replay(ctx context.Context, id string) error {
	return s.gql.do(ctx, replayDLQMutation, map[string]any{"id": id}, nil)
}

// ReplayAll re-queues every entry currently in the DLQ.
func (s *DLQService) ReplayAll(ctx context.Context) (*ReplayResult, error) {
	var result struct {
		ReplayAllDLQ ReplayResult `json:"replayAllDLQ"`
	}
	if err := s.gql.do(ctx, replayAllDLQMutation, nil, &result); err != nil {
		return nil, err
	}
	return &result.ReplayAllDLQ, nil
}

// Purge deletes DLQ entries older than the given duration string.
func (s *DLQService) Purge(ctx context.Context, olderThan string) (*PurgeResult, error) {
	vars := map[string]any{}
	if olderThan != "" {
		vars["olderThan"] = olderThan
	}
	var result struct {
		PurgeDLQ PurgeResult `json:"purgeDLQ"`
	}
	if err := s.gql.do(ctx, purgeDLQMutation, vars, &result); err != nil {
		return nil, err
	}
	return &result.PurgeDLQ, nil
}
