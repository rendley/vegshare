// Пакет handler отвечает за обработку HTTP-запросов, связанных с арендой.
package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/rendley/vegshare/backend/internal/leasing/service"
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
