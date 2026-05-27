// Package hivehook is the official Go client for Hivehook, webhook
// infrastructure for modern teams (inbound and outbound).
//
// Construct a client with [New] and an [Option] for the base URL and API key:
//
//	client := hivehook.New(
//	    hivehook.WithBaseURL("http://localhost:8080"),
//	    hivehook.WithAPIKey(os.Getenv("HIVEHOOK_API_KEY")),
//	)
//
// Resources are exposed as fields on the client (Sources, Destinations,
// Subscriptions, Events, Deliveries, DLQ, APIKeys, AlertRules, Bookmarks,
// EventTypeSchemas, Applications, Endpoints, Messages, OutboundDeliveries,
// OutboundDLQ, Transformations, Status, Portal, Streams, StreamConsumers,
// StreamSinks, Organizations, Users, AuditLogs). Every method takes a
// [context.Context] as its first argument.
//
// The hivehook/webhook subpackage verifies inbound HMAC signatures from
// Hivehook's outbound deliveries.
//
// See https://hivehook.com/docs for the full reference.
package hivehook
