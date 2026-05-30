package hivehook

import "time"

// SourceStatus is the lifecycle state of an inbound webhook source.
type SourceStatus string

// SourceStatus values.
const (
	SourceStatusActive   SourceStatus = "ACTIVE"
	SourceStatusInactive SourceStatus = "INACTIVE"
)

// DestinationStatus is the lifecycle state of an outbound destination.
type DestinationStatus string

// DestinationStatus values.
const (
	DestinationStatusActive   DestinationStatus = "ACTIVE"
	DestinationStatusInactive DestinationStatus = "INACTIVE"
)

// EventStatus is the lifecycle state of an ingested event.
type EventStatus string

// EventStatus values.
const (
	EventStatusPending   EventStatus = "PENDING"
	EventStatusDelivered EventStatus = "DELIVERED"
	EventStatusFailed    EventStatus = "FAILED"
)

// DeliveryStatus is the lifecycle state of a delivery attempt.
type DeliveryStatus string

// DeliveryStatus values.
const (
	DeliveryStatusPending    DeliveryStatus = "PENDING"
	DeliveryStatusDelivered  DeliveryStatus = "DELIVERED"
	DeliveryStatusFailed     DeliveryStatus = "FAILED"
	DeliveryStatusSuperseded DeliveryStatus = "SUPERSEDED"
	DeliveryStatusSkipped    DeliveryStatus = "SKIPPED"
)

// AuthType is the authentication scheme a destination uses to call upstream.
type AuthType string

// AuthType values.
const (
	AuthTypeNone   AuthType = "NONE"
	AuthTypeHeader AuthType = "HEADER"
	AuthTypeOAuth2 AuthType = "OAUTH2"
	AuthTypeMTLS   AuthType = "MTLS"
)

// DeliveryMode is push (server initiates) versus poll (client pulls).
type DeliveryMode string

// DeliveryMode values.
const (
	DeliveryModePush DeliveryMode = "PUSH"
	DeliveryModePoll DeliveryMode = "POLL"
)

// DestinationType is the protocol/transport of a destination.
type DestinationType string

// DestinationType values.
const (
	DestinationTypeHTTP        DestinationType = "HTTP"
	DestinationTypeSQS         DestinationType = "SQS"
	DestinationTypeEventBridge DestinationType = "EVENTBRIDGE"
	DestinationTypePubSub      DestinationType = "PUBSUB"
	DestinationTypeKafka       DestinationType = "KAFKA"
	DestinationTypeRabbitMQ    DestinationType = "RABBITMQ"
	DestinationTypeMock        DestinationType = "MOCK"
	DestinationTypeConnector   DestinationType = "CONNECTOR"
)

// OAuth2Config holds client-credentials parameters for OAuth2-authenticated destinations.
type OAuth2Config struct {
	TokenURL     string `json:"tokenUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Scopes       string `json:"scopes"`
	Audience     string `json:"audience"`
}

// HealthConfig sets the auto-disable threshold for an unhealthy destination.
type HealthConfig struct {
	WindowHours  int     `json:"windowHours"`
	DisableBelow float64 `json:"disableBelow"`
}

// EndpointStatus is the lifecycle state of an outbound endpoint.
type EndpointStatus string

// EndpointStatus values.
const (
	EndpointStatusActive   EndpointStatus = "ACTIVE"
	EndpointStatusInactive EndpointStatus = "INACTIVE"
)

// MessageStatus is the lifecycle state of an outbound message.
type MessageStatus string

// MessageStatus values.
const (
	MessageStatusPending   MessageStatus = "PENDING"
	MessageStatusDelivered MessageStatus = "DELIVERED"
	MessageStatusFailed    MessageStatus = "FAILED"
)

// Source is an inbound webhook source bound to a single upstream provider.
type Source struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	ProviderType    string          `json:"providerType"`
	VerifyConfig    map[string]any  `json:"verifyConfig"`
	Status          SourceStatus    `json:"status"`
	RateLimitRps    int             `json:"rateLimitRps"`
	SpikeProtection bool            `json:"spikeProtection"`
	MaxIngestRps    int             `json:"maxIngestRps"`
	BrokerConfig    map[string]any  `json:"brokerConfig,omitempty"`
	ResponseConfig  *ResponseConfig `json:"responseConfig,omitempty"`
	DedupConfig     *DedupConfig    `json:"dedupConfig,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	Subscriptions   []Subscription  `json:"subscriptions,omitempty"`
}

