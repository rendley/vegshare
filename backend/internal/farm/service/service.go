// Пакет service содержит бизнес-логику, связанную с фермами.
package service

import (
	"context"
	"fmt"

	"github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/rendley/vegshare/backend/internal/farm/repository"
)

// Service - это интерфейс, определяющий контракт для сервиса фермы.
type Service interface {
	GetAllCrops(ctx context.Context) ([]models.Crop, error)
}

// service - это приватная структура, реализующая интерфейс Service.
type service struct {
	repo repository.Repository
}

// NewFarmService - это конструктор для нашего сервиса.
func NewFarmService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// GetAllCrops вызывает репозиторий для получения всех культур.
func (s *service) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	crops, err := s.repo.GetAllCrops(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить культуры в сервисе: %w", err)
	}
	return crops, nil
}