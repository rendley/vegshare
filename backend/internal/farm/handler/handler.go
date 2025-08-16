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

// --- Handler ---

// FarmHandler - это структура, которая содержит зависимости для обработчиков фермы.
type FarmHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

// NewFarmHandler - конструктор для FarmHandler.
func NewFarmHandler(s service.Service, l *logrus.Logger) *FarmHandler {
	return &FarmHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}
