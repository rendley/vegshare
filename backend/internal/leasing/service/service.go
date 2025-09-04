package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/leasing/domain"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	"github.com/rendley/vegshare/backend/internal/leasing/repository"
)

// Service определяет контракт для бизнес-логики аренды.
type Service interface {
	CreateLease(ctx context.Context, userID, unitID uuid.UUID, unitType models.UnitType) (*models.Lease, error)
	GetMyEnrichedLeases(ctx context.Context, userID uuid.UUID) ([]models.EnrichedLease, error)
	GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Lease, error) // Added for internal use
	RegisterUnitManager(unitType models.UnitType, manager domain.UnitManager)
}

// service реализует Service.
type service struct {
	db               *sqlx.DB
	repo             repository.Repository
	unitManagers     map[models.UnitType]domain.UnitManager
}

// NewLeasingService - конструктор для нового, универсального сервиса аренды.
func NewLeasingService(db *sqlx.DB, repo repository.Repository) Service {
	return &service{
		db:               db,
		repo:             repo,
		unitManagers:     make(map[models.UnitType]domain.UnitManager),
	}
}

func (s *service) RegisterUnitManager(unitType models.UnitType, manager domain.UnitManager) {
	s.unitManagers[unitType] = manager
}

func (s *service) CreateLease(ctx context.Context, userID, unitID uuid.UUID, unitType models.UnitType) (*models.Lease, error) {
	manager, ok := s.unitManagers[unitType]
	if !ok {
		return nil, fmt.Errorf("менеджер для типа юнита '%s' не зарегистрирован", unitType)
	}

	unit, err := manager.GetLeasableUnit(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("юнит с ID %s не найден: %w", unitID, err)
	}
	if unit.GetStatus() != "available" {
		return nil, fmt.Errorf("юнит %s недоступен для аренды", unitID)
	}

	now := time.Now()
	lease := &models.Lease{
		ID:        uuid.New(),
		UnitID:    unitID,
		UnitType:  unitType,
		UserID:    userID,
		StartDate: now,
		EndDate:   now.AddDate(0, 3, 0),
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	leasingRepoTx := repository.NewRepository(tx)
	unitManagerTx := manager.WithTx(tx)

	if err := leasingRepoTx.CreateLease(ctx, lease); err != nil {
		return nil, fmt.Errorf("не удалось создать запись аренды: %w", err)
	}

	if err := unitManagerTx.UpdateUnitStatus(ctx, unitID, "rented"); err != nil {
		return nil, fmt.Errorf("не удалось обновить статус юнита: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось закоммитить транзакцию: %w", err)
	}

	return lease, nil
}

// GetMyEnrichedLeases вызывает репозиторий для получения обогащенных данных.
func (s *service) GetMyEnrichedLeases(ctx context.Context, userID uuid.UUID) ([]models.EnrichedLease, error) {
	return s.repo.GetEnrichedLeasesByUserID(ctx, userID)
}

func (s *service) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Lease, error) {
	return s.repo.GetLeasesByUserID(ctx, userID)
}