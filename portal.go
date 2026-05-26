package hivehook

import (
	"context"
	"time"
)

// PortalService generates Consumer Portal sessions.
type PortalService struct {
	gql *graphqlClient
}

// PortalToken is a short-lived Consumer Portal session token.
type PortalToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

const generatePortalTokenMutation = `
	mutation GeneratePortalToken($applicationId: UUID!) {
		generatePortalToken(applicationId: $applicationId) {
			token
			expiresAt
		}
	}
`

// GenerateToken issues a short-lived Portal session token for an Application.
func (s *PortalService) GenerateToken(ctx context.Context, applicationID string) (*PortalToken, error) {
	vars := map[string]any{"applicationId": applicationID}
	var result struct {
		GeneratePortalToken PortalToken `json:"generatePortalToken"`
	}
	if err := s.gql.do(ctx, generatePortalTokenMutation, vars, &result); err != nil {
		return nil, err
	}
	return &result.GeneratePortalToken, nil
}
