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
	plothandler "github.com/rendley/vegshare/backend/internal/plot/handler"
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
	PlotHandler       *plothandler.PlotHandler
}

// New - это конструктор для `Server`.
func New(cfg *config.Config, mw *middleware.Middleware, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler, leasing *leasinghandler.LeasingHandler, ops *operationshandler.OperationsHandler, catalog *cataloghandler.CatalogHandler, camera *camerahandler.CameraHandler, plot *plothandler.PlotHandler) *Server {
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
		PlotHandler:       plot,
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

			// User routes
			r.Mount("/users", s.UserHandler.Routes())

			// Leasing routes
			r.Mount("/leasing", s.LeasingHandler.Routes())

			// Operations routes
			r.Mount("/operations", s.OperationsHandler.Routes())

			// Farm routes (Regions, LandParcels, Greenhouses)
			r.Route("/farm", func(r chi.Router) {
				r.Route("/regions", func(r chi.Router) {
					r.Post("/", s.FarmHandler.CreateRegion)
					r.Get("/", s.FarmHandler.GetAllRegions)
					r.Route("/{regionID}", func(r chi.Router) {
						r.Get("/", s.FarmHandler.GetRegionByID)
						r.Put("/", s.FarmHandler.UpdateRegion)
						r.Delete("/", s.FarmHandler.DeleteRegion)
						r.Get("/land-parcels", s.FarmHandler.GetLandParcelsByRegion)
						r.Post("/land-parcels", s.FarmHandler.CreateLandParcelForRegion)
					})
				})
				r.Route("/land-parcels", func(r chi.Router) {
					r.Route("/{parcelID}", func(r chi.Router) {
						r.Get("/", s.FarmHandler.GetLandParcelByID)
						r.Put("/", s.FarmHandler.UpdateLandParcel)
						r.Delete("/", s.FarmHandler.DeleteLandParcel)
						r.Get("/greenhouses", s.FarmHandler.GetGreenhousesByLandParcel)
						r.Post("/greenhouses", s.FarmHandler.CreateGreenhouseForLandParcel)
					})
				})
				r.Route("/greenhouses", func(r chi.Router) {
					r.Route("/{greenhouseID}", func(r chi.Router) {
						r.Get("/", s.FarmHandler.GetGreenhouseByID)
						r.Put("/", s.FarmHandler.UpdateGreenhouse)
						r.Delete("/", s.FarmHandler.DeleteGreenhouse)
					})
				})
			})

			// Plot routes (now a top-level resource)
			r.Route("/plots", func(r chi.Router) {
				r.Post("/", s.PlotHandler.CreatePlot)
				r.Get("/", s.PlotHandler.GetPlots) // Handles ?greenhouse_id=...
				r.Route("/{plotID}", func(r chi.Router) {
					r.Get("/", s.PlotHandler.GetPlotByID)
					r.Put("/", s.PlotHandler.UpdatePlot)
					r.Delete("/", s.PlotHandler.DeletePlot)
					// Nested camera routes
					r.Get("/cameras", s.CameraHandler.GetCamerasByPlotID)
					r.Post("/cameras", s.CameraHandler.CreateCamera)
				})
			})

			// Camera routes (for top-level access like delete)
			r.Route("/cameras", func(r chi.Router) {
				r.Delete("/{cameraID}", s.CameraHandler.DeleteCamera)
			})
		})
	})

	return http.ListenAndServe(addr, r)
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}