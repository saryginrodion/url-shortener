package unit

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"roadmap.restapi/internal/token"
)

func TestJWTClaimsExtractor_ParseAndValidate_ValidToken(t *testing.T) {
	secret := []byte("test-secret-key-very-long-for-hs256")
	generator := token.NewJWTGenerator(secret, jwt.SigningMethodHS256, time.Hour, 24*time.Hour)
	extractor := token.NewJWTClaimsExtractor(secret, jwt.SigningMethodHS256)

	userClaims := &token.UserClaims{UID: uuid.New()}
	jti := uuid.New()

	tokenStr, err := generator.Generate(context.Background(), userClaims, token.ACCESS, jti)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := extractor.ParseAndValidate(context.Background(), tokenStr)
	if err != nil {
		t.Fatalf("ParseAndValidate failed for valid token: %v", err)
	}
	if claims == nil {
		t.Fatal("ParseAndValidate returned nil claims")
	}

	if claims.UID != userClaims.UID {
		t.Errorf("Uid mismatch: got %v, want %v", claims.UID, userClaims.UID)
	}
	if claims.Type != token.ACCESS {
		t.Errorf("Type = %q, want %q", claims.Type, token.ACCESS)
	}
	if claims.ID != jti {
		t.Errorf("ID mismatch: got %v, want %v", claims.ID, jti)
	}
}

func TestJWTClaimsExtractor_ParseAndValidate_InvalidSignature(t *testing.T) {
	secret := []byte("correct-secret")
	wrongSecret := []byte("wrong-secret")
	generator := token.NewJWTGenerator(secret, jwt.SigningMethodHS256, time.Hour, time.Hour)
	extractor := token.NewJWTClaimsExtractor(wrongSecret, jwt.SigningMethodHS256)

	tokenStr, _ := generator.Generate(context.Background(), &token.UserClaims{UID: uuid.New()}, token.ACCESS, uuid.New())

	_, err := extractor.ParseAndValidate(context.Background(), tokenStr)
	if err == nil {
		t.Fatal("ParseAndValidate accepted token with wrong signature")
	}
	if !strings.Contains(err.Error(), "signature is invalid") {
		t.Errorf("Expected signature error, got: %v", err)
	}
}

func TestJWTClaimsExtractor_ParseAndValidate_ExpiredToken(t *testing.T) {
	secret := []byte("secret")
	generator := token.NewJWTGenerator(secret, jwt.SigningMethodHS256, 10*time.Millisecond, time.Hour)
	extractor := token.NewJWTClaimsExtractor(secret, jwt.SigningMethodHS256)

	tokenStr, _ := generator.Generate(context.Background(), &token.UserClaims{}, token.ACCESS, uuid.New())

	time.Sleep(50 * time.Millisecond) // Ждём истечения

	_, err := extractor.ParseAndValidate(context.Background(), tokenStr)
	if err == nil {
		t.Fatal("ParseAndValidate accepted expired token")
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("Expected 'expired' in error, got: %v", err)
	}
}

func TestJWTClaimsExtractor_ParseAndValidate_MalformedToken(t *testing.T) {
	extractor := token.NewJWTClaimsExtractor("any", jwt.SigningMethodHS256)

	_, err := extractor.ParseAndValidate(context.Background(), "this.is.not.a.jwt")
	if err == nil {
		t.Fatal("ParseAndValidate accepted malformed token")
	}
	if !errors.Is(err, token.ErrTokenMalformed) {
		t.Errorf("expected malformed token error, got: %v", err)
	}
}
