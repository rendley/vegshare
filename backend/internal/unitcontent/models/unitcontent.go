package models

import (
	"github.com/google/uuid"
	"time"
)

// UnitContent представляет собой запись о содержимом юнита.
// Эта модель универсальна и может представлять что угодно: от растения на грядке до пчел в улье.
type UnitContent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UnitID    uuid.UUID `json:"unit_id" db:"unit_id"`
	UnitType  string    `json:"unit_type" db:"unit_type"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"` // Ссылка на catalog_items
	Quantity  int       `json:"quantity" db:"quantity"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
