//go-go-goauth/auth/user/register_test.go

package user

import (
	"testing"

	"go-go-goauth/storage/memory"

)

func TestRegisterUser(t *testing.T) {
	//Case 1: Successful registration
	t.Run("Successful registration", func(t *testing.T) {
		store := memory.NewStore()

		req := RegisterRequest{
			Email:		"test@example.com",
			Password:	"sEcurecutie123!",
		}

		resp, err := RegisterUser(req, store)
		if err != nil {
			t.Fatalf("RegisterUser failed: %v", err)
		}

		if resp.Message != "User registered successfully" {
			t.Errorf("Expected success message, got: %s", resp.Message)
		}
		if resp.UserID == "" {
			t.Error("Expected a non-empty user ID")
		}

		// Verify the user was actually stored
		user, err := store.GetUserByEmail(req.Email)
		if err != nil {
			t.Fatalf("Failed to fetch user: %v", err)
		}
		if user.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, user.Email)
		}
	})

	// Test case 2: Invalid email
	t.Run("Invalid email", func(t *testing.T) {
		store := memory.NewStore()

		req := RegisterRequest{
			Email:    "",
			Password: "sEcurecutie123!",
		}

		_, err := RegisterUser(req, store)
		if err != ErrEmailRequired {
			t.Errorf("Expected ErrEmailRequired, got: %v", err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		store := memory.NewStore()

		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "securePassword123!",
		}

		//expected to fail because Password is a disallowed word

		_, err := RegisterUser(req, store)
		if err == nil {
			t.Error("Expected an error for invalid password, got nil")
		}
	})

	// Test case 4: User already exists
	t.Run("User already exists", func(t *testing.T) {
		store := memory.NewStore()
		
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "sEcurecutie123!",
		}

		// Register the user once
		_, err := RegisterUser(req, store)
		if err != nil {
			t.Fatalf("First registration failed: %v", err)
		}

		// Attempt to register the same user again
		_, err = RegisterUser(req, store)
		if err != ErrUserExists {
			t.Errorf("Expected ErrUserExists, got: %v", err)
		}
	})
}

func TestLoginUser(t *testing.T) {
	//
}