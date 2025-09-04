package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
	"github.com/rendley/vegshare/backend/internal/catalog/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// --- DTOs ---

type CreateItemRequest struct {
	ItemType    string         `json:"item_type" validate:"required"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Attributes  models.JSONB   `json:"attributes"`
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

// GetItems handles GET requests to fetch catalog items, with filtering by type.
func (h *CatalogHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр 'type' из URL
	itemType := r.URL.Query().Get("type")
	if itemType == "" {
		api.RespondWithError(w, "query parameter 'type' is required", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetItems(r.Context(), itemType)
	if err != nil {
		h.logger.Errorf("ошибка при получении элементов каталога: %v", err)
		api.RespondWithError(w, "could not retrieve catalog items", http.StatusInternalServerError)
		return
	}
	api.RespondWithJSON(h.logger, w, items, http.StatusOK)
}

// CreateItem handles POST requests to create a new catalog item.
func (h *CatalogHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.service.CreateItem(r.Context(), req.ItemType, req.Name, req.Description, req.Attributes)
	if err != nil {
		h.logger.Errorf("ошибка при создании элемента каталога: %v", err)
		api.RespondWithError(w, "could not create catalog item", http.StatusInternalServerError)
		return
	}

	api.RespondWithJSON(h.logger, w, item, http.StatusCreated)
}