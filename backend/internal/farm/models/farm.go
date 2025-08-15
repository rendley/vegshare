// Пакет models содержит структуры данных, которые отражают таблицы из базы данных.
package models

import (
	"github.com/google/uuid"
	"time"
)

// Farm представляет тепличный комплекс (ферму) из таблицы 'farms'.
type Farm struct {
	// Уникальный идентификатор фермы.
	ID uuid.UUID `db:"id" json:"id"`
	// Название фермы.
	Name string `db:"name" json:"name"`
	// Местоположение или адрес фермы.
	Location string `db:"location" json:"location"`
}

// Plot представляет участок земли на ферме из таблицы 'plots'.
type Plot struct {
	// Уникальный идентификатор участка.
	ID uuid.UUID `db:"id" json:"id"`
	// Идентификатор фермы, к которой относится участок.
	FarmID uuid.UUID `db:"farm_id" json:"farm_id"`
	// Идентификатор пользователя-владельца. Может быть пустым (NULL).
	OwnerID uuid.NullUUID `db:"owner_id" json:"owner_id"`
	// Размер участка.
	Size float32 `db:"size" json:"size"`
	// Дата и время создания записи.
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Crop представляет культуру из справочника 'crops'.
type Crop struct {
	// Уникальный идентификатор культуры.
	ID uuid.UUID `db:"id" json:"id"`
	// Название культуры.
	Name string `db:"name" json:"name"`
	// Описание особенностей выращивания.
	Description string `db:"description" json:"description"`
	// Среднее время роста в днях.
	GrowingTimeDays int `db:"growing_time_days" json:"growing_time_days"`
}

// PlotCrop представляет конкретную посадку на участке из таблицы 'plot_crops'.
type PlotCrop struct {
	// Уникальный идентификатор посадки.
	ID uuid.UUID `db:"id" json:"id"`
	// Идентификатор участка.
	PlotID uuid.UUID `db:"plot_id" json:"plot_id"`
	// Идентификатор культуры.
	CropID uuid.UUID `db:"crop_id" json:"crop_id"`
	// Статус посадки.
	Status string `db:"status" json:"status"`
	// Дата посадки.
	PlantedAt time.Time `db:"planted_at" json:"planted_at"`
}