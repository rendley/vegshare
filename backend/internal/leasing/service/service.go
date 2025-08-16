// Пакет service содержит бизнес-логику, связанную с арендой.
package service

import (
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	"github.com/rendley/vegshare/backend/internal/leasing/repository"
)

// Service определяет контракт для бизнес-логики аренды.
type Service interface {
	// Здесь будет метод LeasePlot(ctx, userID, plotID)
}

type service struct {
	repo     repository.Repository
	farmRepo farmRepository.Repository // Зависимость от репозитория другого модуля
}

// NewLeasingService - конструктор для сервиса аренды.
func NewLeasingService(repo repository.Repository, farmRepo farmRepository.Repository) Service {
	return &service{
		repo:     repo,
		farmRepo: farmRepo,
	}
}
