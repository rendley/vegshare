package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/models"
	"github.com/rendley/vegshare/backend/internal/camera/repository"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
)

// Service defines the contract for the camera service.
type Service interface {
	CreateCamera(ctx context.Context, name, rtspPathName string, plotID uuid.UUID) (*models.Camera, error)
	GetCamerasByPlotID(ctx context.Context, plotID uuid.UUID) ([]models.Camera, error)
	GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error)
	DeleteCamera(ctx context.Context, cameraID uuid.UUID) error
}

// service implements the Service interface.
type service struct {
	repo    repository.Repository
	plotSvc plotService.Service
}

// NewService is a constructor for the camera service.
func NewService(repo repository.Repository, plotSvc plotService.Service) Service {
	return &service{repo: repo, plotSvc: plotSvc}
}

func (s *service) CreateCamera(ctx context.Context, name, rtspPathName string, plotID uuid.UUID) (*models.Camera, error) {
	// Check if the plot exists
	_, err := s.plotSvc.GetPlotByID(ctx, plotID)
	if err != nil {
		return nil, fmt.Errorf("грядка с ID %s не найдена: %w", plotID, err)
	}

	now := time.Now()
	camera := &models.Camera{
		ID:        uuid.New(),
		PlotID:    plotID,
		Name:      name,
		RTSPPathName:   rtspPathName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateCamera(ctx, camera); err != nil {
		return nil, err
	}

	return camera, nil
}

func (s *service) GetCamerasByPlotID(ctx context.Context, plotID uuid.UUID) ([]models.Camera, error) {
	// We should also check if the plot exists before getting cameras
	_, err := s.plotSvc.GetPlotByID(ctx, plotID)
	if err != nil {
		return nil, fmt.Errorf("грядка с ID %s не найдена: %w", plotID, err)
	}
	return s.repo.GetCamerasByPlotID(ctx, plotID)
}

func (s *service) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	return s.repo.DeleteCamera(ctx, cameraID)
}

func (s *service) GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error) {
	return s.repo.GetCameraByID(ctx, cameraID)
}