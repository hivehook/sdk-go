package hivehook

import "context"

const alertRuleFragment = `
	id name conditionType threshold webhookUrl channel
	emailConfig { to subjectTemplate }
	slackConfig { webhookUrl channel }
	cooldown enabled createdAt
`

const listAlertRulesQuery = `query($enabled: Boolean, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	alertRules(enabled: $enabled, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + alertRuleFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getAlertRuleQuery = `query($id: UUID!) {
	alertRule(id: $id) {` + alertRuleFragment + `}
}`

const createAlertRuleMutation = `mutation($input: CreateAlertRuleInput!) {
	createAlertRule(input: $input) {` + alertRuleFragment + `}
}`

const updateAlertRuleMutation = `mutation($id: UUID!, $input: UpdateAlertRuleInput!) {
	updateAlertRule(id: $id, input: $input) {` + alertRuleFragment + `}
}`

const deleteAlertRuleMutation = `mutation($id: UUID!) {
	deleteAlertRule(id: $id)
}`

const testAlertRuleMutation = `mutation($id: UUID!) {
	testAlertRule(id: $id)
}`

// AlertRuleService manages alert rules.
type AlertRuleService struct {
	gql *graphqlClient
}

// ListAlertRulesOptions filters AlertRuleService.List results.
type ListAlertRulesOptions struct {
	ListOptions
	Enabled *bool
}

// CreateAlertRuleInput is the input for AlertRuleService.Create.
type CreateAlertRuleInput struct {
	Name          string            `json:"name"`
	ConditionType string            `json:"conditionType"`
	Threshold     int               `json:"threshold"`
	WebhookURL    string            `json:"webhookUrl,omitempty"`
	Channel       *AlertChannel     `json:"channel,omitempty"`
	EmailConfig   *EmailAlertConfig `json:"emailConfig,omitempty"`
	SlackConfig   *SlackAlertConfig `json:"slackConfig,omitempty"`
	Cooldown      *string           `json:"cooldown,omitempty"`
	Enabled       *bool             `json:"enabled,omitempty"`
}

// UpdateAlertRuleInput is the input for AlertRuleService.Update.
type UpdateAlertRuleInput struct {
	Name          *string           `json:"name,omitempty"`
	ConditionType *string           `json:"conditionType,omitempty"`
	Threshold     *int              `json:"threshold,omitempty"`
	WebhookURL    *string           `json:"webhookUrl,omitempty"`
	Channel       *AlertChannel     `json:"channel,omitempty"`
	EmailConfig   *EmailAlertConfig `json:"emailConfig,omitempty"`
	SlackConfig   *SlackAlertConfig `json:"slackConfig,omitempty"`
	Cooldown      *string           `json:"cooldown,omitempty"`
	Enabled       *bool             `json:"enabled,omitempty"`
}

// List returns alert rules matching opts and a PageInfo cursor.
func (s *AlertRuleService) List(ctx context.Context, opts *ListAlertRulesOptions) ([]*AlertRule, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.Enabled != nil {
			vars["enabled"] = *opts.Enabled
		}
	}
	var result struct {
		AlertRules struct {
			Nodes    []*AlertRule `json:"nodes"`
			PageInfo PageInfo     `json:"pageInfo"`
		} `json:"alertRules"`
	}
	if err := s.gql.do(ctx, listAlertRulesQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.AlertRules.Nodes, &result.AlertRules.PageInfo, nil
}

// Get fetches a single AlertRule by ID.
func (s *AlertRuleService) Get(ctx context.Context, id string) (*AlertRule, error) {
	var result struct {
		AlertRule *AlertRule `json:"alertRule"`
	}
	if err := s.gql.do(ctx, getAlertRuleQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.AlertRule, nil
}

// Create persists a new AlertRule.
func (s *AlertRuleService) Create(ctx context.Context, input *CreateAlertRuleInput) (*AlertRule, error) {
	m := map[string]any{
		"name":          input.Name,
		"conditionType": input.ConditionType,
		"threshold":     input.Threshold,
	}
	if input.WebhookURL != "" {
		m["webhookUrl"] = input.WebhookURL
	}
	if input.Channel != nil {
		m["channel"] = string(*input.Channel)
	}
	if input.EmailConfig != nil {
		m["emailConfig"] = map[string]any{
			"to":              input.EmailConfig.To,
			"subjectTemplate": input.EmailConfig.SubjectTemplate,
		}
	}
	if input.SlackConfig != nil {
		sc := map[string]any{"webhookUrl": input.SlackConfig.WebhookURL}
		if input.SlackConfig.Channel != "" {
			sc["channel"] = input.SlackConfig.Channel
		}
		m["slackConfig"] = sc
	}
	if input.Cooldown != nil {
		m["cooldown"] = *input.Cooldown
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	var result struct {
		CreateAlertRule AlertRule `json:"createAlertRule"`
	}
	if err := s.gql.do(ctx, createAlertRuleMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateAlertRule, nil
}

// Update partially updates an existing AlertRule by ID.
func (s *AlertRuleService) Update(ctx context.Context, id string, input *UpdateAlertRuleInput) (*AlertRule, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.ConditionType != nil {
		m["conditionType"] = *input.ConditionType
	}
	if input.Threshold != nil {
		m["threshold"] = *input.Threshold
	}
	if input.WebhookURL != nil {
		m["webhookUrl"] = *input.WebhookURL
	}
	if input.Channel != nil {
		m["channel"] = string(*input.Channel)
	}
	if input.EmailConfig != nil {
		m["emailConfig"] = map[string]any{
			"to":              input.EmailConfig.To,
			"subjectTemplate": input.EmailConfig.SubjectTemplate,
		}
	}
	if input.SlackConfig != nil {
		sc := map[string]any{"webhookUrl": input.SlackConfig.WebhookURL}
		if input.SlackConfig.Channel != "" {
			sc["channel"] = input.SlackConfig.Channel
		}
		m["slackConfig"] = sc
	}
	if input.Cooldown != nil {
		m["cooldown"] = *input.Cooldown
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	var result struct {
		UpdateAlertRule AlertRule `json:"updateAlertRule"`
	}
	if err := s.gql.do(ctx, updateAlertRuleMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateAlertRule, nil
}

// Delete removes an AlertRule by ID.
func (s *AlertRuleService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteAlertRuleMutation, map[string]any{"id": id}, nil)
}

// Test triggers a synthetic firing of the rule (for verifying delivery configuration).
func (s *AlertRuleService) Test(ctx context.Context, id string) error {
	return s.gql.do(ctx, testAlertRuleMutation, map[string]any{"id": id}, nil)
}
