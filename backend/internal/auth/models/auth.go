package models

import (
	"time"
)

// AuthCredentials - данные для входа (хранятся в БД)
type AuthCredentials struct {
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// TokenPair - пара токенов для клиента
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