// ResponseConfig controls the static response Hivehook returns to providers.
type ResponseConfig struct {
	StatusCode  int    `json:"statusCode"`
	Body        string `json:"body"`
	ContentType string `json:"contentType"`
}

// DedupConfig configures duplicate-event suppression.
type DedupConfig struct {
	Strategy string   `json:"strategy"`
	Fields   []string `json:"fields,omitempty"`
	Window   string   `json:"window,omitempty"`
}

// Destination is an outbound target that receives routed events.
type Destination struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	URL               string            `json:"url"`
	SigningSecret     string            `json:"signingSecret"`
	Status            DestinationStatus `json:"status"`
	Type              DestinationType   `json:"type"`
	TypeConfig        map[string]any    `json:"typeConfig,omitempty"`
	TimeoutMs         int               `json:"timeoutMs"`
	RateLimitRps      int               `json:"rateLimitRps"`
	RetryPolicy       *RetryPolicy      `json:"retryPolicy"`
	Headers           map[string]any    `json:"headers"`
	AuthType          AuthType          `json:"authType"`
	OAuth2Config      *OAuth2Config     `json:"oauth2Config"`
	MTLSCert          string            `json:"mtlsCert"`
	MTLSKey           string            `json:"mtlsKey"`
	DeliveryMode      DeliveryMode      `json:"deliveryMode"`
	PollAPIKeyPrefix  string            `json:"pollApiKeyPrefix"`
	PollAPIKey        string            `json:"pollApiKey"`
	Ordered           bool              `json:"ordered"`
	BlockedDeliveryID *string           `json:"blockedDeliveryId"`
	HealthScore       float64           `json:"healthScore"`
	DisabledReason    *string           `json:"disabledReason"`
	HealthConfig      *HealthConfig     `json:"healthConfig"`
	OutputFormat      string            `json:"outputFormat"`
	CreatedAt         time.Time         `json:"createdAt"`
	Subscriptions     []Subscription    `json:"subscriptions,omitempty"`
}

// RetryPolicy controls retry attempts and exponential backoff for deliveries.
type RetryPolicy struct {
	MaxAttempts   int     `json:"maxAttempts"`
	InitialDelay  string  `json:"initialDelay"`
	MaxDelay      string  `json:"maxDelay"`
	BackoffFactor float64 `json:"backoffFactor"`
}

// Subscription routes events from one Source to one Destination.
type Subscription struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	SourceID        string           `json:"sourceId"`
	DestinationID   string           `json:"destinationId"`
	FilterConfig    *FilterConfig    `json:"filterConfig"`
	TransformConfig *TransformConfig `json:"transformConfig"`
	Enabled         bool             `json:"enabled"`
	CreatedAt       time.Time        `json:"createdAt"`
	Source          *Source          `json:"source,omitempty"`
	Destination     *Destination     `json:"destination,omitempty"`
}

// FilterConfig restricts which events flow through a Subscription.
type FilterConfig struct {
	EventTypes []string        `json:"eventTypes,omitempty"`
	Regex      []string        `json:"regex,omitempty"`
	BodyMatch  []BodyMatchRule `json:"bodyMatch,omitempty"`
	Rules      []FilterRule    `json:"rules,omitempty"`
}

