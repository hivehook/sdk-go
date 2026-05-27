package hivehook

import "context"

func paginate[T any](fetch func(after string) ([]*T, *PageInfo, error), yield func(*T) bool) error {
	after := ""
	for {
		nodes, pi, err := fetch(after)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			if !yield(n) {
				return nil
			}
		}
		if pi == nil || !pi.HasNextPage || pi.EndCursor == nil {
			return nil
		}
		after = *pi.EndCursor
	}
}

func (s *APIKeyService) Iterate(ctx context.Context, opts *ListOptions, yield func(*APIKey) bool) error {
	o := ListOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*APIKey, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *AlertRuleService) Iterate(ctx context.Context, opts *ListAlertRulesOptions, yield func(*AlertRule) bool) error {
	o := ListAlertRulesOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*AlertRule, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *AuditLogService) Iterate(ctx context.Context, opts *ListAuditLogsOptions, yield func(*AuditLog) bool) error {
	o := ListAuditLogsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*AuditLog, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *BookmarkService) Iterate(ctx context.Context, opts *ListBookmarksOptions, yield func(*Bookmark) bool) error {
	o := ListBookmarksOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Bookmark, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *DestinationService) Iterate(ctx context.Context, opts *ListDestinationsOptions, yield func(*Destination) bool) error {
	o := ListDestinationsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Destination, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *ApplicationService) Iterate(ctx context.Context, opts *ListOptions, yield func(*Application) bool) error {
	o := ListOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Application, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *DeliveryService) Iterate(ctx context.Context, opts *ListDeliveriesOptions, yield func(*Delivery) bool) error {
	o := ListDeliveriesOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Delivery, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *EndpointService) Iterate(ctx context.Context, opts *ListEndpointsOptions, yield func(*Endpoint) bool) error {
	o := ListEndpointsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Endpoint, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *DLQService) Iterate(ctx context.Context, opts *ListDLQOptions, yield func(*DLQEntry) bool) error {
	o := ListDLQOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*DLQEntry, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *EventService) Iterate(ctx context.Context, opts *ListEventsOptions, yield func(*Event) bool) error {
	o := ListEventsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Event, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *EventTypeSchemaService) Iterate(ctx context.Context, opts *ListOptions, yield func(*EventTypeSchema) bool) error {
	o := ListOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*EventTypeSchema, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *OutboundDeliveryService) Iterate(ctx context.Context, opts *ListOutboundDeliveriesOptions, yield func(*OutboundDelivery) bool) error {
	o := ListOutboundDeliveriesOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*OutboundDelivery, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *MessageService) Iterate(ctx context.Context, opts *ListMessagesOptions, yield func(*Message) bool) error {
	o := ListMessagesOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Message, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *OrganizationService) Iterate(ctx context.Context, opts *ListOptions, yield func(*Organization) bool) error {
	o := ListOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Organization, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *SourceService) Iterate(ctx context.Context, opts *ListSourcesOptions, yield func(*Source) bool) error {
	o := ListSourcesOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Source, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *OutboundDLQService) Iterate(ctx context.Context, opts *ListOutboundDLQOptions, yield func(*OutboundDLQEntry) bool) error {
	o := ListOutboundDLQOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*OutboundDLQEntry, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *UserService) Iterate(ctx context.Context, opts *ListUsersOptions, yield func(*User) bool) error {
	o := ListUsersOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*User, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *StreamService) Iterate(ctx context.Context, opts *ListStreamsOptions, yield func(*Stream) bool) error {
	o := ListStreamsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Stream, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *SubscriptionService) Iterate(ctx context.Context, opts *ListSubscriptionsOptions, yield func(*Subscription) bool) error {
	o := ListSubscriptionsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Subscription, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *TransformationService) Iterate(ctx context.Context, opts *ListTransformationsOptions, yield func(*Transformation) bool) error {
	o := ListTransformationsOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*Transformation, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, &o)
	}, yield)
}

func (s *StreamConsumerService) Iterate(ctx context.Context, streamID string, opts *ListStreamConsumersOptions, yield func(*StreamConsumer) bool) error {
	o := ListStreamConsumersOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*StreamConsumer, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, streamID, &o)
	}, yield)
}

func (s *StreamSinkService) Iterate(ctx context.Context, streamID string, opts *ListStreamSinksOptions, yield func(*StreamSink) bool) error {
	o := ListStreamSinksOptions{}
	if opts != nil {
		o = *opts
	}
	return paginate(func(after string) ([]*StreamSink, *PageInfo, error) {
		a := after
		o.After = &a
		return s.List(ctx, streamID, &o)
	}, yield)
}
