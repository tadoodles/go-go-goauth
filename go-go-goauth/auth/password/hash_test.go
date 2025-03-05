//go-go-goauth/auth/password/hash_test.go
package password

import (
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	// Test successful hashing
	password := "Hello_world2025"
	hash, err := Hash(password, nil)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == "" {
		t.Fatal("Expected hash to not be empty")
	}

	// Test successful verification
	if !VerifyHash(password, hash) {
		t.Fatal("Expected password verification to succeed")
	}

	wrongPassword := "Not_Hello_world2025"
	if VerifyHash(wrongPassword, hash) {
		t.Fatal("Expected wrong password verification to fail")
	}
}

func TestEmptyPassword(t *testing.T) {
	// Test error handling for empty password
	hash, err := Hash("", nil)
	if err == nil {
		t.Fatal("Expected error for empty password")
	}

	if hash != "" {
		t.Fatal("Expected empty hash for failed hashing")
	}
}

func TestCustomCost(t *testing.T) {
	// Test custom cost 
	password := "Hello_world2025"
	customCost := &Options{Cost: 4}

	hash, err := Hash(password, customCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == ""{
		t.Fatal("Expected hash to not be empty")
	}

	if !VerifyHash(password, hash) {
		t.Fatal("Password with custom cost should verify correctly")
	}

	wrongPassword := "Not_Hello_world2025"
	if VerifyHash(wrongPassword, hash) {
		t.Fatal("Expected wrong password verification to fail")
	}
}

func TestInvalidCost(t *testing.T) {
    password := "test_password"
    // Test with cost too low
    lowCost := &Options{Cost: 3} // bcrypt minimum is 4
    _, err := Hash(password, lowCost)
    if err == nil {
        t.Fatal("Expected error with too low cost")
    }
    
    // Test with cost too high
    highCost := &Options{Cost: 32} // bcrypt maximum is 31
    _, err = Hash(password, highCost)
    if err == nil {
        t.Fatal("Expected error with too low cost")
	}
}