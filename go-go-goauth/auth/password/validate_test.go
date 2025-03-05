//go-go-goauth/auth/password/validate_test.go

package password

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	policy := NewDefaultPasswordPolicy()

	tests := []struct {
		password string
		valid    bool
	}{
		// Valid passwords
		{"StrongPass!123", true},
		{"AnotherValid1!", true},

		// Invalid passwords
		{"weak", false},                      // Too short
		{"password123!", false},              // Contains disallowed word
		{"1234567890", false},                // No uppercase, lowercase, or special characters
		{"PasswordWithoutSpecial", false},    // No special characters
		{"PASSWORD123!", false},              // No lowercase letters
		{"lowercaseonly!", false},            // No uppercase letters
		{"NoNumbers!", false},                // No numbers
	}

	for _, tt := range tests {
		err := ValidatePassword(tt.password, policy)
		if tt.valid && err != nil {
			t.Errorf("Expected password '%s' to be valid, got error: %v", tt.password, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Expected password '%s' to be invalid, got no error", tt.password)
		}
	}
}