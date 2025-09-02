package models

import (
	"time"
	"github.com/google/uuid"
)

// UnitType определяет тип арендуемого юнита.
type UnitType string

const (
	// UnitTypePlot определяет юнит типа "Грядка".
	UnitTypePlot UnitType = "plot"
)

// Lease представляет собой универсальную модель аренды для любого типа юнита.
type Lease struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UnitID    uuid.UUID `db:"unit_id" json:"unit_id"`
	UnitType  UnitType  `db:"unit_type" json:"unit_type"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
