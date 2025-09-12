// Пакет models содержит структуры данных, которые отражают таблицы в базе данных.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Region представляет географическую единицу верхнего уровня, например, область или край.
type Region struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// LandParcel представляет большой участок земли в определенном регионе, на котором могут располагаться теплицы.
type LandParcel struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	RegionID  uuid.UUID  `db:"region_id" json:"region_id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Structure представляет строение (теплицу, птичник и т.д.), расположенное на земельном участке.
type Structure struct {
	ID           uuid.UUID `db:"id" json:"id"`
	LandParcelID uuid.UUID `db:"land_parcel_id" json:"land_parcel_id"`
	Name         string    `db:"name" json:"name"`
	// Type определяет тип строения, например, "greenhouse" или "poultry_coop".
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}




