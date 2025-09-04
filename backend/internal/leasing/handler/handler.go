package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	"github.com/rendley/vegshare/backend/internal/leasing/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/sirupsen/logrus"
)

// LeasingHandler содержит зависимости для обработчиков аренды.
type LeasingHandler struct {
	service  service.Service // Зависим от интерфейса из пакета service
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

// createLeaseRequest определяет структуру тела запроса на создание аренды.
type createLeaseRequest struct {
	UnitID   uuid.UUID         `json:"unit_id" validate:"required"`
	UnitType models.UnitType `json:"unit_type" validate:"required"`
}

// CreateLease - новый обработчик для создания универсальной аренды.
func (h *LeasingHandler) CreateLease(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req createLeaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем новый, универсальный метод сервиса
	lease, err := h.service.CreateLease(r.Context(), userID, req.UnitID, req.UnitType)
	if err != nil {
		h.logger.Errorf("ошибка при создании аренды: %v", err)
		api.RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, lease, http.StatusCreated)
}

func (h *LeasingHandler) GetMyLeases(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		api.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	leases, err := h.service.GetMyEnrichedLeases(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении обогащенного списка аренд: %v", err)
		api.RespondWithError(w, "could not retrieve leases", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, leases, http.StatusOK)
}