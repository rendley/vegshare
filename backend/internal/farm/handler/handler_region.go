package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

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
	idStr := chi.URLParam(r, "regionID")
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
	idStr := chi.URLParam(r, "regionID")
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
	idStr := chi.URLParam(r, "regionID")
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

// --- Admin Handlers for Regions ---

func (h *FarmHandler) GetAllRegionsIncludingDeleted(w http.ResponseWriter, r *http.Request) {
	regions, err := h.service.GetAllRegionsIncludingDeleted(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении полного списка регионов: %v", err)
		api.RespondWithError(w, "could not retrieve regions", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, regions, http.StatusOK)
}

func (h *FarmHandler) RestoreRegion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "regionID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID", http.StatusBadRequest)
		return
	}

	err = h.service.RestoreRegion(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при восстановлении региона: %v", err)
		api.RespondWithError(w, "could not restore region", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
