// Пакет service содержит бизнес-логику, связанную с фермами.
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/rendley/vegshare/backend/internal/farm/repository"
)

// Service - это интерфейс, определяющий контракт для сервиса фермы.
type Service interface {
	

	// Region methods
	CreateRegion(ctx context.Context, name string) (*models.Region, error)
	GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error)
	GetAllRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, id uuid.UUID, name string) (*models.Region, error)
	DeleteRegion(ctx context.Context, id uuid.UUID) error

	// LandParcel methods
	CreateLandParcel(ctx context.Context, name string, regionID uuid.UUID) (*models.LandParcel, error)
	GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error)
	GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error)
	UpdateLandParcel(ctx context.Context, id uuid.UUID, name string) (*models.LandParcel, error)
	DeleteLandParcel(ctx context.Context, id uuid.UUID) error

	// Greenhouse methods
	CreateGreenhouse(ctx context.Context, name, typeName string, landParcelID uuid.UUID) (*models.Greenhouse, error)
	GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error)
	GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error)
	UpdateGreenhouse(ctx context.Context, id uuid.UUID, name, typeName string) (*models.Greenhouse, error)
	DeleteGreenhouse(ctx context.Context, id uuid.UUID) error

	// Plot methods
	CreatePlot(ctx context.Context, name, size, cameraURL string, greenhouseID uuid.UUID) (*models.Plot, error)
	GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error)
	GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error)
	UpdatePlot(ctx context.Context, id uuid.UUID, name, size, status, cameraURL string) (*models.Plot, error)
	DeletePlot(ctx context.Context, id uuid.UUID) error
}

// service - это приватная структура, реализующая интерфейс Service.
type service struct {
	repo repository.Repository
}

// NewFarmService - это конструктор для нашего сервиса.
func NewFarmService(repo repository.Repository) Service {
	return &service{repo: repo}
}
