package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// --- Plot Methods ---

func (s *service) CreatePlot(ctx context.Context, name, size, cameraURL string, greenhouseID uuid.UUID) (*models.Plot, error) {
	// Проверяем, существует ли теплица
	_, err := s.repo.GetGreenhouseByID(ctx, greenhouseID)
	if err != nil {
		return nil, fmt.Errorf("теплица с ID %s не найдена: %w", greenhouseID, err)
	}

	now := time.Now()
	plot := &models.Plot{
		ID:           uuid.New(),
		Name:         name,
		Size:         size,
		CameraURL:    cameraURL,
		GreenhouseID: greenhouseID,
		Status:       "available", // Новая грядка всегда доступна
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.repo.CreatePlot(ctx, plot)
	if err != nil {
		return nil, err
	}

	return plot, nil
}

func (s *service) GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error) {
	return s.repo.GetPlotByID(ctx, id)
}

func (s *service) GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error) {
	return s.repo.GetPlotsByGreenhouse(ctx, greenhouseID)
}

func (s *service) UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status, cameraURL string) (*models.Plot, error) {
	plot, err := s.repo.GetPlotByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("грядка для обновления не найдена: %w", err)
	}

	plot.Name = name
	plot.Size = size
	plot.Status = status
	plot.CameraURL = cameraURL
	plot.UpdatedAt = time.Now()

	err = s.repo.UpdatePlot(ctx, plot)
	if err != nil {
		return nil, err
	}

	return plot, nil
}

func (s *service) DeletePlot(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePlot(ctx, id)
}
