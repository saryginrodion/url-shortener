package token

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaimsExtractor struct {
	pubKey        any
	signingMethod jwt.SigningMethod
}

func NewJWTClaimsExtractor(pubKey any, signingMethod jwt.SigningMethod) *JWTClaimsExtractor {
	return &JWTClaimsExtractor{
		pubKey:        pubKey,
		signingMethod: signingMethod,
	}
}

func (e *JWTClaimsExtractor) ParseAndValidate(ctx context.Context, token string) (*TokenClaims, error) {
	tokenClaims := &TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, tokenClaims, func(t *jwt.Token) (any, error) {
		return e.pubKey, nil
	})

	if err != nil {
		return nil, jwtErrMapper.Map(err)
	}

	if !parsedToken.Valid {
		return nil, ErrTokenInvalid
	}

	return tokenClaims, nil
}
