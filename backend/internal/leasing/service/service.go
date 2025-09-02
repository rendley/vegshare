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
	// CreateLease создает аренду для любого типа юнита.
	CreateLease(ctx context.Context, userID, unitID uuid.UUID, unitType models.UnitType) (*models.Lease, error)
	GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.Lease, error)
	// RegisterUnitManager регистрирует "менеджер" для определенного типа юнитов.
	RegisterUnitManager(unitType models.UnitType, manager domain.UnitManager)
}

// service реализует Service.
type service struct {
	db               *sqlx.DB
	repo             repository.Repository
	unitManagers     map[models.UnitType]domain.UnitManager // Наш реестр "плагинов"
}

// NewLeasingService - конструктор для нового, универсального сервиса аренды.
func NewLeasingService(db *sqlx.DB, repo repository.Repository) Service {
	return &service{
		db:               db,
		repo:             repo,
		unitManagers:     make(map[models.UnitType]domain.UnitManager), // Инициализируем реестр
	}
}

// RegisterUnitManager добавляет реализацию UnitManager в реестр.
func (s *service) RegisterUnitManager(unitType models.UnitType, manager domain.UnitManager) {
	s.unitManagers[unitType] = manager
}

// CreateLease - новая универсальная функция создания аренды.
func (s *service) CreateLease(ctx context.Context, userID, unitID uuid.UUID, unitType models.UnitType) (*models.Lease, error) {
	// Шаг 1: Находим нужный менеджер в реестре.
	manager, ok := s.unitManagers[unitType]
	if !ok {
		return nil, fmt.Errorf("менеджер для типа юнита '%s' не зарегистрирован", unitType)
	}

	// Шаг 2: Получаем юнит и проверяем его статус (вне транзакции).
	unit, err := manager.GetLeasableUnit(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("юнит с ID %s не найден: %w", unitID, err)
	}
	if unit.GetStatus() != "available" {
		return nil, fmt.Errorf("юнит %s недоступен для аренды", unitID)
	}

	// Шаг 3: Подготавливаем новую универсальную запись аренды.
	now := time.Now()
	lease := &models.Lease{
		ID:        uuid.New(),
		UnitID:    unitID,
		UnitType:  unitType,
		UserID:    userID,
		StartDate: now,
		EndDate:   now.AddDate(0, 3, 0), // Аренда на 3 месяца
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Шаг 4: Запускаем транзакцию.
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Шаг 5: Создаем транзакционные версии зависимостей.
	leasingRepoTx := repository.NewRepository(tx)
	unitManagerTx := manager.WithTx(tx)

	// Шаг 6: Создаем запись аренды в нашей таблице (внутри транзакции).
	if err := leasingRepoTx.CreateLease(ctx, lease); err != nil {
		return nil, fmt.Errorf("не удалось создать запись аренды: %w", err)
	}

	// Шаг 7: Обновляем статус юнита через универсальный менеджер (внутри транзакции).
	if err := unitManagerTx.UpdateUnitStatus(ctx, unitID, "rented"); err != nil {
		return nil, fmt.Errorf("не удалось обновить статус юнита: %w", err)
	}

	// Шаг 8: Если все успешно, коммитим транзакцию.
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось закоммитить транзакцию: %w", err)
	}

	return lease, nil
}

func (s *service) GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.Lease, error) {
	return s.repo.GetLeasesByUserID(ctx, userID)
}