//go-go-goauth/auth/password/validate.go

package password

import (
	"fmt"
	"errors"
	"strings"
	"unicode"
)

type PasswordPolicy struct {
	MinLength		int
	RequireUpper	bool 
	RequireLower	bool
	RequireNumber	bool
	RequireSpecial	bool
	SpecialChars	string
	DisallowedWords	[]string
}

func NewDefaultPasswordPolicy()  *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:			10,
		RequireUpper:		true,
		RequireLower:		true,
		RequireNumber:		true,
		RequireSpecial: 	true,
		SpecialChars:		"!@#$%^&*()-_+=[]{}|;:,.<>?/\\",
		DisallowedWords:	[]string{"password", "123456","qwerty"},
	}
}

//ValidatePassword checks if a password meets the given policy requirements
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = NewDefaultPasswordPolicy()
	}

	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	if policy.RequireUpper && !containsUppercase(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if policy.RequireLower && !containsLowercase(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if policy.RequireNumber && !containsNumber(password) {
        return errors.New("password must contain at least one number")
    }

	if policy.RequireSpecial && !containsSpecial(password, policy.SpecialChars) {
		return fmt.Errorf("password must contain at least one special character: %s", policy.SpecialChars)
	}

	for _, word := range policy.DisallowedWords {
		if strings.Contains(strings.ToLower(password), strings.ToLower(word)) {
			return fmt.Errorf("password contains a disallowed word: %s", word)
		}
	}

	//Password is valid
	return nil
}

//Helper functions

func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func containsSpecial(s, specialChars string) bool {
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}