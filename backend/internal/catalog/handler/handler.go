package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/catalog/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---

type CreateCropRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	PlantingTime int    `json:"planting_time"`
	HarvestTime  int    `json:"harvest_time"`
}

// --- Handler ---

// CatalogHandler handles HTTP requests for the catalog module.
type CatalogHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

// NewCatalogHandler creates a new CatalogHandler.
func NewCatalogHandler(s service.Service, l *logrus.Logger) *CatalogHandler {
	return &CatalogHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *CatalogHandler) GetAllCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := h.service.GetAllCrops(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка культур: %v", err)
		api.RespondWithError(w, "could not retrieve crops", http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(h.logger, w, crops, http.StatusOK)
}

func (h *CatalogHandler) CreateCrop(w http.ResponseWriter, r *http.Request) {
	var req CreateCropRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	crop, err := h.service.CreateCrop(r.Context(), req.Name, req.Description, req.PlantingTime, req.HarvestTime)
	if err != nil {
		h.logger.Errorf("ошибка при создании культуры: %v", err)
		api.RespondWithError(w, "could not create crop", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, crop, http.StatusCreated)
}
