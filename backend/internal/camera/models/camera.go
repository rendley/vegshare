package models

import (
	"time"

	"github.com/google/uuid"
)

// Camera представляет камеру, установленную на универсальном юните (грядке, клетке и т.д.).
type Camera struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UnitID       uuid.UUID `db:"unit_id" json:"unit_id"`
	UnitType     string    `db:"unit_type" json:"unit_type"`
	Name         string    `db:"name" json:"name"`
	RTSPPathName string    `db:"rtsp_path_name" json:"rtsp_path_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
