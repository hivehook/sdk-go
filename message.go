package hivehook

import "context"

// MessageService sends and queries outbound Messages.
type MessageService struct {
	gql *graphqlClient
}

// ListMessagesOptions filters list results.
type ListMessagesOptions struct {
	ListOptions
	ApplicationID *string
	EventType     *string
	Status        MessageStatus
}

// SendMessageInput is the input for the matching service method.
type SendMessageInput struct {
	ApplicationID  string
	EventType      string
	Payload        []byte
	IdempotencyKey *string
	Broadcast      *bool
}

// BroadcastMessageInput is the input for the matching service method.
type BroadcastMessageInput struct {
	ApplicationID  string
	EventType      string
	Payload        []byte
	IdempotencyKey *string
}

// SendDynamicMessageInput is the input for the matching service method.
type SendDynamicMessageInput struct {
	URL            string
	EventType      string
	Payload        []byte
	Headers        map[string]any
	SigningSecret  *string
	RetryPolicy    *RetryPolicy
	TimeoutMs      *int
	RateLimitRps   *int
	IdempotencyKey *string
}

const messageFragment = `
	fragment MessageFields on Message {
		id applicationId eventType payload idempotencyKey status createdAt
	}
`

const listMessagesQuery = `
	query ListMessages($applicationId: UUID, $eventType: String, $status: MessageStatus, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
		messages(applicationId: $applicationId, eventType: $eventType, status: $status, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
			nodes { ...MessageFields }
			pageInfo { total limit offset endCursor hasNextPage }
		}
	}
` + messageFragment

const getMessageQuery = `
	query GetMessage($id: UUID!) {
		message(id: $id) { ...MessageFields }
	}
` + messageFragment

const sendMessageMutation = `
	mutation SendMessage($input: SendMessageInput!) {
		sendMessage(input: $input) { ...MessageFields }
	}
` + messageFragment

const broadcastMessageMutation = `
	mutation BroadcastMessage($input: BroadcastMessageInput!) {
		broadcastMessage(input: $input) { ...MessageFields }
	}
` + messageFragment

const sendDynamicMessageMutation = `
	mutation SendDynamicMessage($input: SendDynamicMessageInput!) {
		sendDynamicMessage(input: $input) {
			id messageId endpointId status attempts maxAttempts nextAttemptAt createdAt
		}
	}
`

// List returns Messages matching opts and a PageInfo cursor.
func (s *MessageService) List(ctx context.Context, opts *ListMessagesOptions) ([]*Message, *PageInfo, error) {
	var vars map[string]any
	if opts != nil {
		vars = opts.toVars()
		if opts.ApplicationID != nil {
			vars = mergeVars(vars, map[string]any{"applicationId": *opts.ApplicationID})
		}
		if opts.EventType != nil {
			vars = mergeVars(vars, map[string]any{"eventType": *opts.EventType})
		}
		if opts.Status != "" {
			vars = mergeVars(vars, map[string]any{"status": opts.Status})
		}
	}
	var result struct {
		Messages struct {
			Nodes    []*Message `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"messages"`
	}
	if err := s.gql.do(ctx, listMessagesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Messages.Nodes, &result.Messages.PageInfo, nil
}

// Get fetches a single Message by ID.
func (s *MessageService) Get(ctx context.Context, id string) (*Message, error) {
	vars := map[string]any{"id": id}
	var result struct {
		Message *Message `json:"message"`
	}
	if err := s.gql.do(ctx, getMessageQuery, vars, &result); err != nil {
		return nil, err
	}
	return result.Message, nil
}

// Send sends a Message.
func (s *MessageService) Send(ctx context.Context, input *SendMessageInput) (*Message, error) {
	inp := map[string]any{
		"applicationId": input.ApplicationID,
		"eventType":     input.EventType,
		"payload":       input.Payload,
	}
	if input.IdempotencyKey != nil {
		inp["idempotencyKey"] = *input.IdempotencyKey
	}
	if input.Broadcast != nil {
		inp["broadcast"] = *input.Broadcast
	}
	vars := map[string]any{"input": inp}
	var result struct {
		SendMessage *Message `json:"sendMessage"`
	}
	if err := s.gql.do(ctx, sendMessageMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.SendMessage, nil
}

// Broadcast broadcasts a Message to all matching Endpoints.
func (s *MessageService) Broadcast(ctx context.Context, input *BroadcastMessageInput) (*Message, error) {
	inp := map[string]any{
		"applicationId": input.ApplicationID,
		"eventType":     input.EventType,
		"payload":       input.Payload,
	}
	if input.IdempotencyKey != nil {
		inp["idempotencyKey"] = *input.IdempotencyKey
	}
	vars := map[string]any{"input": inp}
	var result struct {
		BroadcastMessage *Message `json:"broadcastMessage"`
	}
	if err := s.gql.do(ctx, broadcastMessageMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.BroadcastMessage, nil
}

// SendDynamic sends a Message to a dynamically-provided URL.
func (s *MessageService) SendDynamic(ctx context.Context, input *SendDynamicMessageInput) (*OutboundDelivery, error) {
	inp := map[string]any{
		"url":       input.URL,
		"eventType": input.EventType,
		"payload":   input.Payload,
	}
	if input.Headers != nil {
		inp["headers"] = input.Headers
	}
	if input.SigningSecret != nil {
		inp["signingSecret"] = *input.SigningSecret
	}
	if input.RetryPolicy != nil {
		inp["retryPolicy"] = input.RetryPolicy
	}
	if input.TimeoutMs != nil {
		inp["timeoutMs"] = *input.TimeoutMs
	}
	if input.RateLimitRps != nil {
		inp["rateLimitRps"] = *input.RateLimitRps
	}
	if input.IdempotencyKey != nil {
		inp["idempotencyKey"] = *input.IdempotencyKey
	}
	vars := map[string]any{"input": inp}
	var result struct {
		SendDynamicMessage *OutboundDelivery `json:"sendDynamicMessage"`
	}
	if err := s.gql.do(ctx, sendDynamicMessageMutation, vars, &result); err != nil {
		return nil, err
	}
	return result.SendDynamicMessage, nil
}
