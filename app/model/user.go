package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserName     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FullName     string     `json:"full_name" db:"full_name"`
	Role         string     `json:"role" db:"role"`
	Preferences  string     `json:"preferences" db:"preferences"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
}

// RegisterUser represents the data required to register a new user
// @Description User registration data
type RegisterUser struct {
	UserName string `json:"username" validate:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
	Role     string `json:"role" validate:"oneof=user admin" example:"user"`
}

// UpdateUser represents the data required to register a new user
// @Description User registration data
type UpdateUser struct {
	UserName string `json:"username" validate:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
	Role     string `json:"role" validate:"oneof=user admin" example:"user"`
}

// LoginUser represents the data required to login
// @Description User login credentials
type LoginUser struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// ChangePassword represents data required to change a user's password
// @Description User password change data
type ChangePassword struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"newpassword123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" example:"newpassword123"`
}

// UserBasic represents basic user information without sensitive data
// @Description Basic user information
type UserBasic struct {
	ID       uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string    `json:"username" example:"johndoe"`
	Email    string    `json:"email" example:"john@example.com"`
	FullName string    `json:"full_name" example:"John Doe"`
	Role     string    `json:"role" example:"user"`
}

// AuthResponse represents the authentication response
// @Description Authentication response with token and user info
type AuthResponse struct {
	Token string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserBasic `json:"user"`
}
