//go-go-goauth/auth/claims.go
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	// Standard Claims
	Subject    	string                 	`json:"sub"`
	Issuer     	string                 	`json:"iss"`
	Audience   	[]string                `json:"aud"`
	ID			string 					`json:"jti"`
	Expiration 	time.Time              	`json:"exp"`
	IssuedAt   	time.Time              	`json:"iat"`

	// Custom Claims
	Username	string					`json:"username,omitempty"`
	Email		string					`json:"email,omitempty"`
	Roles		[]string				`json:"roles,omitempty"`
}

// Implement jwt.Claims interface
func (c TokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings(c.Audience), nil
}

func (c TokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.Expiration), nil
}

func (c TokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.IssuedAt), nil
}

func (c TokenClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c TokenClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c TokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil // NotBefore is not used in this example
}
