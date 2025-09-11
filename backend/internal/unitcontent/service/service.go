package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/unitcontent/models"
	"github.com/rendley/vegshare/backend/internal/unitcontent/repository"
	"time"
)

// Service определяет интерфейс для бизнес-логики управления содержимым юнитов.
type Service interface {
	CreateOrUpdateContent(ctx context.Context, unitID, itemID uuid.UUID, unitType string, quantity int) error
}

type service struct {
	repo repository.Repository
}

// NewService - конструктор для сервиса.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// CreateOrUpdateContent - основной метод, который создает или обновляет содержимое юнита.
// Он не знает о существовании "операций" или "задач", что делает модуль независимым.
func (s *service) CreateOrUpdateContent(ctx context.Context, unitID, itemID uuid.UUID, unitType string, quantity int) error {
	content := &models.UnitContent{
		ID:        uuid.New(), // ID нужен для новой записи, при обновлении он игнорируется в ON CONFLICT
		UnitID:    unitID,
		UnitType:  unitType,
		ItemID:    itemID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateOrUpdate(ctx, content)
}
