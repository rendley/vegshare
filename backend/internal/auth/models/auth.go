package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// User struct for the auth module, represents the users table.
type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password"` // Struct field uses Hash, but DB column is 'password'
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// TokenPair represents the access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Auth is used for storing refresh token data.
type Auth struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
}
