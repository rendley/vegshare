package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

// --- LandParcel Handlers ---

func (h *FarmHandler) CreateLandParcelForRegion(w http.ResponseWriter, r *http.Request) {
	regionIDStr := chi.URLParam(r, "regionID")
	regionID, err := uuid.Parse(regionIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID in URL", http.StatusBadRequest)
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
	regionIDStr := chi.URLParam(r, "regionID")
	regionID, err := uuid.Parse(regionIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid region ID in URL", http.StatusBadRequest)
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
	idStr := chi.URLParam(r, "parcelID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
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
	idStr := chi.URLParam(r, "parcelID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
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
	idStr := chi.URLParam(r, "parcelID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
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
