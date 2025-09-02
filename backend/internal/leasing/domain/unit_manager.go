package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	// Пакет models больше не нужен, так как LeasableUnit теперь в этом же пакете
)

// UnitManager определяет контракт для сервиса, управляющего юнитами,
// которые можно арендовать (например, plot.Service).
type UnitManager interface {
	// GetLeasableUnit получает юнит по ID.
	GetLeasableUnit(ctx context.Context, unitID uuid.UUID) (LeasableUnit, error) // <- Убрали models.
	// UpdateUnitStatus обновляет статус юнита.
	UpdateUnitStatus(ctx context.Context, unitID uuid.UUID, status string) error
	// WithTx создает экземпляр менеджера, работающий в контексте транзакции.
	WithTx(tx *sqlx.Tx) UnitManager
}
