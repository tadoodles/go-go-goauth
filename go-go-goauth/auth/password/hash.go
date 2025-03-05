//go-go-goauth/auth/password/hash.go

package password

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

const DefaultCost = 12

// Options specifies the cost for bcrypt hashing.
// Cost must be between bcrypt.MinCost (4) and bcrypt.MaxCost (31).
type Options struct {
	Cost int
}

//Hash securely hashes a password using the bcrypt algorithm
func Hash(password string, opts *Options) (string, error) {
	if password == "" {
		return "", errors.New("Password cannot be empty")
	}

	cost := DefaultCost
	if opts != nil && opts.Cost > 0 {
		if opts.Cost < bcrypt.MinCost {
            return "", fmt.Errorf("cost %d is below minimum allowed cost %d", opts.Cost, bcrypt.MinCost)
        }
        if opts.Cost > bcrypt.MaxCost {
            return "", fmt.Errorf("cost %d is above maximum allowed cost %d", opts.Cost, bcrypt.MaxCost)
        }
		cost = opts.Cost
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}


// Verify checks if a plaintext password matches a hashed password
func VerifyHash(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
