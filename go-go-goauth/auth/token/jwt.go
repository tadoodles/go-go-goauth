//go-go-goauth/auth/token/jwt.go
package token

import (
    "time"
	"fmt"
    
    "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ErrFailedToSign			= "failed to sign token"
	ErrInvalidToken       	= "invalid token"
    ErrInvalidTokenClaims 	= "invalid token claims"
    ErrInvalidIssuer      	= "invalid issuer"
    ErrUnexpectedSigning  	= "unexpected signing method"
)

// JWTService implements the Tokenizer
type JWTService struct {
	secretKey 	[]byte
	issuer		string
}

func NewJWTService(secretKey string, issuer string) *JWTService {
	return &JWTService{
		secretKey: 	[]byte(secretKey),
		issuer:		issuer,
	}
}

func (s *JWTService) CreateToken(userID, username string, options TokenOptions) (string, error) {
	now := time.Now()

	tokenID := uuid.New().String()

	claims := TokenClaims{
		Subject:    userID,
		Issuer:     options.Issuer,
		Audience:   options.Audience,
		ID:			tokenID,
		IssuedAt:   now,
		Expiration: now.Add(options.ExpiresIn),
		Username:   username,
		Email:      options.Email,
		Roles:      options.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrFailedToSign, err)
	}

	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	//Parse token
	token, err := jwt.ParseWithClaims(
					tokenString,
					claims,
					jwtKeyValidationFunc(s.secretKey))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrInvalidToken, err)
	}

	//Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("%s", ErrInvalidToken)
	}

	//Extract the claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("%s", ErrInvalidTokenClaims)
	}

	//Validate issuer
	if claims.Issuer != s.issuer {
		return nil, fmt.Errorf("%s", ErrInvalidIssuer)
	}

	//Checking user status in database

	return claims, nil
}

//Helper function to validate JWT tokens
func jwtKeyValidationFunc(secretKey []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("%s: %v", ErrUnexpectedSigning, token.Header["alg"])
		}

		return secretKey, nil
	}
}