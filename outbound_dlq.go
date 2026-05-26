package hivehook

import "context"

// OutboundDLQService manages the outbound dead-letter queue.
type OutboundDLQService struct {
	gql *graphqlClient
}

// ListOutboundDLQOptions filters list results.
type ListOutboundDLQOptions struct {
	ListOptions
	MessageID *string
	Replayed  *bool
}

const outboundDLQFragment = `
	fragment OutboundDLQFields on OutboundDLQEntry {
		id deliveryId messageId lastError replayedAt createdAt
	}
`

const listOutboundDLQQuery = `
	query ListOutboundDLQ($messageId: UUID, $replayed: Boolean, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
		outboundDlqEntries(messageId: $messageId, replayed: $replayed, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
			nodes { ...OutboundDLQFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + outboundDLQFragment

const replayOutboundDLQMutation = `
	mutation ReplayOutboundDLQ($id: UUID!) {
		replayOutboundDlqEntry(id: $id)
	}
`

const replayAllOutboundDLQMutation = `
	mutation ReplayAllOutboundDLQ {
		replayAllOutboundDlq { deliveries }
	}
`

const purgeOutboundDLQMutation = `
	mutation PurgeOutboundDLQ($olderThan: String) {
		purgeOutboundDlq(olderThan: $olderThan) { purged }
	}
`

// List returns OutboundDLQ entries matching opts and a PageInfo cursor.
func (s *OutboundDLQService) List(ctx context.Context, opts *ListOutboundDLQOptions) ([]*OutboundDLQEntry, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
		if opts.MessageID != nil {
			vars = mergeVars(vars, map[string]any{"messageId": *opts.MessageID})
		}
		if opts.Replayed != nil {
			vars = mergeVars(vars, map[string]any{"replayed": *opts.Replayed})
		}
	}
	var result struct {
		OutboundDlqEntries struct {
			Nodes    []*OutboundDLQEntry `json:"nodes"`
			PageInfo PageInfo            `json:"pageInfo"`
		} `json:"outboundDlqEntries"`
	}
	if err := s.gql.do(ctx, listOutboundDLQQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.OutboundDlqEntries.Nodes, &result.OutboundDlqEntries.PageInfo, nil
}

// Replay re-queues a single DLQ entry by ID for delivery.
func (s *OutboundDLQService) Replay(ctx context.Context, id string) error {
	vars := map[string]any{"id": id}
	return s.gql.do(ctx, replayOutboundDLQMutation, vars, nil)
}

// ReplayAll re-queues every entry currently in the DLQ.
func (s *OutboundDLQService) ReplayAll(ctx context.Context) (*ReplayResult, error) {
	var result struct {
		ReplayAllOutboundDlq *ReplayResult `json:"replayAllOutboundDlq"`
	}
	if err := s.gql.do(ctx, replayAllOutboundDLQMutation, nil, &result); err != nil {
		return nil, err
	}
	return result.ReplayAllOutboundDlq, nil
}

// Purge deletes DLQ entries older than the given duration string.
func (s *OutboundDLQService) Purge(ctx context.Context, olderThan string) (*PurgeResult, error) {
	vars := map[string]any{"olderThan": olderThan}
	var result struct {
		PurgeOutboundDlq *PurgeResult `json:"purgeOutboundDlq"`
	}
	if err := s.gql.do(ctx, purgeOutboundDLQMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.PurgeOutboundDlq, nil
}
