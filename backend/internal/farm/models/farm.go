// Пакет models содержит структуры данных, которые отражают таблицы в базе данных.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Region представляет географическую единицу верхнего уровня, например, область или край.
type Region struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// LandParcel представляет большой участок земли в определенном регионе, на котором могут располагаться теплицы.
type LandParcel struct {
	ID        uuid.UUID `db:"id" json:"id"`
	RegionID  uuid.UUID `db:"region_id" json:"region_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Greenhouse представляет теплицу, расположенную на земельном участке.
type Greenhouse struct {
	ID           uuid.UUID `db:"id" json:"id"`
	LandParcelID uuid.UUID `db:"land_parcel_id" json:"land_parcel_id"`
	Name         string    `db:"name" json:"name"`
	// Type может определять тип теплицы, например, "гидропоника" или "грунтовая".
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Plot представляет конечную, арендуемую пользователем грядку в теплице.
type Plot struct {
	ID           uuid.UUID `db:"id" json:"id"`
	GreenhouseID uuid.UUID `db:"greenhouse_id" json:"greenhouse_id"`
	Name         string    `db:"name" json:"name"`
	Size         string    `db:"size" json:"size"`
	// Status показывает текущее состояние грядки: 'available', 'rented', 'maintenance'.
	Status    string    `db:"status" json:"status"`
	CameraURL string    `db:"camera_url" json:"camera_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// PlotLease представляет собой факт аренды грядки пользователем на определенный срок.
// Эта таблица связывает пользователей и грядки.
type PlotLease struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PlotID    uuid.UUID `db:"plot_id" json:"plot_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	// Status показывает состояние аренды: 'active', 'expired'.
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}


