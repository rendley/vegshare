package domain

import "github.com/google/uuid"

// LeasableUnit определяет контракт для любой сущности, которую можно арендовать.
type LeasableUnit interface {
	GetID() uuid.UUID
	GetStatus() string
	GetUnitType() string // Возвращаем string, чтобы избежать импорта пакета models
}
