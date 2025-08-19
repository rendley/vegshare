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

// PlotCrop - это связующая запись, показывающая, какая культура была посажена на какой грядке в рамках какой аренды.
type PlotCrop struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PlotID    uuid.UUID `db:"plot_id" json:"plot_id"`
	CropID    uuid.UUID `db:"crop_id" json:"crop_id"`
	LeaseID   uuid.UUID `db:"lease_id" json:"lease_id"`
	PlantedAt time.Time `db:"planted_at" json:"planted_at"`
	// Status показывает состояние посадки: 'growing', 'harvested', 'failed'.
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
