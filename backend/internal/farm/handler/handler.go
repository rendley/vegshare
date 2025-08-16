// Пакет handler отвечает за обработку HTTP-запросов, связанных с фермами.
package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

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

// GetAllCrops - это публичный метод-обработчик для эндпоинта получения списка культур.
func (h *FarmHandler) GetAllCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := h.service.GetAllCrops(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка культур: %v", err)
		api.RespondWithError(w, "could not retrieve crops", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, crops, http.StatusOK)
}