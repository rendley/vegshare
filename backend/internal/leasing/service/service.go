package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	"github.com/rendley/vegshare/backend/internal/leasing/repository"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
)

// Service определяет контракт для бизнес-логики аренды.
type Service interface {
	LeasePlot(ctx context.Context, userID, plotID uuid.UUID) (*models.PlotLease, error)
	GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error)
}

type service struct {
	db      *sqlx.DB
	repo    repository.Repository
	plotSvc plotService.Service
}

// NewLeasingService - конструктор для сервиса аренды.
func NewLeasingService(db *sqlx.DB, repo repository.Repository, plotSvc plotService.Service) Service {
	return &service{
		db:      db,
		repo:    repo,
		plotSvc: plotSvc,
	}
}

func (s *service) LeasePlot(ctx context.Context, userID, plotID uuid.UUID) (*models.PlotLease, error) {
	// Шаг 1: Получаем грядку из сервиса plot, чтобы проверить ее статус (вне транзакции).
	plot, err := s.plotSvc.GetPlotByID(ctx, plotID)
	if err != nil {
		return nil, fmt.Errorf("грядка с ID %s не найдена: %w", plotID, err)
	}

	// Шаг 2: Проверяем бизнес-правило: можно арендовать только доступную грядку.
	if plot.Status != "available" {
		return nil, fmt.Errorf("грядка %s недоступна для аренды, ее текущий статус: %s", plotID, plot.Status)
	}

	// Шаг 3: Подготавливаем новую запись аренды.
	now := time.Now()
	lease := &models.PlotLease{
		ID:        uuid.New(),
		PlotID:    plotID,
		UserID:    userID,
		StartDate: now,
		EndDate:   now.AddDate(0, 3, 0), // Аренда на 3 месяца вперед
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Шаг 4: Запускаем транзакцию
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	// Гарантируем откат транзакции в случае любой паники или ошибки
	defer tx.Rollback()

	// Шаг 5: Создаем транзакционные версии зависимостей
	leasingRepoTx := repository.NewRepository(tx)
	plotSvcTx := s.plotSvc.WithTx(tx)

	// Шаг 6: Создаем запись аренды в нашей таблице (внутри транзакции)
	if err := leasingRepoTx.CreateLease(ctx, lease); err != nil {
		return nil, fmt.Errorf("не удалось создать запись аренды в транзакции: %w", err)
	}

	// Шаг 7: Обновляем статус грядки через сервис plot (внутри транзакции)
	_, err = plotSvcTx.UpdatePlot(ctx, plot.ID, plot.Name, plot.Size, "rented")
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить статус грядки в транзакции: %w", err)
	}

	// Шаг 8: Если все успешно, коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось закоммитить транзакцию: %w", err)
	}

	return lease, nil
}

func (s *service) GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error) {
	return s.repo.GetLeasesByUserID(ctx, userID)
}
