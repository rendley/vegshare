// Пакет handler отвечает за обработку HTTP-запросов, связанных с фермами.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/farm/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTO (Data Transfer Objects) ---

// CreateFarmRequest - это структура для данных, приходящих в теле запроса на создание фермы.
type CreateFarmRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Location string `json:"location" validate:"required,min=5"`
}

// --- Handler ---

// FarmHandler - это структура, которая содержит зависимости для обработчиков фермы.
type FarmHandler struct {
	// service - ссылка на сервис бизнес-логики.
	service service.Service
	// logger - логгер для записи информации о событиях и ошибках.
	logger *logrus.Logger
	// validate - экземпляр валидатора для проверки входящих данных.
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

// CreateFarm - это публичный метод-обработчик для эндпоинта создания фермы.
func (h *FarmHandler) CreateFarm(w http.ResponseWriter, r *http.Request) {
	var req CreateFarmRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	farm, err := h.service.CreateFarm(r.Context(), req.Name, req.Location)
	if err != nil {
		h.logger.Errorf("ошибка при создании фермы: %v", err)
		api.RespondWithError(w, "could not create farm", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, farm, http.StatusCreated)
}

// GetAllFarms - это публичный метод-обработчик для эндпоинта получения списка ферм.
func (h *FarmHandler) GetAllFarms(w http.ResponseWriter, r *http.Request) {
	farms, err := h.service.GetAllFarms(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка ферм: %v", err)
		api.RespondWithError(w, "could not retrieve farms", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, farms, http.StatusOK)
}
