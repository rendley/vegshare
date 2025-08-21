package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	farmService "github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/rendley/vegshare/backend/internal/plot/repository"
)

// Service defines the contract for the plot service.
type Service interface {
	CreatePlot(ctx context.Context, name, size string, greenhouseID uuid.UUID) (*models.Plot, error)
	GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error)
	GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error)
	UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status string) (*models.Plot, error)
	DeletePlot(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface.
type service struct {
	repo    repository.Repository
	farmSvc farmService.Service
}

// NewService is a constructor for the plot service.
func NewService(repo repository.Repository, farmSvc farmService.Service) Service {
	return &service{repo: repo, farmSvc: farmSvc}
}

func (s *service) CreatePlot(ctx context.Context, name, size string, greenhouseID uuid.UUID) (*models.Plot, error) {
	// Check if the greenhouse exists
	_, err := s.farmSvc.GetGreenhouseByID(ctx, greenhouseID)
	if err != nil {
		return nil, fmt.Errorf("теплица с ID %s не найдена: %w", greenhouseID, err)
	}

	now := time.Now()
	plot := &models.Plot{
		ID:           uuid.New(),
		Name:         name,
		Size:         size,
		GreenhouseID: greenhouseID,
		Status:       "available",
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

func (s *service) UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status string) (*models.Plot, error) {
	plot, err := s.repo.GetPlotByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("грядка для обновления не найдена: %w", err)
	}

	plot.Name = name
	plot.Size = size
	plot.Status = status
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
