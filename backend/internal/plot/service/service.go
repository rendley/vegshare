package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	farmService "github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/rendley/vegshare/backend/internal/leasing/domain"
	"github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/rendley/vegshare/backend/internal/plot/repository"
)

// Service определяет контракт для plot service
type Service interface {
	CreatePlot(ctx context.Context, name, size string, structureID uuid.UUID) (*models.Plot, error)
	GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error)
	GetPlotsByStructure(ctx context.Context, structureID uuid.UUID) ([]models.Plot, error)
	UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status string) (*models.Plot, error)
	DeletePlot(ctx context.Context, id uuid.UUID) error

	// Методы из domain.UnitManager
	GetLeasableUnit(ctx context.Context, unitID uuid.UUID) (domain.LeasableUnit, error)
	UpdateUnitStatus(ctx context.Context, unitID uuid.UUID, status string) error
	WithTx(tx *sqlx.Tx) domain.UnitManager
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

// WithTx создает новый сервис в рамках транзакции.
func (s *service) WithTx(tx *sqlx.Tx) domain.UnitManager { // <- ИСПРАВЛЕНО
	return &service{
		repo:    repository.NewRepository(tx),
		farmSvc: s.farmSvc,
	}
}

// ... (остальные методы без изменений) ...
func (s *service) CreatePlot(ctx context.Context, name, size string, structureID uuid.UUID) (*models.Plot, error) {
	_, err := s.farmSvc.GetStructureByID(ctx, structureID)
	if err != nil {
		return nil, fmt.Errorf("строение с ID %s не найдено: %w", structureID, err)
	}

	now := time.Now()
	plot := &models.Plot{
		ID:          uuid.New(),
		Name:        name,
		Size:        size,
		StructureID: structureID,
		Status:      "available",
		CreatedAt:   now,
		UpdatedAt:   now,
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

func (s *service) GetPlotsByStructure(ctx context.Context, structureID uuid.UUID) ([]models.Plot, error) {
	return s.repo.GetPlotsByStructure(ctx, structureID)
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

func (s *service) GetLeasableUnit(ctx context.Context, unitID uuid.UUID) (domain.LeasableUnit, error) {
	return s.GetPlotByID(ctx, unitID)
}

func (s *service) UpdateUnitStatus(ctx context.Context, unitID uuid.UUID, status string) error {
	plot, err := s.GetPlotByID(ctx, unitID)
	if err != nil {
		return err
	}
	_, err = s.UpdatePlot(ctx, unitID, plot.Name, plot.Size, status)
	return err
}