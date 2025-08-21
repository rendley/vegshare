package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	authhandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	cataloghandler "github.com/rendley/vegshare/backend/internal/catalog/handler"
	farmhandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	leasinghandler "github.com/rendley/vegshare/backend/internal/leasing/handler"
	operationshandler "github.com/rendley/vegshare/backend/internal/operations/handler"
	camerahandler "github.com/rendley/vegshare/backend/internal/camera/handler"
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
	CameraHandler     *camerahandler.CameraHandler
}

// New - это конструктор для `Server`.
func New(cfg *config.Config, mw *middleware.Middleware, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler, leasing *leasinghandler.LeasingHandler, ops *operationshandler.OperationsHandler, catalog *cataloghandler.CatalogHandler, camera *camerahandler.CameraHandler) *Server {
	return &Server{
		cfg:               cfg,
		mw:                mw,
		AuthHandler:       auth,
		UserHandler:       user,
		FarmHandler:       farm,
		LeasingHandler:    leasing,
		OperationsHandler: ops,
		CatalogHandler:    catalog,
		CameraHandler:     camera,
	}
}

// Start запускает HTTP-сервер.
func (s *Server) Start() error {
	addr := s.cfg.HTTP.Host + ":" + s.cfg.HTTP.Port
	log.Printf("Starting server on %s", addr)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Logger) // Logs the start and end of each request with the path and duration
	r.Use(chi_middleware.Recoverer) // Recovers from panics without crashing server
	r.Use(middleware.CorsMiddleware)   // Handles CORS headers
	r.Use(chi_middleware.RedirectSlashes) // Redirects slashes on paths
	r.Use(chi_middleware.StripSlashes)   // Strips slashes from paths

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(chi_middleware.Timeout(60 * time.Second))

	// Health check endpoint
	r.Get("/api/v1/health", s.healthCheckHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Mount("/auth", s.AuthHandler.Routes())
		r.Mount("/catalog", s.CatalogHandler.Routes())

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.mw.AuthMiddleware)
			r.Mount("/users", s.UserHandler.Routes())
			r.Mount("/farm", s.FarmHandler.Routes())
			r.Mount("/leasing", s.LeasingHandler.Routes())
			r.Mount("/operations", s.OperationsHandler.Routes())

			// Отдельный маршрут для удаления камеры по ее ID
			r.Delete("/cameras/{cameraID}", s.CameraHandler.DeleteCamera)
		})
	})

	return http.ListenAndServe(addr, r)
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}