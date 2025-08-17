package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/leasing/service"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// LeasingHandler содержит зависимости для обработчиков аренды.
type LeasingHandler struct {
	service  service.Service
	logger   *logrus.Logger
	validate *validator.Validate
}

// NewLeasingHandler - конструктор для LeasingHandler.
func NewLeasingHandler(s service.Service, l *logrus.Logger) *LeasingHandler {
	return &LeasingHandler{
		service:  s,
		logger:   l,
		validate: validator.New(),
	}
}

func (h *LeasingHandler) LeasePlot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	plotIDStr := chi.URLParam(r, "plotID")
	plotID, err := uuid.Parse(plotIDStr)
	if err != nil {
		api.RespondWithError(w, "invalid plot ID in URL", http.StatusBadRequest)
		return
	}

	lease, err := h.service.LeasePlot(r.Context(), userID, plotID)
	if err != nil {
		h.logger.Errorf("ошибка при аренде грядки: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, lease, http.StatusCreated)
}

func (h *LeasingHandler) GetMyLeases(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	leases, err := h.service.GetMyLeases(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении списка аренд: %v", err)
		api.RespondWithError(w, "could not retrieve leases", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, leases, http.StatusOK)
}
