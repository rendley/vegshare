package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/camera/models"
	"github.com/rendley/vegshare/backend/internal/camera/repository"
)

// Service defines the contract for the camera service.
type Service interface {
	CreateCamera(ctx context.Context, name, rtspPathName string, unitID uuid.UUID, unitType string) (*models.Camera, error)
	GetCamerasByUnitID(ctx context.Context, unitID uuid.UUID, unitType string) ([]models.Camera, error)
	GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error)
	DeleteCamera(ctx context.Context, cameraID uuid.UUID) error
}

// service implements the Service interface.
type service struct {
	repo repository.Repository
}

// NewService is a constructor for the camera service.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCamera(ctx context.Context, name, rtspPathName string, unitID uuid.UUID, unitType string) (*models.Camera, error) {
	// TODO: В будущем здесь может быть проверка существования юнита через полиморфный интерфейс

	now := time.Now()
	camera := &models.Camera{
		ID:           uuid.New(),
		UnitID:       unitID,
		UnitType:     unitType,
		Name:         name,
		RTSPPathName: rtspPathName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateCamera(ctx, camera); err != nil {
		return nil, err
	}

	return camera, nil
}

func (s *service) GetCamerasByUnitID(ctx context.Context, unitID uuid.UUID, unitType string) ([]models.Camera, error) {
	// TODO: Проверка существования юнита и прав доступа
	return s.repo.GetCamerasByUnitID(ctx, unitID, unitType)
}

func (s *service) DeleteCamera(ctx context.Context, cameraID uuid.UUID) error {
	return s.repo.DeleteCamera(ctx, cameraID)
}

func (s *service) GetCameraByID(ctx context.Context, cameraID uuid.UUID) (*models.Camera, error) {
	return s.repo.GetCameraByID(ctx, cameraID)
}