// BodyMatchRule matches a JSON path against an operator and value.
type BodyMatchRule struct {
	Path     string `json:"path"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// FilterRule is a composable filter clause (single match or nested group).
type FilterRule struct {
	Path     string       `json:"path,omitempty"`
	Operator string       `json:"operator"`
	Value    any          `json:"value,omitempty"`
	Rules    []FilterRule `json:"rules,omitempty"`
}

// TransformConfig controls how the payload is envelope-wrapped and header-augmented.
type TransformConfig struct {
	Envelope bool           `json:"envelope"`
	Headers  map[string]any `json:"headers,omitempty"`
}

// Event is an ingested webhook event awaiting (or completing) delivery.
type Event struct {
	ID             string         `json:"id"`
	SourceID       string         `json:"sourceId"`
	IdempotencyKey string         `json:"idempotencyKey"`
	EventType      string         `json:"eventType"`
	Headers        map[string]any `json:"headers"`
	RawBody        []byte         `json:"rawBody"`
	Status         EventStatus    `json:"status"`
	ReceivedAt     time.Time      `json:"receivedAt"`
	Source         *Source        `json:"source,omitempty"`
	Deliveries     []Delivery     `json:"deliveries,omitempty"`
}

// Delivery is a single Event-to-Destination delivery record.
type Delivery struct {
	ID               string            `json:"id"`
	EventID          string            `json:"eventId"`
	SubscriptionID   string            `json:"subscriptionId"`
	DestinationID    string            `json:"destinationId"`
	Status           DeliveryStatus    `json:"status"`
	Attempts         int               `json:"attempts"`
	MaxAttempts      int               `json:"maxAttempts"`
	NextAttemptAt    *time.Time        `json:"nextAttemptAt"`
	CreatedAt        time.Time         `json:"createdAt"`
	Event            *Event            `json:"event,omitempty"`
	Subscription     *Subscription     `json:"subscription,omitempty"`
	Destination      *Destination      `json:"destination,omitempty"`
	DeliveryAttempts []DeliveryAttempt `json:"deliveryAttempts,omitempty"`
}

// DeliveryAttempt is one HTTP attempt (with status, latency, error) for a Delivery.
type DeliveryAttempt struct {
	ID             string    `json:"id"`
	DeliveryID     string    `json:"deliveryId"`
	AttemptNumber  int       `json:"attemptNumber"`
	ResponseStatus int       `json:"responseStatus"`
	ResponseBody   string    `json:"responseBody"`
	Error          string    `json:"error"`
	DurationMs     int       `json:"durationMs"`
	AttemptedAt    time.Time `json:"attemptedAt"`
}

// DLQEntry is a permanently-failed Delivery parked in the dead-letter queue.
type DLQEntry struct {
	ID         string     `json:"id"`
	DeliveryID string     `json:"deliveryId"`
	EventID    string     `json:"eventId"`
	LastError  string     `json:"lastError"`
	ReplayedAt *time.Time `json:"replayedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	Delivery   *Delivery  `json:"delivery,omitempty"`
	Event      *Event     `json:"event,omitempty"`
}

// APIKey is an API key metadata record (the secret bytes are never returned).
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	Scopes     []string   `json:"scopes"`
	SourceIDs  []string   `json:"sourceIds"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	RevokedAt  *time.Time `json:"revokedAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
}

// APIKeyWithSecret is an APIKey paired with the one-time-visible raw secret.
type APIKeyWithSecret struct {
	APIKey APIKey `json:"apiKey"`
	RawKey string `json:"rawKey"`
}

// AlertChannel selects how an AlertRule notifies (webhook, email, Slack).
type AlertChannel string

// AlertChannel values.
const (
	AlertChannelWebhook AlertChannel = "WEBHOOK"
	AlertChannelEmail   AlertChannel = "EMAIL"
	AlertChannelSlack   AlertChannel = "SLACK"
)

// EmailAlertConfig configures the email notification target.
type EmailAlertConfig struct {
	To              []string `json:"to"`
	SubjectTemplate string   `json:"subjectTemplate,omitempty"`
}

// SlackAlertConfig configures the Slack webhook notification target.
type SlackAlertConfig struct {
	WebhookURL string `json:"webhookUrl"`
	Channel    string `json:"channel,omitempty"`
}

