package hivehook

import (
	"net/http"
	"time"
)

// Version is the SDK semver string sent in the User-Agent header.
const Version = "0.1.0"

// userAgent is the User-Agent header value used on every outbound request.
const userAgent = "hivehook-sdk-go/" + Version

// Client is the top-level Hivehook API client. Resource services are exposed
// as fields; every method takes a context.Context as its first argument.
type Client struct {
	apiKey     string
	baseURL    string
	http       *http.Client
	maxRetries int

	// Sources manages inbound webhook sources (one per upstream provider).
	Sources *SourceService
	// Destinations manages outbound destinations that receive routed events.
	Destinations *DestinationService
	// Subscriptions manages source-to-destination routing rules.
	Subscriptions *SubscriptionService
	// Events queries ingested events.
	Events *EventService
	// Deliveries queries individual delivery attempts.
	Deliveries *DeliveryService
	// DLQ manages the inbound dead-letter queue.
	DLQ *DLQService
	// APIKeys manages API keys for authenticating SDK callers.
	APIKeys *APIKeyService
	// AlertRules manages alerting rules.
	AlertRules *AlertRuleService
	// Bookmarks manages saved searches/bookmarks.
	Bookmarks *BookmarkService
	// EventTypeSchemas manages event-type JSON schemas.
	EventTypeSchemas *EventTypeSchemaService
	// Applications manages outbound webhook applications.
	Applications *ApplicationService
	// Endpoints manages per-application destination endpoints.
	Endpoints *EndpointService
	// Messages sends outbound webhook messages.
	Messages *MessageService
	// OutboundDeliveries queries outbound delivery attempts.
	OutboundDeliveries *OutboundDeliveryService
	// OutboundDLQ manages the outbound dead-letter queue.
	OutboundDLQ *OutboundDLQService
	// Transformations manages payload transformation rules.
	Transformations *TransformationService
	// Status reports server health and metadata.
	Status *StatusService
	// Portal generates Consumer Portal sessions.
	Portal *PortalService
	// Streams manages event streams.
	Streams *StreamService
	// StreamConsumers manages consumer state for streams.
	StreamConsumers *StreamConsumerService
	// StreamSinks manages stream sinks (Kafka, S3, etc.).
	StreamSinks *StreamSinkService
	// Organizations manages organizations on multi-tenant deployments.
	Organizations *OrganizationService
	// Users manages user accounts.
	Users *UserService
	// AuditLogs queries the audit log.
	AuditLogs *AuditLogService
}

// Option configures a Client passed to New.
type Option func(*Client)

// WithAPIKey sets the API key sent as a Bearer token on every request.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithBaseURL overrides the API base URL (default http://localhost:8080).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient injects a custom *http.Client (overrides WithTimeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithTimeout sets the per-request timeout on the default HTTP client.
// Ignored when WithHTTPClient is also used.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.http == nil {
			c.http = &http.Client{}
		}
		c.http.Timeout = d
	}
}

// WithRetry configures the maximum number of automatic retries for transient
// failures (HTTP 5xx and 429). Default is 2. Pass 0 to disable retries.
// Auth/NotFound/Conflict/Validation errors are never retried.
func WithRetry(n int) Option {
	return func(c *Client) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	}
}

// New constructs a Client with the provided options.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:    "http://localhost:8080",
		maxRetries: 2,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}

	gql := &graphqlClient{
		baseURL:    c.baseURL,
		apiKey:     c.apiKey,
		http:       c.http,
		maxRetries: c.maxRetries,
	}

	c.Sources = &SourceService{gql: gql}
	c.Destinations = &DestinationService{gql: gql}
	c.Subscriptions = &SubscriptionService{gql: gql}
	c.Events = &EventService{gql: gql}
	c.Deliveries = &DeliveryService{gql: gql}
	c.DLQ = &DLQService{gql: gql}
	c.APIKeys = &APIKeyService{gql: gql}
	c.AlertRules = &AlertRuleService{gql: gql}
	c.Bookmarks = &BookmarkService{gql: gql}
	c.EventTypeSchemas = &EventTypeSchemaService{gql: gql}
	c.Applications = &ApplicationService{gql: gql}
	c.Endpoints = &EndpointService{gql: gql}
	c.Messages = &MessageService{gql: gql}
	c.OutboundDeliveries = &OutboundDeliveryService{gql: gql}
	c.OutboundDLQ = &OutboundDLQService{gql: gql}
	c.Transformations = &TransformationService{gql: gql}
	c.Status = &StatusService{gql: gql}
	c.Portal = &PortalService{gql: gql}
	c.Streams = &StreamService{gql: gql}
	c.StreamConsumers = &StreamConsumerService{gql: gql}
	c.StreamSinks = &StreamSinkService{gql: gql}
	c.Organizations = &OrganizationService{gql: gql}
	c.Users = &UserService{gql: gql}
	c.AuditLogs = &AuditLogService{gql: gql}

	return c
}
