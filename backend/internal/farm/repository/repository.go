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
	// Crop methods
	GetAllCrops(ctx context.Context) ([]models.Crop, error)

	// Region methods
	CreateRegion(ctx context.Context, region *models.Region) error
	GetRegionByID(ctx context.Context, id uuid.UUID) (*models.Region, error)
	GetAllRegions(ctx context.Context) ([]models.Region, error)
	UpdateRegion(ctx context.Context, region *models.Region) error
	DeleteRegion(ctx context.Context, id uuid.UUID) error

	// LandParcel methods
	CreateLandParcel(ctx context.Context, parcel *models.LandParcel) error
	GetLandParcelByID(ctx context.Context, id uuid.UUID) (*models.LandParcel, error)
	GetLandParcelsByRegion(ctx context.Context, regionID uuid.UUID) ([]models.LandParcel, error)
	UpdateLandParcel(ctx context.Context, parcel *models.LandParcel) error
	DeleteLandParcel(ctx context.Context, id uuid.UUID) error

	// Greenhouse methods
	CreateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error
	GetGreenhouseByID(ctx context.Context, id uuid.UUID) (*models.Greenhouse, error)
	GetGreenhousesByLandParcel(ctx context.Context, landParcelID uuid.UUID) ([]models.Greenhouse, error)
	UpdateGreenhouse(ctx context.Context, greenhouse *models.Greenhouse) error
	DeleteGreenhouse(ctx context.Context, id uuid.UUID) error

	// Plot methods
	CreatePlot(ctx context.Context, plot *models.Plot) error
	GetPlotByID(ctx context.Context, id uuid.UUID) (*models.Plot, error)
	GetPlotsByGreenhouse(ctx context.Context, greenhouseID uuid.UUID) ([]models.Plot, error)
	UpdatePlot(ctx context.Context, plot *models.Plot) error
	DeletePlot(ctx context.Context, id uuid.UUID) error
}

// repository - это СТРУКТУРА, которая реализует интерфейс Repository.
type repository struct {
	db *sqlx.DB // Подключение к базе данных.
}

// NewRepository - это функция-конструктор.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
