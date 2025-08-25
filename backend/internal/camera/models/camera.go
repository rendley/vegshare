package models

import (
	"time"
	"github.com/google/uuid"
)

// Camera представляет камеру, установленную на участке.
type Camera struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PlotID    uuid.UUID `db:"plot_id" json:"plot_id"`
	Name      string    `db:"name" json:"name"`
	RTSPPathName   string    `db:"rtsp_path_name" json:"rtsp_path_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
