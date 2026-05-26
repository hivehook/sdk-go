package hivehook

// ListOptions is the common pagination/search payload accepted by every List endpoint.
type ListOptions struct {
	Limit  *int
	Offset *int
	After  *string
	First  *int
	Search *string
}

func (o *ListOptions) toVars() map[string]any {
	if o == nil {
		return nil
	}
	vars := map[string]any{}
	if o.Limit != nil {
		vars["limit"] = *o.Limit
	}
	if o.Offset != nil {
		vars["offset"] = *o.Offset
	}
	if o.After != nil {
		vars["after"] = *o.After
	}
	if o.First != nil {
		vars["first"] = *o.First
	}
	if o.Search != nil {
		vars["search"] = *o.Search
	}
	return vars
}

func mergeVars(base, extra map[string]any) map[string]any {
	if base == nil {
		base = map[string]any{}
	}
	for k, v := range extra {
		base[k] = v
	}
	return base
}