// AlertRule is a server-side alerting rule (threshold, cooldown, channel).
type AlertRule struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	ConditionType string            `json:"conditionType"`
	Threshold     int               `json:"threshold"`
	WebhookURL    string            `json:"webhookUrl"`
	Channel       AlertChannel      `json:"channel"`
	EmailConfig   *EmailAlertConfig `json:"emailConfig,omitempty"`
	SlackConfig   *SlackAlertConfig `json:"slackConfig,omitempty"`
	Cooldown      string            `json:"cooldown"`
	Enabled       bool              `json:"enabled"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// Bookmark is a saved reference to an Event.
type Bookmark struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventId"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	Event     *Event    `json:"event,omitempty"`
}

// EventTypeSchema is a JSON schema describing the payload shape of an event type.
type EventTypeSchema struct {
	ID          string         `json:"id"`
	EventType   string         `json:"eventType"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"schema"`
	Example     map[string]any `json:"example"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// Application is a tenant of the outbound webhook gateway.
type Application struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	UID       string     `json:"uid"`
	CreatedAt time.Time  `json:"createdAt"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// Endpoint is an outbound destination owned by an Application.
type Endpoint struct {
	ID                string          `json:"id"`
	ApplicationID     string          `json:"applicationId"`
	URL               string          `json:"url"`
	SigningSecret     string          `json:"signingSecret"`
	FilterConfig      *FilterConfig   `json:"filterConfig"`
	Status            EndpointStatus  `json:"status"`
	Type              DestinationType `json:"type"`
	TypeConfig        map[string]any  `json:"typeConfig,omitempty"`
	RateLimitRps      int             `json:"rateLimitRps"`
	TimeoutMs         int             `json:"timeoutMs"`
	RetryPolicy       *RetryPolicy    `json:"retryPolicy"`
	Headers           map[string]any  `json:"headers"`
	AuthType          AuthType        `json:"authType"`
	OAuth2Config      *OAuth2Config   `json:"oauth2Config"`
	MTLSCert          string          `json:"mtlsCert"`
	MTLSKey           string          `json:"mtlsKey"`
	DeliveryMode      DeliveryMode    `json:"deliveryMode"`
	PollAPIKeyPrefix  string          `json:"pollApiKeyPrefix"`
	PollAPIKey        string          `json:"pollApiKey"`
	Ordered           bool            `json:"ordered"`
	BlockedDeliveryID *string         `json:"blockedDeliveryId"`
	HealthScore       float64         `json:"healthScore"`
	DisabledReason    *string         `json:"disabledReason"`
	HealthConfig      *HealthConfig   `json:"healthConfig"`
	OutputFormat      string          `json:"outputFormat"`
	CreatedAt         time.Time       `json:"createdAt"`
	Application       *Application    `json:"application,omitempty"`
}

// Message is an outbound webhook message addressed to an Application's Endpoints.
type Message struct {
	ID                 string             `json:"id"`
	ApplicationID      string             `json:"applicationId"`
	EventType          string             `json:"eventType"`
	Payload            []byte             `json:"payload"`
	IdempotencyKey     string             `json:"idempotencyKey"`
	Status             MessageStatus      `json:"status"`
	CreatedAt          time.Time          `json:"createdAt"`
	Application        *Application       `json:"application,omitempty"`
	OutboundDeliveries []OutboundDelivery `json:"outboundDeliveries,omitempty"`
}

// OutboundDelivery is a Message-to-Endpoint delivery record.
type OutboundDelivery struct {
	ID               string                    `json:"id"`
	MessageID        string                    `json:"messageId"`
	EndpointID       string                    `json:"endpointId"`
	Status           DeliveryStatus            `json:"status"`
	Attempts         int                       `json:"attempts"`
	MaxAttempts      int                       `json:"maxAttempts"`
	NextAttemptAt    *time.Time                `json:"nextAttemptAt"`
	CreatedAt        time.Time                 `json:"createdAt"`
	Message          *Message                  `json:"message,omitempty"`
	Endpoint         *Endpoint                 `json:"endpoint,omitempty"`
	DeliveryAttempts []OutboundDeliveryAttempt `json:"deliveryAttempts,omitempty"`
}

// OutboundDeliveryAttempt is one HTTP attempt for an OutboundDelivery.
type OutboundDeliveryAttempt struct {
	ID             string    `json:"id"`
	DeliveryID     string    `json:"deliveryId"`
	AttemptNumber  int       `json:"attemptNumber"`
	ResponseStatus int       `json:"responseStatus"`
	ResponseBody   string    `json:"responseBody"`
	Error          string    `json:"error"`
	DurationMs     int       `json:"durationMs"`
	AttemptedAt    time.Time `json:"attemptedAt"`
}

// OutboundDLQEntry is a permanently-failed OutboundDelivery parked in the DLQ.
type OutboundDLQEntry struct {
	ID         string            `json:"id"`
	DeliveryID string            `json:"deliveryId"`
	MessageID  string            `json:"messageId"`
	LastError  string            `json:"lastError"`
	ReplayedAt *time.Time        `json:"replayedAt"`
	CreatedAt  time.Time         `json:"createdAt"`
	Delivery   *OutboundDelivery `json:"delivery,omitempty"`
	Message    *Message          `json:"message,omitempty"`
}

// Transformation is a user-supplied payload transformation script.
type Transformation struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Enabled     bool      `json:"enabled"`
	FailOpen    bool      `json:"failOpen"`
	TimeoutMs   int       `json:"timeoutMs"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ListTransformationsOptions filters Transformations.List results.
