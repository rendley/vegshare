package models

import (
	cameraModels "github.com/rendley/vegshare/backend/internal/camera/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
)

// EnrichedPlot представляет собой грядку с привязанными к ней камерами.
// Эта структура используется для обогащения ответа API.	
type EnrichedPlot struct {
	plotModels.Plot
	Cameras []cameraModels.Camera `json:"cameras"`
}

// EnrichedLease - это обогащенная модель аренды, которая включает в себя
// полную информацию об арендованном юните (пока только грядке).
type EnrichedLease struct {
	Lease
	Plot *EnrichedPlot `json:"plot,omitempty"`
	// В будущем здесь могут быть другие типы юнитов, например, Coop *EnrichedCoop
}