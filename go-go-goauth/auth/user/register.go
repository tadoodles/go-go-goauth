//go-go-goauth/auth/user/register.go

package user

import (
	"errors"
	"fmt"
	"strings"

	"go-go-goauth/auth/password"
	"go-go-goauth/storage"
)

type RegisterRequest struct {
	Email		string
	Password	string
}

type RegisterResponse struct {
	Message		string
	UserID		string
}

var (
	ErrEmailRequired   		= errors.New("email required")
	ErrPasswordRequired		= errors.New("password required")
	ErrUserExists     		= errors.New("user already exists")
	ErrInvalidCredentials 	= errors.New("invalid credentials")
	ErrInternalStorage 		= errors.New("internal storage error")
)

func RegisterUser(req RegisterRequest, store storage.UserStore) (RegisterResponse, error) {
	//Step 1: Validate input
	if strings.TrimSpace(req.Email) == "" {
		return RegisterResponse{}, ErrEmailRequired
	}
	if strings.TrimSpace(req.Password) == "" {
		return RegisterResponse{}, ErrPasswordRequired
	}

	//add sanitization steps and more


	//Step 2: Validate the password
	policy := password.NewDefaultPasswordPolicy()
	err := password.ValidatePassword(req.Password, policy)
	if err != nil {
		return RegisterResponse{}, err
	}

	//Step 3: Check if the user already exists
	exists, err := store.UserExists(req.Email)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("%w: %v", ErrInternalStorage, err)
	}
	if exists {
		return RegisterResponse{}, ErrUserExists
	}

	//Step 4: Hash the password
	hashedPassword, err := password.Hash(req.Password, nil)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("%w: %v", ErrInternalStorage, err)
	}

	//Step 5: Create the user
	userID, err := store.CreateUser(req.Email, hashedPassword)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("%w: %v", ErrInternalStorage, err)
	}

	//Step 6: Return success response
	return RegisterResponse{
		Message: 	"User registered successfully",
		UserID:		userID,
	}, nil
}