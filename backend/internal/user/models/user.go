package models

import (
	"github.com/google/uuid"
	"time"
)

type UserProfile struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	FullName  string    `db:"full_name" json:"full_name"`
	AvatarURL string    `db:"avatar_url" json:"avatar_url"`
	FarmID    uuid.UUID `db:"farm_id" json:"farm_id"` // Связь с модулем farm
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
