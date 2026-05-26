package hivehook

import "context"

const eventFragment = `
	id sourceId idempotencyKey eventType headers rawBody status receivedAt
`

const listEventsQuery = `query($sourceId: UUID, $eventType: String, $status: EventStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	events(sourceId: $sourceId, eventType: $eventType, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + eventFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getEventQuery = `query($id: UUID!) {
	event(id: $id) {` + eventFragment + `}
}`

// EventService queries ingested Events.
type EventService struct {
	gql *graphqlClient
}

// ListEventsOptions filters list results.
type ListEventsOptions struct {
	ListOptions
	SourceID  string
	EventType string
	Status    EventStatus
}

// List returns Events matching opts and a PageInfo cursor.
func (s *EventService) List(ctx context.Context, opts *ListEventsOptions) ([]*Event, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.SourceID != "" {
			vars["sourceId"] = opts.SourceID
		}
		if opts.EventType != "" {
			vars["eventType"] = opts.EventType
		}
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
	}
	var result struct {
		Events struct {
			Nodes    []*Event `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"events"`
	}
	if err := s.gql.do(ctx, listEventsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Events.Nodes, &result.Events.PageInfo, nil
}

// Get fetches a single Event by ID.
func (s *EventService) Get(ctx context.Context, id string) (*Event, error) {
	var result struct {
		Event *Event `json:"event"`
	}
	if err := s.gql.do(ctx, getEventQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Event, nil
}
