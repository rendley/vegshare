package models

import (
	"github.com/google/uuid"
	"time"
)

// UserProfile представляет данные профиля пользователя, которые можно безопасно отдавать на фронтенд.
// Модель умышленно не содержит связей с другими модулями (например, с арендами),
// чтобы сохранять модульность. Связанные данные, такие как аренды пользователя,
// получаются через отдельные эндпоинты (например, GET /api/v1/leases/my).
type UserProfile struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
