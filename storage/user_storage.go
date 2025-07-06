package storage

import (
	"errors"
	"loyalty-core/models"
	"sync"
)

// UserStorage provides in-memory storage for users
type UserStorage struct {
	users        map[string]*models.User // userID -> User
	usersByEmail map[string]*models.User // email -> User
	mu           sync.RWMutex
}

// NewUserStorage creates a new user storage instance
func NewUserStorage() *UserStorage {
	return &UserStorage{
		users:        make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
}

// CreateUser adds a new user to storage
func (us *UserStorage) CreateUser(user *models.User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	// Check if user already exists
	if _, exists := us.usersByEmail[user.Email]; exists {
		return errors.New("user already exists")
	}

	us.users[user.ID] = user
	us.usersByEmail[user.Email] = user
	return nil
}

// GetUserByID retrieves a user by ID
func (us *UserStorage) GetUserByID(userID string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (us *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.usersByEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates an existing user
func (us *UserStorage) UpdateUser(user *models.User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	us.users[user.ID] = user
	us.usersByEmail[user.Email] = user
	return nil
}

// GetAllUsers returns all users (for testing)
func (us *UserStorage) GetAllUsers() map[string]*models.User {
	us.mu.RLock()
	defer us.mu.RUnlock()

	users := make(map[string]*models.User)
	for k, v := range us.users {
		users[k] = v
	}
	return users
}

// Global user storage instance
var globalUserStorage *UserStorage

// GetGlobalUserStorage returns the global user storage instance
func GetGlobalUserStorage() *UserStorage {
	if globalUserStorage == nil {
		globalUserStorage = NewUserStorage()
	}
	return globalUserStorage
}
