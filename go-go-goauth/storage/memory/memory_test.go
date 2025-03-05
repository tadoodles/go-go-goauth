//go-go-goauth/storage/memory/memory_test.go

package memory

import (
    "testing"
    "time"

    "go-go-goauth/storage"
    "github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
    store := NewStore()

    email := "test@example.com"
    hashedPassword := "hashedpassword123"

    id, err := store.CreateUser(email, hashedPassword)
    if err != nil {
        t.Fatalf("CreateUser failed: %v", err)
    }

    if id == "" {
        t.Error("Expected non-empty user ID, got empty string")
    }

    user, err := store.GetUserByEmail(email)
    if err != nil {
        t.Fatalf("GetUserByEmail failed: %v", err)
    }

    if user.Email != email {
        t.Errorf("Expected email %s, got %s", email, user.Email)
    }
}

func TestUserExists(t *testing.T) {
    //Case 1: User exists
    
    t.Run("User exists", func(t *testing.T) {
    store := NewStore()

    email := "test@example.com"
    hashedPassword := "hashedpassword123"

    //Create User

    _, err := store.CreateUser(email, hashedPassword)
    if err != nil {
        t.Fatalf("CreateUser failed: %v", err)
    }

    exists, err := store.UserExists(email)
    if err != nil {
        t.Fatalf("UserExists failed: %v", err)
    }
    if !exists {
        t.Error("Expected user to exist, but UserExists returned false")
    }
    })

    //Case 2: User does not exist
    t.Run("User does not exist", func(t *testing.T) {
        store := NewStore()
        unstoredEmail := "unstored@example.com"

        // Check if the user exists
        exists, err := store.UserExists(unstoredEmail)
        if err != nil {
            t.Fatalf("UserExists failed: %v", err)
        }
        if exists {
            t.Error("Expected user to not exist, but UserExists returned true")
        }
    })    
}

func TestGetUserByEmail(t *testing.T) {
    //Case 1: Retrieve a user by their email when the user exists.
    t.Run("Retrieve user by existing email", func(t *testing.T) {
        store := NewStore()
    
        email := "test@example.com"
        hashedPassword := "hashedpassword123"
    
        //Create User
    
        _, err := store.CreateUser(email, hashedPassword)
        if err != nil {
            t.Fatalf("CreateUser failed: %v", err)
        }
    
        //Retrieve user
        user, err := store.GetUserByEmail(email)
        if err != nil {
            t.Fatalf("GetUserByEmail failed: %v", err)
        }

        //Verify the returned user
        if user.Email != email {
            t.Errorf("Expected email %s, got %s", email, user.Email)
        }
        if user.HashedPassword != hashedPassword {
            t.Errorf("Expected hashed password %s, got %s", hashedPassword, user.HashedPassword)
        }
    })

    // Case 2: Attempt to retrieve a user by an email that doesn’t exist.
    t.Run("Retrieve user by non-existent email", func(t *testing.T) {
        store := NewStore()

        unstoredEmail := "unstored@example.com"

        // Retrieve User
        _, err := store.GetUserByEmail(unstoredEmail)
        if err == nil {
            t.Fatal("Expected an error, but got nil")
        }
        if err != storage.ErrUserNotFound {
            t.Errorf("Expected error %v, got %v", storage.ErrUserNotFound, err)
        }
    })
}

