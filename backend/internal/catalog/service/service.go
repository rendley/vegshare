package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	"github.com/rendley/vegshare/backend/internal/catalog/repository"
)

// Service defines the contract for the catalog service.
type Service interface {
	GetAllCrops(ctx context.Context) ([]models.Crop, error)
	GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error)
	CreateCrop(ctx context.Context, name, description string, plantingTime, harvestTime int) (*models.Crop, error)
}

// service implements the Service interface.
type service struct {
	repo repository.Repository
}

// NewService is a constructor for the service.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllCrops(ctx context.Context) ([]models.Crop, error) {
	crops, err := s.repo.GetAllCrops(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить культуры в сервисе: %w", err)
	}
	return crops, nil
}

func (s *service) GetCropByID(ctx context.Context, id uuid.UUID) (*models.Crop, error) {
	return s.repo.GetCropByID(ctx, id)
}

func (s *service) CreateCrop(ctx context.Context, name, description string, plantingTime, harvestTime int) (*models.Crop, error) {
	now := time.Now()
	crop := &models.Crop{
		ID:           uuid.New(),
		Name:         name,
		Description:  description,
		PlantingTime: plantingTime,
		HarvestTime:  harvestTime,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateCrop(ctx, crop); err != nil {
		return nil, err
	}
	return crop, nil
}