type ListTransformationsOptions struct {
	ListOptions
	Enabled *bool
}

// CreateTransformationInput is the input for TransformationService.Create.
type CreateTransformationInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Code        string `json:"code"`
	FailOpen    *bool  `json:"failOpen,omitempty"`
	TimeoutMs   *int   `json:"timeoutMs,omitempty"`
}

// UpdateTransformationInput is the input for TransformationService.Update.
type UpdateTransformationInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Code        *string `json:"code,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	FailOpen    *bool   `json:"failOpen,omitempty"`
	TimeoutMs   *int    `json:"timeoutMs,omitempty"`
}

// TestTransformationInput is the input for TransformationService.Test (dry-run).
type TestTransformationInput struct {
	Code      string         `json:"code"`
	Payload   map[string]any `json:"payload"`
	EventType string         `json:"eventType"`
	Headers   map[string]any `json:"headers,omitempty"`
}

// TransformTestResult is the result of a Transformation dry-run.
type TransformTestResult struct {
	Success    bool           `json:"success"`
	Output     map[string]any `json:"output,omitempty"`
	Error      string         `json:"error,omitempty"`
	DurationMs int            `json:"durationMs"`
}

// SystemStatus is a snapshot of server health and aggregate counters.
type SystemStatus struct {
	Status                    string `json:"status"`
	DLQSize                   int    `json:"dlqSize"`
	OutboundDLQSize           int    `json:"outboundDlqSize"`
	QueueDepth                int    `json:"queueDepth"`
	ActiveWorkers             int    `json:"activeWorkers"`
	TotalWorkers              int    `json:"totalWorkers"`
	Uptime                    int    `json:"uptime"`
	Version                   string `json:"version"`
	SourcesTotal              int    `json:"sourcesTotal"`
	DestinationsTotal         int    `json:"destinationsTotal"`
	SubscriptionsTotal        int    `json:"subscriptionsTotal"`
	EventsTotal               int    `json:"eventsTotal"`
	EventsFailed              int    `json:"eventsFailed"`
	DeliveriesTotal           int    `json:"deliveriesTotal"`
	DeliveriesPending         int    `json:"deliveriesPending"`
	DeliveriesDelivered       int    `json:"deliveriesDelivered"`
	MessagesTotal             int    `json:"messagesTotal"`
	OutboundDeliveriesTotal   int    `json:"outboundDeliveriesTotal"`
	OutboundDeliveriesPending int    `json:"outboundDeliveriesPending"`
	OutboundDeliveriesFailed  int    `json:"outboundDeliveriesFailed"`
}

// ReplayResult reports how many deliveries were replayed by a DLQ operation.
type ReplayResult struct {
	Deliveries int `json:"deliveries"`
}

// PurgeResult reports how many entries were purged by a DLQ operation.
type PurgeResult struct {
	Purged int `json:"purged"`
}

// StreamStatus is the lifecycle state of a Stream.
type StreamStatus string

// StreamStatus values.
const (
	StreamStatusActive StreamStatus = "ACTIVE"
	StreamStatusPaused StreamStatus = "PAUSED"
)

// SinkType is the protocol/destination of a StreamSink.
type SinkType string

// SinkType values.
const (
	SinkTypeS3       SinkType = "S3"
	SinkTypeWebhook  SinkType = "WEBHOOK"
	SinkTypePostgres SinkType = "POSTGRES"
)

// SinkStatus is the lifecycle state of a StreamSink.
type SinkStatus string

// SinkStatus values.
const (
	SinkStatusActive SinkStatus = "ACTIVE"
	SinkStatusPaused SinkStatus = "PAUSED"
)

// Stream is a durable ordered log of events for an Application.
type Stream struct {
	ID            string       `json:"id"`
	ApplicationID string       `json:"applicationId"`
	Name          string       `json:"name"`
	Status        StreamStatus `json:"status"`
	RetentionDays int          `json:"retentionDays"`
	CreatedAt     time.Time    `json:"createdAt"`
}

// StreamConsumer tracks a cursor position on a Stream.
type StreamConsumer struct {
	ID             string    `json:"id"`
	StreamID       string    `json:"streamId"`
	Name           string    `json:"name"`
	CursorSequence int       `json:"cursorSequence"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// StreamSink ships Stream events to an external sink (S3, webhook, Postgres).