func TestGetUserByID(t *testing.T) {
    //Case 1: Retrieve a user by their ID when the user exists
    t.Run("Retrieve user by existing user ID", func(t *testing.T) {
        store := NewStore()
    
        email := "test@example.com"
        hashedPassword := "hashedpassword123"
    
        //Create User
    
        userID, err := store.CreateUser(email, hashedPassword)
        if err != nil {
            t.Fatalf("CreateUser failed: %v", err)
        }
    
        //Check if we can get user by their userID
        user, err := store.GetUserByID(userID)
        if err != nil {
            t.Fatalf("GetUserByID failed: %v", err)
        }

        //Verify the returned user
        if user.Email != email {
            t.Errorf("Expected email %s, got %s", email, user.Email)
        }
        if user.HashedPassword != hashedPassword {
            t.Errorf("Expected hashed password %s, got %s", hashedPassword, user.HashedPassword)
        }
    })

    //Case 2: Attempt to retrieve a user by an ID that doesn’t exist
    t.Run("Retrieve user by non-existent user ID", func(t *testing.T) {
        store := NewStore()

        nonExistentID := "nonexistent-id"

        // Retrieve User by ID
        _, err := store.GetUserByID(nonExistentID)
        if err == nil {
            t.Fatal("Expected an error, but got nil")
        }
        if err != storage.ErrUserNotFound {
            t.Errorf("Expected error %v, got %v", storage.ErrUserNotFound, err)
        }
    })
}

func TestUpdateUser(t *testing.T) {
    //Case 1: Update an existing user successfully
    t.Run("Update existing user", func(t *testing.T) {
        store := NewStore()
    
        email := "test@example.com"
        hashedPassword := "hashedpassword123"
    
        //Create User
    
        userID, err := store.CreateUser(email, hashedPassword)
        if err != nil {
            t.Fatalf("CreateUser failed: %v", err)
        }

        //Retrieve the user by their userID
        user, err := store.GetUserByID(userID)
        if err != nil {
            t.Fatalf("GetUserByID failed: %v", err)
        }
    
        //Modify user's hashed password and update
        modifiedHashedPassword := "newhashedpassword123"
        user.HashedPassword = modifiedHashedPassword
        
        err = store.UpdateUser(user)
        if err != nil {
            t.Fatalf("UpdateUser failed: %v", err)
        }
      
        //Retrieve the updated user by their userID
        updatedUser, err := store.GetUserByID(userID)
        if err != nil {
            t.Fatalf("GetUserByID failed: %v", err)
        }

        //Verify the returned user
        if updatedUser.HashedPassword != modifiedHashedPassword {
            t.Errorf("Expected hashed password %s, got %s", modifiedHashedPassword, updatedUser.HashedPassword)
        }
    })
    //Case 2: Attempt to update a user that doesn’t exist
    t.Run("Update non-existent user", func(t *testing.T) {
        store := NewStore()

        nonExistentUser := storage.User{
            ID:             "nonexistent-id",
            Email:          "nonexistent@example.com",
            HashedPassword: "hashedpassword123",
            CreatedAt:      time.Now(),
            UpdatedAt:      time.Now(),
        }

        // Attempt to update the non-existent user
        err := store.UpdateUser(nonExistentUser) // Pass the user struct directly
        if err == nil {
            t.Fatal("Expected an error, but got nil")
        }
        if err != storage.ErrUserNotFound {
            t.Errorf("Expected error %v, got %v", storage.ErrUserNotFound, err)
        }
    })
}

func TestDeleteUser(t *testing.T) {
    //Case 1: Delete an existing user successfully
    t.Run("Delete existing user", func(t *testing.T) {
        store := NewStore()
    
        email := "test@example.com"
        hashedPassword := "hashedpassword123"
    
        //Create User
    
        userID, err := store.CreateUser(email, hashedPassword)
        if err != nil {
            t.Fatalf("CreateUser failed: %v", err)
        }

        //Delete the user by their ID
        err = store.DeleteUser(userID)
        if err != nil {
            t.Fatalf("DeleteUser failed: %v", err)
        }
        _, err = store.GetUserByID(userID)
        if err == nil {
            t.Fatal("Expected user to be deleted, but GetUserByID returned no error")
        }
        if err != storage.ErrUserNotFound {
            t.Errorf("Expected error %v, got %v", storage.ErrUserNotFound, err)
        }
    })

    //Case 2: Attempt to delete a user that doesn’t exist 
    t.Run("Delete non-existent user", func(t *testing.T) {
        store := NewStore()

        nonExistentID := "nonexistent-id"

        // Attempt to delete the non-existent user
        err := store.DeleteUser(nonExistentID)
        if err == nil {
            t.Fatal("Expected an error, but got nil")
        }
        if err != storage.ErrUserNotFound {
            t.Errorf("Expected error %v, got %v", storage.ErrUserNotFound, err)
        }
    })
}

