package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

// --- Greenhouse Handlers ---

func (h *FarmHandler) CreateGreenhouseForLandParcel(w http.ResponseWriter, r *http.Request) {
	landParcelIDStr := chi.URLParam(r, "parcelID")
	landParcelID, err := uuid.Parse(landParcelIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
		return
	}

	var req GreenhouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	greenhouse, err := h.service.CreateGreenhouse(r.Context(), req.Name, req.Type, landParcelID)
	if err != nil {
		h.logger.Errorf("ошибка при создании теплицы: %v", err)
		api.RespondWithError(w, "could not create greenhouse", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, greenhouse, http.StatusCreated)
}

func (h *FarmHandler) GetGreenhousesByLandParcel(w http.ResponseWriter, r *http.Request) {
	landParcelIDStr := chi.URLParam(r, "parcelID")
	landParcelID, err := uuid.Parse(landParcelIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
		return
	}

	greenhouses, err := h.service.GetGreenhousesByLandParcel(r.Context(), landParcelID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка теплиц: %v", err)
		api.RespondWithError(w, "could not retrieve greenhouses", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, greenhouses, http.StatusOK)
}

func (h *FarmHandler) GetGreenhouseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "greenhouseID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse ID in URL", http.StatusBadRequest)
		return
	}

	greenhouse, err := h.service.GetGreenhouseByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при получении теплицы: %v", err)
		api.RespondWithError(w, "greenhouse not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, greenhouse, http.StatusOK)
}

func (h *FarmHandler) UpdateGreenhouse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "greenhouseID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse ID in URL", http.StatusBadRequest)
		return
	}

	var req GreenhouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	greenhouse, err := h.service.UpdateGreenhouse(r.Context(), id, req.Name, req.Type)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении теплицы: %v", err)
		api.RespondWithError(w, "could not update greenhouse", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, greenhouse, http.StatusOK)
}

func (h *FarmHandler) DeleteGreenhouse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "greenhouseID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid greenhouse ID in URL", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteGreenhouse(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при удалении теплицы: %v", err)
		api.RespondWithError(w, "could not delete greenhouse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}