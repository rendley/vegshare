package handler

import (
	"github.com/google/uuid"
	"time"
)

// LoginRequest - запрос на вход.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse ответ
type LoginResponse struct {
	TokenPair
	UserID string `json:"user_id"`
}

// RegisterRequest - запрос на регистрацию.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterResponse ответ
type RegisterResponse struct {
	TokenPair
	UserID string `json:"user_id"`
}

// TokenPair - ответ с токенами.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// UserResponse - структура для ответа API (без чувствительных данных)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
