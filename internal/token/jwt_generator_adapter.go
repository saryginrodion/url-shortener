package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"roadmap.restapi/internal/errormapper"
)

type JWTGenerator struct {
	secretKey     any
	signingMethod jwt.SigningMethod
	accessExp     time.Duration
	refreshExp    time.Duration
	errMap        *errormapper.ErrorMapper
}

func NewJWTGenerator(
	secret any,
	signingMethod jwt.SigningMethod,
	accessExpirationTime time.Duration,
	refreshExpirationTime time.Duration,
) *JWTGenerator {
	return &JWTGenerator{
		secretKey:     secret,
		signingMethod: signingMethod,
		accessExp:     accessExpirationTime,
		refreshExp:    refreshExpirationTime,
	}
}

func (g *JWTGenerator) Generate(
	ctx context.Context,
	userClaims *UserClaims,
	tokenType TokenType,
	jti uuid.UUID,
) (string, error) {
	expTime := g.accessExp
	if tokenType == REFRESH {
		expTime = g.refreshExp
	}

	claims := &TokenClaims{
		UserClaims: *userClaims,
		Type:       tokenType,
		ID:         jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "url-shortener",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expTime)),
		},
	}

	token := jwt.NewWithClaims(
		g.signingMethod,
		claims,
	)

	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		return "", jwtErrMapper.Map(err)
	}

	return tokenString, nil
}
