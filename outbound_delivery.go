package hivehook

import "context"

// OutboundDeliveryService queries OutboundDeliveries.
type OutboundDeliveryService struct {
	gql *graphqlClient
}

// ListOutboundDeliveriesOptions filters list results.
type ListOutboundDeliveriesOptions struct {
	ListOptions
	MessageID  *string
	EndpointID *string
	Status     DeliveryStatus
}

const outboundDeliveryFragment = `
	fragment OutboundDeliveryFields on OutboundDelivery {
		id messageId endpointId status attempts maxAttempts nextAttemptAt createdAt
	}
`

const listOutboundDeliveriesQuery = `
	query ListOutboundDeliveries($messageId: UUID, $endpointId: UUID, $status: DeliveryStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
		outboundDeliveries(messageId: $messageId, endpointId: $endpointId, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
			nodes { ...OutboundDeliveryFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + outboundDeliveryFragment

const getOutboundDeliveryQuery = `
	query GetOutboundDelivery($id: UUID!) {
		outboundDelivery(id: $id) {
			...OutboundDeliveryFields
			deliveryAttempts {
				id deliveryId attemptNumber responseStatus responseBody error durationMs attemptedAt
			}
		}
	}
` + outboundDeliveryFragment

// List returns OutboundDeliveries matching opts and a PageInfo cursor.
func (s *OutboundDeliveryService) List(ctx context.Context, opts *ListOutboundDeliveriesOptions) ([]*OutboundDelivery, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
		if opts.MessageID != nil {
			vars = mergeVars(vars, map[string]any{"messageId": *opts.MessageID})
		}
		if opts.EndpointID != nil {
			vars = mergeVars(vars, map[string]any{"endpointId": *opts.EndpointID})
		}
		if opts.Status != "" {
			vars = mergeVars(vars, map[string]any{"status": opts.Status})
		}
	}
	var result struct {
		OutboundDeliveries struct {
			Nodes    []*OutboundDelivery `json:"nodes"`
			PageInfo PageInfo            `json:"pageInfo"`
		} `json:"outboundDeliveries"`
	}
	if err := s.gql.do(ctx, listOutboundDeliveriesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.OutboundDeliveries.Nodes, &result.OutboundDeliveries.PageInfo, nil
}

// Get fetches a single OutboundDelivery by ID.
func (s *OutboundDeliveryService) Get(ctx context.Context, id string) (*OutboundDelivery, error) {
	vars := map[string]any{"id": id}
	var result struct {
		OutboundDelivery *OutboundDelivery `json:"outboundDelivery"`
	}
	if err := s.gql.do(ctx, getOutboundDeliveryQuery, vars, &result); err != nil {
		return nil, err
	}
	return result.OutboundDelivery, nil
}
