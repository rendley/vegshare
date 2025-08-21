// Пакет handler отвечает за обработку HTTP-запросов, связанных с фермами.
package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---

type RegionRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type LandParcelRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type GreenhouseRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Type string `json:"type" validate:"omitempty,min=2,max=50"`
}

type PlotRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Size      string `json:"size" validate:"omitempty,min=1,max=50"`
}

type CreateCropRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"omitempty"`
	PlantingTime int    `json:"planting_time" validate:"omitempty,gte=0"`
	HarvestTime  int    `json:"harvest_time" validate:"omitempty,gte=0"`
}

type UpdatePlotRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Size      string `json:"size" validate:"omitempty,min=1,max=50"`
	Status    string `json:"status" validate:"required,oneof=available rented maintenance"`
}

// --- Handler ---

// FarmHandler - это структура, которая содержит зависимости для обработчиков фермы.
type FarmHandler struct {
	service       service.Service
	logger        *logrus.Logger
	validate      *validator.Validate
}

// NewFarmHandler - конструктор для FarmHandler.
func NewFarmHandler(s service.Service, l *logrus.Logger) *FarmHandler {
	return &FarmHandler{
		service:       s,
		logger:        l,
		validate:      validator.New(),
	}
}
