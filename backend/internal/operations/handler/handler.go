package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/operations/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---
type PlantCropRequest struct {
	CropID uuid.UUID `json:"crop_id" validate:"required"`
}

// --- Handler ---

type OperationsHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

func NewOperationsHandler(s service.Service, l *logrus.Logger) *OperationsHandler {
	return &OperationsHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *OperationsHandler) PlantCrop(w http.ResponseWriter, r *http.Request) {
	// TODO: Получить userID из JWT токена
	userID := uuid.MustParse("607526d4-b782-4c8f-95c3-5b4da1ed3312") // Заглушка с ID нашего тестового пользователя

	plotIDStr := chi.URLParam(r, "plotID")
	plotID, err := uuid.Parse(plotIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	var req PlantCropRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	plotCrop, err := h.service.PlantCrop(r.Context(), userID, plotID, req.CropID)
	if err != nil {
		h.logger.Errorf("ошибка при посадке культуры: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, plotCrop, http.StatusCreated)
}