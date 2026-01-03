package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenType string

const (
	ACCESS  TokenType = "access"
	REFRESH TokenType = "refresh"
)

type UserClaims struct {
	UID uuid.UUID `json:"uid"`
}

type TokenClaims struct {
	UserClaims
	jwt.RegisteredClaims
	Type TokenType `json:"type"`
	ID   uuid.UUID `json:"id"`
}

type TokenPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
