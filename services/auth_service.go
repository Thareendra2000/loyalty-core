package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/storage"
	"loyalty-core/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	config      *config.Config
	userStorage *storage.UserStorage
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config:      cfg,
		userStorage: storage.GetGlobalUserStorage(),
	}
}

// generateLoyaltyID generates a unique loyalty ID
func (as *AuthService) generateLoyaltyID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return "LOY" + string(b)
}

// generateUserID generates a unique user ID
func (as *AuthService) generateUserID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// validateEmail checks if email is valid format
func (as *AuthService) validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// validatePassword checks if password meets requirements
func (as *AuthService) validatePassword(password string) bool {
	return len(password) >= 6
}

// hashPassword hashes the password using bcrypt
func (as *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares password with hash
func (as *AuthService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SignupUser creates a new user account
func (as *AuthService) SignupUser(req models.SignupRequest) (*models.SignupResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("all fields are required")
	}

	// Validate email format
	if !as.validateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password strength
	if !as.validatePassword(req.Password) {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Check if user already exists
	if _, err := as.userStorage.GetUserByEmail(req.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := as.hashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("internal server error")
	}

	// Create new user
	user := &models.User{
		ID:        as.generateUserID(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		LoyaltyID: as.generateLoyaltyID(),
		Points:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user (using shared storage)
	if err := as.userStorage.CreateUser(user); err != nil {
		return nil, err
	}

	// Create response (exclude password)
	responseUser := *user
	responseUser.Password = ""

	response := &models.SignupResponse{
		Message: "User created successfully",
		User:    responseUser,
	}

	log.Printf("User created: %s", user.Email)
	return response, nil
}

// LoginUser authenticates a user and returns login response
func (as *AuthService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user by email
	foundUser, err := as.userStorage.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !as.checkPasswordHash(req.Password, foundUser.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(foundUser.ID, foundUser.Email, as.config.JWTSecret)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("internal server error")
	}

	// Create response (exclude password)
	responseUser := *foundUser
	responseUser.Password = ""

	response := &models.LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    responseUser,
	}

	log.Printf("User logged in: %s", foundUser.Email)
	return response, nil
}

// GetAllUsers returns all users (for testing purposes)
func (as *AuthService) GetAllUsers() map[string]*models.User {
	return as.userStorage.GetAllUsers()
}

// ValidateToken validates JWT token and returns user claims
func (as *AuthService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return utils.ValidateToken(tokenString, as.config.JWTSecret)
}

// GetUserProfile retrieves user profile by user ID
func (as *AuthService) GetUserProfile(userID string) (*models.User, error) {
	user, err := as.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Create response (exclude password)
	responseUser := *user
	responseUser.Password = ""

	return &responseUser, nil
}
