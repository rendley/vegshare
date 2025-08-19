package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	authhandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	cataloghandler "github.com/rendley/vegshare/backend/internal/catalog/handler"
	farmhandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	leasinghandler "github.com/rendley/vegshare/backend/internal/leasing/handler"
	operationshandler "github.com/rendley/vegshare/backend/internal/operations/handler"
	userhandler "github.com/rendley/vegshare/backend/internal/user/handler"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/middleware"
)

// Server - это наша основная структура сервера, которая объединяет все зависимости.
type Server struct {
	cfg               *config.Config
	mw                *middleware.Middleware
	AuthHandler       *authhandler.AuthHandler
	UserHandler       *userhandler.UserHandler
	FarmHandler       *farmhandler.FarmHandler
	LeasingHandler    *leasinghandler.LeasingHandler
	OperationsHandler *operationshandler.OperationsHandler
	CatalogHandler    *cataloghandler.CatalogHandler
}

// New - это конструктор для `Server`.
func New(cfg *config.Config, mw *middleware.Middleware, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler, leasing *leasinghandler.LeasingHandler, ops *operationshandler.OperationsHandler, catalog *cataloghandler.CatalogHandler) *Server {
	return &Server{
		cfg:               cfg,
		mw:                mw,
		AuthHandler:       auth,
		UserHandler:       user,
		FarmHandler:       farm,
		LeasingHandler:    leasing,
		OperationsHandler: ops,
		CatalogHandler:    catalog,
	}
}

// Start запускает HTTP-сервер.
func (s *Server) Start() error {
	addr := s.cfg.HTTP.Host + ":" + s.cfg.HTTP.Port
	log.Printf("Starting server on %s", addr)

	r := chi.NewRouter()

	// Стандартные и CORS middleware
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(middleware.CorsMiddleware) // Добавляем наш CORS middleware
	r.Use(chi_middleware.RedirectSlashes)
	r.Use(chi_middleware.StripSlashes)

	// Публичные маршруты (аутентификация, каталог)
	r.Group(func(r chi.Router) {
		s.AuthHandler.RegisterRoutes(r)
		s.CatalogHandler.RegisterRoutes(r) // Маршруты каталога считаем публичными
	})

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(s.mw.AuthMiddleware) // Применяем Auth middleware

		s.UserHandler.RegisterRoutes(r)
		s.FarmHandler.RegisterRoutes(r)
		s.LeasingHandler.RegisterRoutes(r)
		s.OperationsHandler.RegisterRoutes(r)
	})

	return http.ListenAndServe(addr, r)
}