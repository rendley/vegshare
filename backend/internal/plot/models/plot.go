package models

import (
	"time"

	"github.com/google/uuid"
	// Импорт leasing/models больше не нужен
)

// Plot представляет конечную, арендуемую пользователем грядку в теплице.
type Plot struct {
	ID           uuid.UUID `db:"id" json:"id"`
	GreenhouseID uuid.UUID `db:"greenhouse_id" json:"greenhouse_id"`
	Name         string    `db:"name" json:"name"`
	Size         string    `db:"size" json:"size"`
	// Status показывает текущее состояние грядки: 'available', 'rented', 'maintenance'.
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// GetID возвращает ID грядки. удовлетворяя интерфейсу LeasableUnit.
func (p Plot) GetID() uuid.UUID {
	return p.ID
}

// GetStatus возвращает текущий статус грядки.удовлетворяя интерфейсу LeasableUnit.
func (p Plot) GetStatus() string {
	return p.Status
}

// GetUnitType возвращает тип юнита "plot". удовлетворяя интерфейсу LeasableUnit.
// Теперь возвращает string, чтобы соответствовать новому интерфейсу.
func (p Plot) GetUnitType() string {
	return "plot"
}
