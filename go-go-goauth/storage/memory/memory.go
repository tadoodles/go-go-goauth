//go-go-goauth/storage/memory/memory.go

package memory

import (
	"sync"
	"time"

	"go-go-goauth/storage"
	"github.com/google/uuid"
)

type Store struct {
	//users by ID
	users			map[string]*storage.User
	//user by email for email lookup
	usersByEmail	map[string]string
	//tokens
	tokens			map[string]*storage.Token
	//revoked Tokens
	revokedTokens	map[string]time.Time

	//mutex for concurrent access
	mu		sync.RWMutex
}

func NewStore() *Store {
	return &Store {
		users:			make(map[string]*storage.User),
		usersByEmail:	make(map[string]string),
		tokens:			make(map[string]*storage.Token),
		revokedTokens:	make(map[string]time.Time),		
	}	
}

var _ storage.UserStore = (*Store)(nil)
var _ storage.TokenStore = (*Store)(nil)

func (s *Store) UserExists(email string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.usersByEmail[email]
	return exists, nil
}

func (s *Store) CreateUser(email, hashedPassword string) (string, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[email]; exists {
		return "", storage.ErrUserExists
	}

	user := &storage.User{
		ID:				uuid.New().String(),
		Email:			email,
		HashedPassword:	hashedPassword,
		CreatedAt:		time.Now(),
		UpdatedAt:		time.Now(),
	}

	s.users[user.ID] = user
	s.usersByEmail[user.Email] = user.ID

	return user.ID, nil
}

func (s *Store) GetUserByEmail(email string) (storage.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.usersByEmail[email]
	if !exists {
		return storage.User{}, storage.ErrUserNotFound
	}

	user, exists := s.users[userID]
	if !exists {
		return storage.User{}, storage.ErrUserNotFound
	}

	return *user, nil
}

func (s *Store) GetUserByID(id string) (storage.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return storage.User{}, storage.ErrUserNotFound
	}

	return *user, nil
}

func (s *Store) UpdateUser(user storage.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return storage.ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	s.users[user.ID] = &user

	return nil
}

func (s *Store) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id];
	if !exists {
		return storage.ErrUserNotFound
	}

	delete(s.usersByEmail, user.Email)
	delete(s.users, id)
	return nil
}

func (s *Store) GetToken(tokenID string) (storage.Token, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    token, exists := s.tokens[tokenID]
    if !exists {
        return storage.Token{}, storage.ErrTokenNotFound
    }

    return *token, nil
}

func (s *Store) StoreToken(userID, tokenID string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[tokenID] = &storage.Token{
		UserID:		userID,
		ID:			tokenID,
		ExpiresAt:	expiresAt,
	}

	return nil
}


func (s *Store) RevokeToken(tokenID string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

	// Check if the token exists
    if _, exists := s.tokens[tokenID]; !exists {
        return storage.ErrTokenNotFound
    }
    
    s.revokedTokens[tokenID] = time.Now()
    return nil
}


func (s *Store) IsTokenRevoked(tokenID string) (bool, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    _, isRevoked := s.revokedTokens[tokenID]
    return isRevoked, nil
}


func (s *Store) CleanExpiredTokens(expirationDuration time.Duration) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    cutoff := time.Now().Add(-expirationDuration)
    
    for tokenID, revokedAt := range s.revokedTokens {
        if revokedAt.Before(cutoff) {
            delete(s.revokedTokens, tokenID)
        }
    }
    
    return nil
}