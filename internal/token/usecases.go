package token

import (
	"context"

	"github.com/google/uuid"
)

type UseCases struct {
	extractor ClaimsExtractor
	generator Generator
	whitelist WhitelistRepository
}

func NewUseCases(
	extactor ClaimsExtractor,
	generator Generator,
	whitelist WhitelistRepository,
) *UseCases {
	return &UseCases{
		extractor: extactor,
		generator: generator,
		whitelist: whitelist,
	}
}

func (u *UseCases) NewPair(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	accessTokenID := uuid.New()
	refreshTokenID := uuid.New()

	userClaims := &UserClaims{
		UID: userID,
	}

	access, err := u.generator.Generate(ctx, userClaims, ACCESS, accessTokenID)
	if err != nil {
		return nil, err
	}

	refresh, err := u.generator.Generate(ctx, userClaims, REFRESH, refreshTokenID)
	if err != nil {
		return nil, err
	}

	if err = u.whitelist.Add(ctx, refreshTokenID); err != nil {
		return nil, err
	}

	return &TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (u *UseCases) RevokeToken(ctx context.Context, refreshToken string) error {
	claims, err := u.extractor.ParseAndValidate(ctx, refreshToken)
	if err != nil {
		return err
	}

	return u.whitelist.Remove(ctx, claims.ID)
}

func (u *UseCases) RefreshTokenPair(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := u.extractor.ParseAndValidate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	exists, err := u.whitelist.Exists(ctx, claims.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrRefreshNotWhitelisted
	}

	if err = u.whitelist.Remove(ctx, claims.ID); err != nil {
		return nil, err
	}

	return u.NewPair(ctx, claims.UID)
}
