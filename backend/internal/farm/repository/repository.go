// Пакет repository отвечает за прямой доступ к базе данных для сущностей фермы.
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository - это ИНТЕРФЕЙС, который определяет "контракт" для нашего хранилища.
type Repository interface {
	// Region methods
	CreateRegion(ctx context.Context, region *models.Region) error
	GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error)
	GetAllRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, region *models.Region) error
	DeleteRegion(ctx context.Context, id uuid.UUID) error
	RestoreRegion(ctx context.Context, id uuid.UUID) error
	GetAllRegionsIncludingDeleted(ctx context.Context) ([]models.Region, error)

	// LandParcel methods
	CreateLandParcel(ctx context.Context, parcel *models.LandParcel) error
	GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error)
	GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error)
	UpdateLandParcel(ctx context.Context, parcel *models.LandParcel) error
	DeleteLandParcel(ctx context.Context, id uuid.UUID) error

	// Structure methods
	CreateStructure(ctx context.Context, structure *models.Structure) error
	GetStructureByID(ctx context.Context, id uuid.UUID) (*models.Structure, error)
	GetStructuresByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Structure, error)
	UpdateStructure(ctx context.Context, structure *models.Structure) error
	DeleteStructure(ctx context.Context, id uuid.UUID) error
	GetStructureTypes(ctx context.Context) ([]string, error)
}

// repository - это СТРУКТУРА, которая реализует интерфейс Repository.
type repository struct {
	db *sqlx.DB // Подключение к базе данных.
}

// NewRepository - это функция-конструктор.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
