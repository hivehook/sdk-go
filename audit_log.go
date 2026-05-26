package hivehook

import (
	"context"
	"time"
)

const auditLogFragment = `
	id actorType actorId actorName action resourceType resourceId orgId ipAddress userAgent details createdAt
`

const listAuditLogsQuery = `query($actorType: String, $resourceType: String, $resourceId: UUID, $action: String, $since: Time, $until: Time, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	auditLogs(actorType: $actorType, resourceType: $resourceType, resourceId: $resourceId, action: $action, since: $since, until: $until, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + auditLogFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getAuditLogQuery = `query($id: UUID!) {
	auditLog(id: $id) {` + auditLogFragment + `}
}`

// AuditLogService queries the audit log.
type AuditLogService struct {
	gql *graphqlClient
}

// ListAuditLogsOptions filters list results.
type ListAuditLogsOptions struct {
	ListOptions
	ActorType    string
	ResourceType string
	ResourceID   string
	Action       string
	Since        *time.Time
	Until        *time.Time
}

// List returns AuditLogs matching opts and a PageInfo cursor.
func (s *AuditLogService) List(ctx context.Context, opts *ListAuditLogsOptions) ([]*AuditLog, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.ActorType != "" {
			vars["actorType"] = opts.ActorType
		}
		if opts.ResourceType != "" {
			vars["resourceType"] = opts.ResourceType
		}
		if opts.ResourceID != "" {
			vars["resourceId"] = opts.ResourceID
		}
		if opts.Action != "" {
			vars["action"] = opts.Action
		}
		if opts.Since != nil {
			vars["since"] = *opts.Since
		}
		if opts.Until != nil {
			vars["until"] = *opts.Until
		}
	}
	var result struct {
		AuditLogs struct {
			Nodes    []*AuditLog `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"auditLogs"`
	}
	if err := s.gql.do(ctx, listAuditLogsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.AuditLogs.Nodes, &result.AuditLogs.PageInfo, nil
}

// Get fetches a single AuditLog by ID.
func (s *AuditLogService) Get(ctx context.Context, id string) (*AuditLog, error) {
	var result struct {
		AuditLog *AuditLog `json:"auditLog"`
	}
	if err := s.gql.do(ctx, getAuditLogQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.AuditLog, nil
}