func TestGetToken(t *testing.T){
    //Case 1: Get an existing token
    t.Run("Get existing token", func(t *testing.T) {
		store := NewStore()

		userID := uuid.New().String()
		tokenID := uuid.New().String()
		expiresAt := time.Now().Add(1 * time.Hour)

		// Store the token
		err := store.StoreToken(userID, tokenID, expiresAt)
		if err != nil {
			t.Fatalf("StoreToken failed: %v", err)
		}

		// Retrieve the token
		token, err := store.GetToken(tokenID)
		if err != nil {
			t.Fatalf("GetToken failed: %v", err)
		}

		// Verify token details
		if token.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, token.UserID)
		}
		if token.ID != tokenID {
			t.Errorf("Expected token ID %s, got %s", tokenID, token.ID)
		}
		if !token.ExpiresAt.Equal(expiresAt) {
			t.Errorf("Expected expiration time %v, got %v", expiresAt, token.ExpiresAt)
		}
	})

    // Case 2: Get a non-existing token
	t.Run("Get non-existing token", func(t *testing.T) {
		store := NewStore()

		// Attempt to retrieve a non-existing token
		_, err := store.GetToken("non-existing-token-id")
		if err != storage.ErrTokenNotFound {
			t.Errorf("Expected ErrTokenNotFound, got %v", err)
		}
	})
}

func TestStoreToken(t *testing.T) {
    //Case 1: Store a new token successfully
    t.Run("Store new token", func(t *testing.T) {
        store := NewStore()

        userID := uuid.New().String()
        tokenID := uuid.New().String()
        expiresAt := time.Now().Add(1 * time.Hour)

        //Store the token
        err := store.StoreToken(userID, tokenID, expiresAt)
        if err != nil {
            t.Fatalf("StoreToken failed: %v", err)
        }

        //Verify the token was stored
        token, err := store.GetToken(tokenID)
        if err != nil {
            t.Fatalf("GetToken failed: %v", err)
        }
        
        //Verify token details
        if token.UserID != userID {
            t.Errorf("Expected user ID %s, got %s", userID, token.UserID)
        }
        if token.ID != tokenID {
            t.Errorf("Expected token ID %s, got %s", tokenID, token.ID)
        }
        if !token.ExpiresAt.Equal(expiresAt) {
            t.Errorf("Expected expiration time %v, got %v", expiresAt, token.ExpiresAt)
        }
    })

    //Case 2: Store a token with an existing ID
    t.Run("Store token with existing ID", func(t *testing.T) {
		store := NewStore()

		userID1 := uuid.New().String()
		userID2 := uuid.New().String()
		tokenID := uuid.New().String()
		expiresAt1 := time.Now().Add(1 * time.Hour)
		expiresAt2 := time.Now().Add(2 * time.Hour)

		// Store the first token
		err := store.StoreToken(userID1, tokenID, expiresAt1)
		if err != nil {
			t.Fatalf("StoreToken failed: %v", err)
		}

		// Attempt to store a second token with the same ID
		err = store.StoreToken(userID2, tokenID, expiresAt2)
		if err != nil {
			t.Fatalf("StoreToken failed: %v", err)
		}

		// Verify the token was overwritten
		token, err := store.GetToken(tokenID)
		if err != nil {
			t.Fatalf("GetToken failed: %v", err)
		}

		// Verify token details
		if token.UserID != userID2 {
			t.Errorf("Expected user ID %s, got %s", userID2, token.UserID)
		}
		if token.ID != tokenID {
			t.Errorf("Expected token ID %s, got %s", tokenID, token.ID)
		}
		if !token.ExpiresAt.Equal(expiresAt2) {
			t.Errorf("Expected expiration time %v, got %v", expiresAt2, token.ExpiresAt)
        }
    })
}

