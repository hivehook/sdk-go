package hivehook

import "context"

const transformationFragment = `
	id name description code enabled failOpen timeoutMs createdAt updatedAt
`

const listTransformationsQuery = `query($enabled: Boolean, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	transformations(enabled: $enabled, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + transformationFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getTransformationQuery = `query($id: UUID!) {
	transformation(id: $id) {` + transformationFragment + `}
}`

const createTransformationMutation = `mutation($input: CreateTransformationInput!) {
	createTransformation(input: $input) {` + transformationFragment + `}
}`

const updateTransformationMutation = `mutation($id: UUID!, $input: UpdateTransformationInput!) {
	updateTransformation(id: $id, input: $input) {` + transformationFragment + `}
}`

const deleteTransformationMutation = `mutation($id: UUID!) {
	deleteTransformation(id: $id)
}`

const testTransformationMutation = `mutation($input: TestTransformationInput!) {
	testTransformation(input: $input) {
		success output error durationMs
	}
}`

// TransformationService manages Transformation scripts.
type TransformationService struct {
	gql *graphqlClient
}

// List returns Transformations matching opts and a PageInfo cursor.
func (s *TransformationService) List(ctx context.Context, opts *ListTransformationsOptions) ([]*Transformation, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.Enabled != nil {
			vars["enabled"] = *opts.Enabled
		}
	}
	var result struct {
		Transformations struct {
			Nodes    []*Transformation `json:"nodes"`
			PageInfo PageInfo          `json:"pageInfo"`
		} `json:"transformations"`
	}
	if err := s.gql.do(ctx, listTransformationsQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Transformations.Nodes, &result.Transformations.PageInfo, nil
}

// Get fetches a single Transformation by ID.
func (s *TransformationService) Get(ctx context.Context, id string) (*Transformation, error) {
	var result struct {
		Transformation *Transformation `json:"transformation"`
	}
	if err := s.gql.do(ctx, getTransformationQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Transformation, nil
}

// Create persists a new Transformation.
func (s *TransformationService) Create(ctx context.Context, input *CreateTransformationInput) (*Transformation, error) {
	m := map[string]any{
		"name": input.Name,
		"code": input.Code,
	}
	if input.Description != "" {
		m["description"] = input.Description
	}
	if input.FailOpen != nil {
		m["failOpen"] = *input.FailOpen
	}
	if input.TimeoutMs != nil {
		m["timeoutMs"] = *input.TimeoutMs
	}
	var result struct {
		CreateTransformation Transformation `json:"createTransformation"`
	}
	if err := s.gql.do(ctx, createTransformationMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.CreateTransformation, nil
}

// Update partially updates an existing Transformation by ID.
func (s *TransformationService) Update(ctx context.Context, id string, input *UpdateTransformationInput) (*Transformation, error) {
	m := map[string]any{}
	if input.Name != nil {
		m["name"] = *input.Name
	}
	if input.Description != nil {
		m["description"] = *input.Description
	}
	if input.Code != nil {
		m["code"] = *input.Code
	}
	if input.Enabled != nil {
		m["enabled"] = *input.Enabled
	}
	if input.FailOpen != nil {
		m["failOpen"] = *input.FailOpen
	}
	if input.TimeoutMs != nil {
		m["timeoutMs"] = *input.TimeoutMs
	}
	var result struct {
		UpdateTransformation Transformation `json:"updateTransformation"`
	}
	if err := s.gql.do(ctx, updateTransformationMutation, map[string]any{"id": id, "input": m}, &result); err != nil {
		return nil, err
	}
	return &result.UpdateTransformation, nil
}

// Delete removes a Transformation by ID.
func (s *TransformationService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteTransformationMutation, map[string]any{"id": id}, nil)
}

// Test runs a synthetic test of the resource (dry-run).
func (s *TransformationService) Test(ctx context.Context, input *TestTransformationInput) (*TransformTestResult, error) {
	m := map[string]any{
		"code":      input.Code,
		"payload":   input.Payload,
		"eventType": input.EventType,
	}
	if input.Headers != nil {
		m["headers"] = input.Headers
	}
	var result struct {
		TestTransformation TransformTestResult `json:"testTransformation"`
	}
	if err := s.gql.do(ctx, testTransformationMutation, map[string]any{"input": m}, &result); err != nil {
		return nil, err
	}
	return &result.TestTransformation, nil
}
