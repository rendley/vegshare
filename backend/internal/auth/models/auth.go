package models

import (
	"github.com/google/uuid"
	"time"
)

// User - соответствует таблице users в БД.
// db: - для соответствия с названиями колонок в БД
// json: - для сериализации в API (пароль исключён через json:"-")
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// RefreshToken - соответствует таблице refresh_tokens в БД.
type RefreshToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"userId"`
	TokenHash string    `db:"token_hash" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// AuthContext - данные пользователя в контексте HTTP (для middleware).
type AuthContext struct {
	UserID uuid.UUID
	Email  string
}
