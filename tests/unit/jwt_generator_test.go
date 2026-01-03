package unit

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"roadmap.restapi/internal/token"
)

func TestJWTGenerator_Generate_AccessToken(t *testing.T) {
	secret := []byte("test-secret-key-very-long-for-hs256")
	generator := token.NewJWTGenerator(
		secret,
		jwt.SigningMethodHS256,
		15*time.Minute,
		7*24*time.Hour,
	)

	userClaims := &token.UserClaims{
		UID: uuid.New(),
	}
	jti := uuid.New()

	tokenStr, err := generator.Generate(context.Background(), userClaims, token.ACCESS, jti)
	if err != nil {
		t.Fatalf("Generate access token failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("Generated access token is empty")
	}

	// Проверяем, что токен можно распарсить с тем же ключом
	claims := &token.TokenClaims{}
	parsed, parseErr := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if parseErr != nil {
		t.Fatalf("Failed to parse generated token: %v", parseErr)
	}
	if !parsed.Valid {
		t.Fatal("Generated token is not valid according to jwt library")
	}

	if claims.Issuer != "url-shortener" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "url-shortener")
	}
	if claims.UID != userClaims.UID {
		t.Errorf("Uid = %v, want %v", claims.UID, userClaims.UID)
	}
	if claims.Type != token.ACCESS {
		t.Errorf("Type = %q, want %q", claims.Type, token.ACCESS)
	}
	if claims.ID != jti {
		t.Errorf("ID = %v, want %v", claims.ID, jti)
	}

	// Примерная проверка exp (в пределах ±2 сек)
	expectedExp := time.Now().Add(15 * time.Minute)
	if claims.RegisteredClaims.ExpiresAt.Sub(expectedExp).Abs() > time.Second*2 {
		t.Errorf("ExpiresAt = %v, expected around %v", claims.ExpiresAt.Time, expectedExp)
	}
}

func TestJWTGenerator_Generate_RefreshToken(t *testing.T) {
	secret := []byte("test-secret-key-very-long-for-hs256")
	generator := token.NewJWTGenerator(
		secret,
		jwt.SigningMethodHS256,
		15*time.Minute,
		72*time.Hour,
	)

	userClaims := &token.UserClaims{UID: uuid.New()}
	jti := uuid.New()

	tokenStr, err := generator.Generate(context.Background(), userClaims, token.REFRESH, jti)
	if err != nil {
		t.Fatalf("Generate refresh token failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("Generated refresh token is empty")
	}

	claims := &token.TokenClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse refresh token: %v", err)
	}

	if claims.Type != token.REFRESH {
		t.Errorf("Type = %q, want %q", claims.Type, token.REFRESH)
	}

	expectedExp := time.Now().Add(72 * time.Hour)
	if claims.ExpiresAt.Time.Before(expectedExp.Add(-5*time.Second)) ||
		claims.ExpiresAt.Time.After(expectedExp.Add(5*time.Second)) {
		t.Errorf("Refresh ExpiresAt not in expected range: got %v", claims.ExpiresAt.Time)
	}
}
