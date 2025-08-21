package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	catalogModels "github.com/rendley/vegshare/backend/internal/catalog/models"
	catalogService "github.com/rendley/vegshare/backend/internal/catalog/service"
	leasingModels "github.com/rendley/vegshare/backend/internal/leasing/models"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	"github.com/rendley/vegshare/backend/internal/operations/repository"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
)

// Service определяет контракт для бизнес-логики.
type Service interface {
	PlantCrop(ctx context.Context, userID, plotID, cropID uuid.UUID) (*catalogModels.PlotCrop, error)
	GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]catalogModels.PlotCrop, error)
	RemoveCrop(ctx context.Context, plantingID uuid.UUID) error
	PerformAction(ctx context.Context, plotID uuid.UUID, action string) error
}

type service struct {
	repo           repository.Repository
	plotSvc        plotService.Service
	leasingRepo    leasingRepository.Repository
	catalogService catalogService.Service
	rabbitmq       rabbitmq.ClientInterface
}

// NewOperationsService - конструктор для сервиса.
func NewOperationsService(repo repository.Repository, plotSvc plotService.Service, leasingRepo leasingRepository.Repository, catalogService catalogService.Service, rabbitmq rabbitmq.ClientInterface) Service {
	return &service{
		repo:           repo,
		plotSvc:        plotSvc,
		leasingRepo:    leasingRepo,
		catalogService: catalogService,
		rabbitmq:       rabbitmq,
	}
}

func (s *service) PlantCrop(ctx context.Context, userID, plotID, cropID uuid.UUID) (*catalogModels.PlotCrop, error) {
	// ... (previous implementation)
	leases, err := s.leasingRepo.GetLeasesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить аренду: %w", err)
	}
	var activeLease *leasingModels.PlotLease
	for _, lease := range leases {
		if lease.PlotID == plotID && lease.Status == "active" {
			activeLease = &lease
			break
		}
	}
	if activeLease == nil {
		return nil, fmt.Errorf("у пользователя нет активной аренды для грядки %s", plotID)
	}

	_, err = s.catalogService.GetCropByID(ctx, cropID)
	if err != nil {
		return nil, fmt.Errorf("культура с ID %s не найдена: %w", cropID, err)
	}

	now := time.Now()
	plotCrop := &catalogModels.PlotCrop{
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

func (s *service) GetPlotCrops(ctx context.Context, plotID uuid.UUID) ([]catalogModels.PlotCrop, error) {
	return s.repo.GetPlotCrops(ctx, plotID)
}

func (s *service) RemoveCrop(ctx context.Context, plantingID uuid.UUID) error {
	return s.repo.DeletePlotCrop(ctx, plantingID)
}

type ActionMessage struct {
	PlotID uuid.UUID `json:"plot_id"`
	Action string    `json:"action"`
}

func (s *service) PerformAction(ctx context.Context, plotID uuid.UUID, action string) error {
	msg := ActionMessage{
		PlotID: plotID,
		Action: action,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal action message: %w", err)
	}

	return s.rabbitmq.Publish("actions", string(body))
}