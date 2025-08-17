package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	farmModels "github.com/rendley/vegshare/backend/internal/farm/models"
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	"github.com/rendley/vegshare/backend/internal/operations/repository"
)

// Service определяет контракт для бизнес-логики.
type Service interface {
	PlantCrop(ctx context.Context, userID, plotID, cropID uuid.UUID) (*farmModels.PlotCrop, error)
	GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]farmModels.PlotCrop, error)
	RemoveCrop(ctx context.Context, plantingID uuid.UUID) error
	PerformAction(ctx context.Context, plotID uuid.UUID, action string) error
}

type service struct {
	repo        repository.Repository
	farmRepo    farmRepository.Repository
	leasingRepo leasingRepository.Repository
}

// NewOperationsService - конструктор для сервиса.
func NewOperationsService(repo repository.Repository, farmRepo farmRepository.Repository, leasingRepo leasingRepository.Repository) Service {
	return &service{
		repo:        repo,
		farmRepo:    farmRepo,
		leasingRepo: leasingRepo,
	}
}

func (s *service) PlantCrop(ctx context.Context, userID, plotID, cropID uuid.UUID) (*farmModels.PlotCrop, error) {
	// ... (previous implementation)
	leases, err := s.leasingRepo.GetLeasesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить аренду: %w", err)
	}
	var activeLease *farmModels.PlotLease
	for _, lease := range leases {
		if lease.PlotID == plotID && lease.Status == "active" {
			activeLease = &lease
			break
		}
	}
	if activeLease == nil {
		return nil, fmt.Errorf("у пользователя нет активной аренды для грядки %s", plotID)
	}

	_, err = s.farmRepo.GetCropByID(ctx, cropID)
	if err != nil {
		return nil, fmt.Errorf("культура с ID %s не найдена: %w", cropID, err)
	}

	now := time.Now()
	plotCrop := &farmModels.PlotCrop{
		ID:        uuid.New(),
		PlotID:    plotID,
		CropID:    cropID,
		LeaseID:   activeLease.ID,
		PlantedAt: now,
		Status:    "growing",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreatePlotCrop(ctx, plotCrop); err != nil {
		return nil, err
	}

	return plotCrop, nil
}

func (s *service) GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]farmModels.PlotCrop, error) {
	return s.repo.GetPlotCrops(ctx, plotID)
}

func (s *service) RemoveCrop(ctx context.Context, plantingID uuid.UUID) error {
	return s.repo.DeletePlotCrop(ctx, plantingID)
}

func (s *service) PerformAction(ctx context.Context, plotID uuid.UUID, action string) error {
	// TODO: Implement RabbitMQ logic
	fmt.Printf("Performing action '%s' on plot %s\n", action, plotID)
	return nil
}