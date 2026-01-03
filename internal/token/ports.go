package token

import (
	"context"

	"github.com/google/uuid"
)

type ClaimsExtractor interface {
	ParseAndValidate(ctx context.Context, token string) (*TokenClaims, error)
}

type Generator interface {
	Generate(ctx context.Context, claims *UserClaims, tokenType TokenType, tokenID uuid.UUID) (string, error)
}

type WhitelistRepository interface {
	Add(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Remove(ctx context.Context, id uuid.UUID) error
}
