//go-go-goauth/storage/storage.go

package storage

import (
	"time"
	"errors"
)

var (
	ErrUserNotFound   		= errors.New("user not found")
	ErrTokenNotFound		= errors.New("token not found")
	ErrUserExists     		= errors.New("user already exists")
	ErrInvalidCredentials 	= errors.New("invalid credentials")
	ErrInternalStorage 		= errors.New("internal storage error")
)

type User struct {
	ID				string
	Email			string
	HashedPassword	string
	CreatedAt		time.Time
	UpdatedAt		time.Time
}

type UserStore interface {
	CreateUser(email, hashedPassword string) (string, error)

	GetUserByEmail(email string) (User, error)

	GetUserByID(id string) (User, error)

	UserExists(email string) (bool, error)

	UpdateUser(user User) error

	DeleteUser(id string) error
}

type Token struct {
	UserID    	string
    ID        	string
    ExpiresAt 	time.Time
}

type TokenStore interface {
	GetToken(tokenID string) (Token, error)

	StoreToken(userID, tokenID string, expiresAt time.Time) error

	IsTokenRevoked(tokenID string) (bool, error)

	RevokeToken(tokenID string) error

	CleanExpiredTokens(expirationDuration time.Duration) error	
}