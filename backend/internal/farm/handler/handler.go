// Пакет handler отвечает за обработку HTTP-запросов, связанных с фермами.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/rendley/vegshare/backend/pkg/api"
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

// --- Crop Handlers ---

func (h *FarmHandler) GetAllCrops(w http.ResponseWriter, r *http.Request) {
	crops, err := h.service.GetAllCrops(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка культур: %v", err)
		api.RespondWithError(w, "could not retrieve crops", http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(h.logger, w, crops, http.StatusOK)
}

// --- Region Handlers ---

func (h *FarmHandler) CreateRegion(w http.ResponseWriter, r *http.Request) {
	var req RegionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	region, err := h.service.CreateRegion(r.Context(), req.Name)
	if err != nil {
		h.logger.Errorf("ошибка при создании региона: %v", err)
		api.RespondWithError(w, "could not create region", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, region, http.StatusCreated)
}

func (h *FarmHandler) GetRegionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	region, err := h.service.GetRegionByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при получении региона: %v", err)
		api.RespondWithError(w, "region not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, region, http.StatusOK)
}

func (h *FarmHandler) GetAllRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.service.GetAllRegions(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка регионов: %v", err)
		api.RespondWithError(w, "could not retrieve regions", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, regions, http.StatusOK)
}

func (h *FarmHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	var req RegionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	region, err := h.service.UpdateRegion(r.Context(), id, req.Name)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении региона: %v", err)
		api.RespondWithError(w, "could not update region", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, region, http.StatusOK)
}

func (h *FarmHandler) DeleteRegion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRegion(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при удалении региона: %v", err)
		api.RespondWithError(w, "could not delete region", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- LandParcel Handlers ---

func (h *FarmHandler) CreateLandParcel(w http.ResponseWriter, r *http.Request) {
	regionIDStr := chi.URLParam(r, "regionID")
	regionID, err := uuid.Parse(regionIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	var req LandParcelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel, err := h.service.CreateLandParcel(r.Context(), req.Name, regionID)
	if err != nil {
		h.logger.Errorf("ошибка при создании земельного участка: %v", err)
		api.RespondWithError(w, "could not create land parcel", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, parcel, http.StatusCreated)
}

func (h *FarmHandler) GetLandParcelsByRegion(w http.ResponseWriter, r *http.Request) {
	regionIDStr := chi.URLParam(r, "regionID") // Получаем ID региона из URL
	regionID, err := uuid.Parse(regionIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	parcels, err := h.service.GetLandParcelsByRegion(r.Context(), regionID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка земельных участков: %v", err)
		api.RespondWithError(w, "could not retrieve land parcels", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, parcels, http.StatusOK)
}

func (h *FarmHandler) GetLandParcelByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID", http.StatusBadRequest)
		return
	}

	parcel, err := h.service.GetLandParcelByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при получении земельного участка: %v", err)
		api.RespondWithError(w, "land parcel not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, parcel, http.StatusOK)
}

func (h *FarmHandler) UpdateLandParcel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID", http.StatusBadRequest)
		return
	}

	var req LandParcelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel, err := h.service.UpdateLandParcel(r.Context(), id, req.Name)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении земельного участка: %v", err)
		api.RespondWithError(w, "could not update land parcel", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, parcel, http.StatusOK)
}

func (h *FarmHandler) DeleteLandParcel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteLandParcel(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при удалении земельного участка: %v", err)
		api.RespondWithError(w, "could not delete land parcel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}