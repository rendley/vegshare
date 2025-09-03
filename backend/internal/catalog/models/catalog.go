package models

import (
	"time"

	"github.com/google/uuid"
)

// Crop представляет собой запись в каталоге доступных для посадки культур.
type Crop struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	// PlantingTime - примерное время посадки в днях.
	PlantingTime int `db:"planting_time" json:"planting_time"`
	// HarvestTime - примерное время сбора урожая в днях после посадки.
	HarvestTime int       `db:"harvest_time" json:"harvest_time"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
