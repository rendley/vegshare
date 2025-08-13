package models

import (
	"github.com/google/uuid"
	"time"
)

type Farm struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OwnerID   uuid.UUID `db:"owner_id" json:"owner_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
