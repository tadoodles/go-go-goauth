//go-go-goauth/auth/token.go
package token

import (
	"time"
)

type Tokenizer interface {
	CreateToken(userID, username string, options TokenOptions) (string, error)
	
	VerifyToken(tokenString string) (*TokenClaims, error)
}

type TokenOptions struct {
	//Required settings
	ExpiresIn	time.Duration

	//Optional settings
	Issuer		string
	Audience	[]string
	Email		string
	Roles		[]string
}

