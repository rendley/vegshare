// Пакет `server` — это коллекция логики, связанной с запуском HTTP-сервера.
package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	authhandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	farmhandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	leasinghandler "github.com/rendley/vegshare/backend/internal/leasing/handler"
	userhandler "github.com/rendley/vegshare/backend/internal/user/handler"
	"github.com/rendley/vegshare/backend/pkg/config"
)

// Server - это наша основная структура сервера, которая объединяет все зависимости.
type Server struct {
	cfg            *config.Config
	AuthHandler    *authhandler.AuthHandler
	UserHandler    *userhandler.UserHandler
	FarmHandler    *farmhandler.FarmHandler
	LeasingHandler *leasinghandler.LeasingHandler
}

// New - это конструктор для `Server`.
func New(cfg *config.Config, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler, leasing *leasinghandler.LeasingHandler) *Server {
	return &Server{
		cfg:            cfg,
		AuthHandler:    auth,
		UserHandler:    user,
		FarmHandler:    farm,
		LeasingHandler: leasing,
	}
}

// Start запускает HTTP-сервер.
func (s *Server) Start() error {
	addr := s.cfg.HTTP.Host + ":" + s.cfg.HTTP.Port
	log.Printf("Starting server on %s", addr)

	r := chi.NewRouter()

	// Стандартные middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.StripSlashes)

	// Регистрируем маршруты всех наших хендлеров
	s.AuthHandler.RegisterRouter(r)
	s.UserHandler.RegisterRouter(r)
	s.FarmHandler.RegisterRouter(r)
	s.LeasingHandler.RegisterRouter(r)

	return http.ListenAndServe(addr, r)
}
