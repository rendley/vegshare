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

	// Structure methods
	CreateStructure(ctx context.Context, name, typeName string, landParcelID uuid.UUID) (*models.Structure, error)
	GetStructureByID(ctx context.Context, id uuid.UUID) (*models.Structure, error)
	GetStructuresByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Structure, error)
	UpdateStructure(ctx context.Context, id uuid.UUID, name, typeName string) (*models.Structure, error)
	DeleteStructure(ctx context.Context, id uuid.UUID) error
	GetStructureTypes(ctx context.Context) ([]string, error)
}

// service - это приватная структура, реализующая интерфейс Service.
type service struct {
	repo repository.Repository
}

// NewFarmService - это конструктор для нашего сервиса.
func NewFarmService(repo repository.Repository) Service {
	return &service{repo: repo}
}
