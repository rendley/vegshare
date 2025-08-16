package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	"github.com/rendley/vegshare/backend/internal/leasing/repository"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Service определяет контракт для бизнес-логики аренды.
type Service interface {
	LeasePlot(ctx context.Context, userID, plotID uuid.UUID) (*models.PlotLease, error)
	GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error)
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

func (s *service) LeasePlot(ctx context.Context, userID, plotID uuid.UUID) (*models.PlotLease, error) {
	// Шаг 1: Получаем грядку из репозитория farm, чтобы проверить ее статус.
	// Это критически важный шаг, который использует зависимость от другого модуля.
	plot, err := s.farmRepo.GetPlotByID(ctx, plotID)
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

	// Шаг 4: Создаем запись аренды в нашей таблице.
	if err := s.repo.CreateLease(ctx, lease); err != nil {
		return nil, fmt.Errorf("не удалось создать запись аренды: %w", err)
	}

	// Шаг 5: Обновляем статус грядки в таблице plots.
	plot.Status = "rented"
	plot.UpdatedAt = time.Now()
	if err := s.farmRepo.UpdatePlot(ctx, plot); err != nil {
		// Здесь в реальном приложении нужно было бы откатить создание аренды,
		// но без транзакции мы этого сделать не можем. Пока оставляем так.
		// TODO: Обернуть всю операцию в транзакцию БД.
		return nil, fmt.Errorf("не удалось обновить статус грядки: %w", err)
	}

	// Шаг 6: Возвращаем созданную аренду.
	return lease, nil
}

func (s *service) GetMyLeases(ctx context.Context, userID uuid.UUID) ([]models.PlotLease, error) {
	return s.repo.GetLeasesByUserID(ctx, userID)
}