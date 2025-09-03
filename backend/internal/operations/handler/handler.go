package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/operations/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
)

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

func (h *OperationsHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req service.ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Add validation for the request

	logEntry, err := h.service.CreateAction(r.Context(), userID, req)
	if err != nil {
		h.logger.Errorf("ошибка при создании действия: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, logEntry, http.StatusCreated)
}

func (h *OperationsHandler) GetActionsForUnit(w http.ResponseWriter, r *http.Request) {
	unitIDStr := chi.URLParam(r, "unitID")
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid unit ID in URL", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetActionsForUnit(r.Context(), unitID)
	if err != nil {
		h.logger.Errorf("ошибка при получении журнала действий: %v", err)
		api.RespondWithError(w, "could not retrieve actions", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, logs, http.StatusOK)
}

func (h *OperationsHandler) CancelAction(w http.ResponseWriter, r *http.Request) {
	logIDStr := chi.URLParam(r, "actionID")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid action ID in URL", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelAction(r.Context(), logID); err != nil {
		h.logger.Errorf("ошибка при отмене действия: %v", err)
		api.RespondWithError(w, "could not cancel action", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