type StreamSink struct {
	ID             string         `json:"id"`
	StreamID       string         `json:"streamId"`
	Name           string         `json:"name"`
	SinkType       SinkType       `json:"sinkType"`
	Config         map[string]any `json:"config"`
	BatchSize      int            `json:"batchSize"`
	FlushInterval  string         `json:"flushInterval"`
	CursorSequence int            `json:"cursorSequence"`
	Status         SinkStatus     `json:"status"`
	LastFlushedAt  *time.Time     `json:"lastFlushedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
}

// StreamEntry is a single message persisted in a Stream's log.
type StreamEntry struct {
	ID        string    `json:"id"`
	StreamID  string    `json:"streamId"`
	Sequence  int       `json:"sequence"`
	MessageID *string   `json:"messageId"`
	EventType string    `json:"eventType"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

// MetaEventConfig describes a webhook that receives meta-events (events about
// events, e.g. delivery.failed, source.created) emitted by Hivehook itself.
type MetaEventConfig struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	SigningSecret string    `json:"signingSecret"`
	EventTypes    []string  `json:"eventTypes"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"createdAt"`
}

// OTLPConfig holds OpenTelemetry exporter settings for an Organization.
type OTLPConfig struct {
	Endpoint   string         `json:"endpoint"`
	Headers    map[string]any `json:"headers,omitempty"`
	Insecure   bool           `json:"insecure"`
	SampleRate float64        `json:"sampleRate"`
}

// Organization is a tenant boundary on a multi-tenant deployment.
type Organization struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Slug              string      `json:"slug"`
	SSOEnabled        bool        `json:"ssoEnabled"`
	SSOProvider       string      `json:"ssoProvider"`
	RetentionEvents   int         `json:"retentionEvents"`
	RetentionMessages int         `json:"retentionMessages"`
	OTLPConfig        *OTLPConfig `json:"otlpConfig,omitempty"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

// User is a human account within an Organization.
type User struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organizationId"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	Role           string     `json:"role"`
	LastLoginAt    *time.Time `json:"lastLoginAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// AuditLog is a single audit-log entry capturing actor, action, and resource.
type AuditLog struct {
	ID           string         `json:"id"`
	ActorType    string         `json:"actorType"`
	ActorID      string         `json:"actorId"`
	ActorName    string         `json:"actorName"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resourceType"`
	ResourceID   string         `json:"resourceId"`
	OrgID        string         `json:"orgId"`
	IPAddress    string         `json:"ipAddress"`
	UserAgent    string         `json:"userAgent"`
	Details      map[string]any `json:"details"`
	CreatedAt    time.Time      `json:"createdAt"`
}

// PageInfo describes the pagination cursor returned by list endpoints.
type PageInfo struct {
	Total       int     `json:"total"`
	Limit       int     `json:"limit"`
	Offset      int     `json:"offset"`
	EndCursor   *string `json:"endCursor"`
	HasNextPage bool    `json:"hasNextPage"`
}
