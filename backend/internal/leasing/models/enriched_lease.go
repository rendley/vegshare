package models

import (
	cameraModels "github.com/rendley/vegshare/backend/internal/camera/models"
	catalogModels "github.com/rendley/vegshare/backend/internal/catalog/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/google/uuid"
)

// EnrichedContent представляет собой содержимое юнита с полной информацией о предмете.
type EnrichedContent struct {
	ID       uuid.UUID               `json:"id"`
	Quantity int                     `json:"quantity"`
	Item     catalogModels.CatalogItem `json:"item"`
}

// EnrichedPlot представляет собой грядку с привязанными к ней камерами и содержимым.
type EnrichedPlot struct {
	plotModels.Plot
	Cameras  []cameraModels.Camera `json:"cameras"`
	Contents []EnrichedContent     `json:"contents"`
}

// EnrichedLease - это обогащенная модель аренды, которая включает в себя
// полную информацию об арендованном юните (пока только грядке).
type EnrichedLease struct {
	Lease
	Plot *EnrichedPlot `json:"plot,omitempty"`
	// В будущем здесь могут быть другие типы юнитов, например, Coop *EnrichedCoop
}
