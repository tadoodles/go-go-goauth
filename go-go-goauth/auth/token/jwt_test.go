//go-go-goauth/auth/token/jwt_test.go

package token

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTService_VerifyToken(t *testing.T) {
	secretKey := "test-secret-key"
	issuer := "expected-issuer"
	service := NewJWTService(secretKey, issuer)

	validClaims := TokenClaims{
		Subject:    "user-id",
		Issuer:     issuer,
		Audience:   []string{"test-audience"},
		ID:         uuid.New().String(),
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(1 * time.Hour),
		Username:   "test-user",
		Email:      "test@example.com",
		Roles:      []string{"user"},
	}

	validToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims).SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create valid token: %v", err)
	}

	//Test valid token
	verifiedClaims, err := service.VerifyToken(validToken)
	if err != nil {
		t.Errorf("Expected no error for valid token, got: %v", err)
	}

	if verifiedClaims.Subject != validClaims.Subject {
		t.Errorf("Expected subject %s, got %s", validClaims.Subject, verifiedClaims.Subject)
	}

	if verifiedClaims.Issuer != validClaims.Issuer {
		t.Errorf("Expected issuer %s, got %s", validClaims.Issuer, verifiedClaims.Issuer)
	}

	//Test expired tokens
	expiredClaims := TokenClaims{
		Subject:    "user-id",
		Issuer:     issuer,
		Audience:   []string{"test-audience"},
		ID:         uuid.New().String(),
		IssuedAt:   time.Now().Add(-2 * time.Hour),
		Expiration: time.Now().Add(-1 * time.Hour),
		Username:   "test-user",
		Email:      "test@example.com",
		Roles:      []string{"user"},
	}

	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	_, err = service.VerifyToken(expiredToken)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	} else if !strings.Contains(err.Error(), "token is expired") {
		t.Errorf("Expected 'token is expired' error, got: %v", err)
	}

	//Test invalid issuer
	invalidIssuerClaims := TokenClaims{
		Subject:    "user-id",
		Issuer:     "invalid-issuer-is-me",
		Audience:   []string{"test-audience"},
		ID:         uuid.New().String(),
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(1 * time.Hour),
		Username:   "test-user",
		Email:      "test@example.com",
		Roles:      []string{"user"},
	}

	invalidIssuerToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidIssuerClaims).SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create token with invalid issuer: %v", err)
	}

	_, err = service.VerifyToken(invalidIssuerToken)
	if err == nil {
		t.Error("Expected error for invalid issuer, got nil")
	} else if !strings.Contains(err.Error(), ErrInvalidIssuer) {
		t.Errorf("Expected '%s' error, got: %v", ErrInvalidIssuer, err)
	}

	// Test invalid signing method
	invalidSigningMethodToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, validClaims).SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create token with invalid signing method: %v", err)
	}

	_, err = service.VerifyToken(invalidSigningMethodToken)
	if err == nil {
		t.Error("Expected error for invalid signing method, got nil")
	} else if !strings.Contains(err.Error(), ErrUnexpectedSigning) {
		t.Errorf("Expected '%s' error, got: %v", ErrUnexpectedSigning, err)
	}

	// Test malformed token
	malformedToken := "not.a.valid.jwt.token"
	_, err = service.VerifyToken(malformedToken)
	if err == nil {
		t.Error("Expected error for malformed token, got nil")
	} else if !strings.Contains(err.Error(), "token contains an invalid number of segments") {
		t.Errorf("Expected invalid token error, got: %v", err)
	}
}