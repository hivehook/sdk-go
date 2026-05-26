package hivehook

import "context"

const deliveryFragment = `
	id eventId subscriptionId destinationId status attempts maxAttempts nextAttemptAt createdAt
`

const listDeliveriesQuery = `query($eventId: UUID, $destinationId: UUID, $subscriptionId: UUID, $status: DeliveryStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	deliveries(eventId: $eventId, destinationId: $destinationId, subscriptionId: $subscriptionId, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + deliveryFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getDeliveryQuery = `query($id: UUID!) {
	delivery(id: $id) {` + deliveryFragment + `
		deliveryAttempts { id deliveryId attemptNumber responseStatus responseBody error durationMs attemptedAt }
	}
}`

// DeliveryService queries inbound Deliveries.
type DeliveryService struct {
	gql *graphqlClient
}

// ListDeliveriesOptions filters list results.
type ListDeliveriesOptions struct {
	ListOptions
	EventID        string
	DestinationID  string
	SubscriptionID string
	Status         DeliveryStatus
}

// List returns Deliveries matching opts and a PageInfo cursor.
func (s *DeliveryService) List(ctx context.Context, opts *ListDeliveriesOptions) ([]*Delivery, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.EventID != "" {
			vars["eventId"] = opts.EventID
		}
		if opts.DestinationID != "" {
			vars["destinationId"] = opts.DestinationID
		}
		if opts.SubscriptionID != "" {
			vars["subscriptionId"] = opts.SubscriptionID
		}
		if opts.Status != "" {
			vars["status"] = opts.Status
		}
	}
	var result struct {
		Deliveries struct {
			Nodes    []*Delivery `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"deliveries"`
	}
	if err := s.gql.do(ctx, listDeliveriesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Deliveries.Nodes, &result.Deliveries.PageInfo, nil
}

// Get fetches a single Delivery by ID.
func (s *DeliveryService) Get(ctx context.Context, id string) (*Delivery, error) {
	var result struct {
		Delivery *Delivery `json:"delivery"`
	}
	if err := s.gql.do(ctx, getDeliveryQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Delivery, nil
}
