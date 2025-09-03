package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/pkg/api"
)

// --- Structure Handlers ---

func (h *FarmHandler) CreateStructureForLandParcel(w http.ResponseWriter, r *http.Request) {
	landParcelIDStr := chi.URLParam(r, "parcelID")
	landParcelID, err := uuid.Parse(landParcelIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
		return
	}

	var req StructureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	structure, err := h.service.CreateStructure(r.Context(), req.Name, req.Type, landParcelID)
	if err != nil {
		h.logger.Errorf("ошибка при создании строения: %v", err)
		api.RespondWithError(w, "could not create structure", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, structure, http.StatusCreated)
}

func (h *FarmHandler) GetStructuresByLandParcel(w http.ResponseWriter, r *http.Request) {
	landParcelIDStr := chi.URLParam(r, "parcelID")
	landParcelID, err := uuid.Parse(landParcelIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid land parcel ID in URL", http.StatusBadRequest)
		return
	}

	structures, err := h.service.GetStructuresByLandParcel(r.Context(), landParcelID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка строений: %v", err)
		api.RespondWithError(w, "could not retrieve structures", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, structures, http.StatusOK)
}

func (h *FarmHandler) GetStructureByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "structureID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid structure ID in URL", http.StatusBadRequest)
		return
	}

	structure, err := h.service.GetStructureByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при получении строения: %v", err)
		api.RespondWithError(w, "structure not found", http.StatusNotFound)
		return
	}

	api.RespondWithJSON(h.logger, w, structure, http.StatusOK)
}

func (h *FarmHandler) UpdateStructure(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "structureID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid structure ID in URL", http.StatusBadRequest)
		return
	}

	var req StructureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	structure, err := h.service.UpdateStructure(r.Context(), id, req.Name, req.Type)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении строения: %v", err)
		api.RespondWithError(w, "could not update structure", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, structure, http.StatusOK)
}

func (h *FarmHandler) DeleteStructure(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "structureID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		api.RespondWithError(w, "invalid structure ID in URL", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteStructure(r.Context(), id)
	if err != nil {
		h.logger.Errorf("ошибка при удалении строения: %v", err)
		api.RespondWithError(w, "could not delete structure", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FarmHandler) GetStructureTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetStructureTypes(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении типов строений: %v", err)
		api.RespondWithError(w, "could not retrieve structure types", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, types, http.StatusOK)
}