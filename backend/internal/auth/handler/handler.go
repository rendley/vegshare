package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/auth/models"
	"github.com/rendley/vegshare/backend/internal/auth/service"
	"github.com/rendley/vegshare/backend/pkg/api"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	service  service.AuthService
	logger   *logrus.Logger
	validate *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(service service.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		service:  service,
		logger:   logger,
		validate: validator.New(),
	}
}

// registerHandler handles user registration.
func (h *AuthHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, tokens, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case models.ErrUserExists:
			api.RespondWithError(w, err.Error(), http.StatusConflict)
		default:
			api.RespondWithError(w, "Registration failed", http.StatusInternalServerError)
		}
		return
	}

	response := LoginResponse{
		UserID: user.ID.String(),
		TokenPair: TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
	api.RespondWithJSON(h.logger, w, response, http.StatusCreated)
}

// loginHandler handles user login.
func (h *AuthHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		api.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, tokens, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case models.ErrInvalidCredentials:
			api.RespondWithError(w, err.Error(), http.StatusUnauthorized)
		default:
			api.RespondWithError(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	response := LoginResponse{
		UserID: user.ID.String(),
		TokenPair: TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
	api.RespondWithJSON(h.logger, w, response, http.StatusOK)
}