func TestRevokeToken(t *testing.T) {
    //Case 1: Revoke an existing token successfully
    t.Run("Store new token", func(t *testing.T) {
        store := NewStore()

        userID := uuid.New().String()
        tokenID := uuid.New().String()
        expiresAt := time.Now().Add(1 * time.Hour)

        //Store the token
        err := store.StoreToken(userID, tokenID, expiresAt)
        if err != nil {
            t.Fatalf("StoreToken failed: %v", err)
        }
        
        err = store.RevokeToken(tokenID)
        if err != nil {
            t.Fatalf("RevokeToken failed: %v", err)
        }

        // Verify the token is marked as revoked
        revokedAt, exists := store.revokedTokens[tokenID]
        if !exists {
            t.Error("Token was not marked as revoked")
        }
        if revokedAt.IsZero() {
            t.Error("Revocation timestamp is zero")
        }
    })
    //Case 2: Revoke a token that doesn’t exist
    t.Run("Revoke non-existing token", func(t *testing.T) {
        store := NewStore()

        // Attempt to revoke a non-existing token
        err := store.RevokeToken("non-existing-token-id")
        if err != storage.ErrTokenNotFound {
            t.Errorf("Expected ErrTokenNotFound, got %v", err)
        }
    })
}

func TestIsTokenRevoked(t *testing.T) {
    //Case 1: Check if a revoked token is marked as revoked
    t.Run("Check revoked token", func(t *testing.T) {
        store := NewStore()
    
        tokenID := uuid.New().String()
        store.revokedTokens[tokenID] = time.Now()

        isRevoked, err := store.IsTokenRevoked(tokenID)
        if err != nil {
            t.Fatalf("IsTokenRevoked failed: %v", err)
        }

        if !isRevoked {
            t.Error("Expected token to be revoked, but it was not")
        }
    })

    //Case 2: Check if a non-revoked token is not marked as revoked
    t.Run("Check non-revoked token", func(t *testing.T) {
        store := NewStore()

        tokenID := uuid.New().String()

        isRevoked, err := store.IsTokenRevoked(tokenID)
        if err != nil {
            t.Fatalf("IsTokenRevoked failed: %v", err)
        }
        if isRevoked {
            t.Error("Expected token to not be revoked, but it was")
        }
    })
}

func TestCleanExpiredTokens(t *testing.T) {
    //Case 1: Clean up expired revoked tokens
    t.Run("Clean expired tokens", func(t *testing.T) {
        store := NewStore()

        tokenID1 := uuid.New().String()
        tokenID2 := uuid.New().String()

        // Mark tokenID1 as revoked 2 hours ago
        store.revokedTokens[tokenID1] = time.Now().Add(-2 * time.Hour)
        // Mark tokenID2 as revoked now
        store.revokedTokens[tokenID2] = time.Now()

        // Clean up tokens older than 1 hour
        err := store.CleanExpiredTokens(1 * time.Hour)
        if err != nil {
            t.Fatalf("CleanExpiredTokens failed: %v", err)
        }

        // Verify tokenID1 was cleaned up
        _, exists := store.revokedTokens[tokenID1]
        if exists {
            t.Error("Expected tokenID1 to be cleaned up, but it was not")
        }

        // Verify tokenID2 was not cleaned up
        _, exists = store.revokedTokens[tokenID2]
        if !exists {
            t.Error("Expected tokenID2 to not be cleaned up, but it was")
        }
    })

    //Case 2: Ensure non-expired tokens are not cleaned up
    t.Run("Do not clean non-expired tokens", func(t *testing.T) {
        store := NewStore()

        tokenID := uuid.New().String()

        // Mark the token as revoked now
        store.revokedTokens[tokenID] = time.Now()

        // Clean up tokens older than 1 hour
        err := store.CleanExpiredTokens(1 * time.Hour)
        if err != nil {
            t.Fatalf("CleanExpiredTokens failed: %v", err)
        }

        // Verify the token was not cleaned up
        _, exists := store.revokedTokens[tokenID]
        if !exists {
            t.Error("Expected token to not be cleaned up, but it was")
        }
    })
}