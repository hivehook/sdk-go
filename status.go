package hivehook

import "context"

// StatusService reports server health and metadata.
type StatusService struct {
	gql *graphqlClient
}

const getStatusQuery = `
	query GetStatus {
		status {
			status
			dlqSize
			outboundDlqSize
			queueDepth
			activeWorkers
			totalWorkers
			uptime
			version
			sourcesTotal
			destinationsTotal
			subscriptionsTotal
			eventsTotal
			eventsFailed
			deliveriesTotal
			deliveriesPending
			deliveriesDelivered
			messagesTotal
			outboundDeliveriesTotal
			outboundDeliveriesPending
			outboundDeliveriesFailed
		}
	}
`

// Get fetches a single SystemStatus by ID.
func (s *StatusService) Get(ctx context.Context) (*SystemStatus, error) {
	var result struct {
		Status *SystemStatus `json:"status"`
	}
	if err := s.gql.do(ctx, getStatusQuery, nil, &result); err != nil {
		return nil, err
	}
	return result.Status, nil
}
