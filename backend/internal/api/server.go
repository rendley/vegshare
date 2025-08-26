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
	streaminghandler "github.com/rendley/vegshare/backend/internal/streaming/handler"
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
	StreamingHandler  *streaminghandler.StreamingHandler
}

// New - это конструктор для `Server`.
func New(cfg *config.Config, mw *middleware.Middleware, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler, leasing *leasinghandler.LeasingHandler, ops *operationshandler.OperationsHandler, catalog *cataloghandler.CatalogHandler, camera *camerahandler.CameraHandler, plot *plothandler.PlotHandler, stream *streaminghandler.StreamingHandler) *Server {
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
		StreamingHandler:  stream,
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
		// --- УРОВЕНЬ 1: ПУБЛИЧНЫЕ РОУТЫ ---
		// Доступны всем, без аутентификации.
		r.Mount("/auth", s.AuthHandler.Routes())
		r.Mount("/catalog", s.CatalogHandler.Routes())

		// --- УРОВЕНЬ 2: РОУТЫ ДЛЯ ЛЮБОГО ЗАЛОГИНЕННОГО ПОЛЬЗОВАТЕЛЯ ---
		// К этой группе применяется AuthMiddleware. Все, что внутри, требует наличия валидного токена.
		r.Group(func(r chi.Router) {
			r.Use(s.mw.AuthMiddleware) // ПРОВЕРКА №1: Пользователь залогинен?

			// --- Роуты, доступные всем аутентифицированным пользователям ---

			// Пользовательские операции (аренда, действия на грядке и т.д.)
			r.Mount("/users", s.UserHandler.Routes())
			r.Mount("/leasing", s.LeasingHandler.Routes())
			r.Mount("/operations", s.OperationsHandler.Routes())

			// Роуты на ЧТЕНИЕ данных о ферме (просмотр иерархии)
			r.Get("/farm/regions", s.FarmHandler.GetAllRegions)
			r.Get("/farm/regions/{regionID}", s.FarmHandler.GetRegionByID)
			r.Get("/farm/regions/{regionID}/land-parcels", s.FarmHandler.GetLandParcelsByRegion)
			r.Get("/farm/land-parcels/{parcelID}", s.FarmHandler.GetLandParcelByID)
			r.Get("/farm/land-parcels/{parcelID}/greenhouses", s.FarmHandler.GetGreenhousesByLandParcel)
			r.Get("/farm/greenhouses/{greenhouseID}", s.FarmHandler.GetGreenhouseByID)
			r.Get("/plots", s.PlotHandler.GetPlots) // Handles ?greenhouse_id=...
			r.Get("/plots/{plotID}", s.PlotHandler.GetPlotByID)
			r.Get("/plots/{plotID}/cameras", s.CameraHandler.GetCamerasByPlotID)

			// --- УРОВЕНЬ 3: РОУТЫ ТОЛЬКО ДЛЯ АДМИНИСТРАТОРА ---
			// Вложенная группа, к которой дополнительно применяется AdminMiddleware.
			// Запрос должен пройти обе проверки: AuthMiddleware и AdminMiddleware.
			r.Group(func(r chi.Router) {
				r.Use(s.mw.AdminMiddleware) // ПРОВЕРКА №2: Является ли пользователь админом?

				// Роуты на СОЗДАНИЕ, ИЗМЕНЕНИЕ и УДАЛЕНИЕ иерархии фермы
				r.Route("/farm", func(r chi.Router) {
					r.Post("/regions", s.FarmHandler.CreateRegion)
					r.Route("/regions/{regionID}", func(r chi.Router) {
						r.Put("/", s.FarmHandler.UpdateRegion)
						r.Delete("/", s.FarmHandler.DeleteRegion)
						r.Post("/land-parcels", s.FarmHandler.CreateLandParcelForRegion)
					})
					r.Route("/land-parcels/{parcelID}", func(r chi.Router) {
						r.Put("/", s.FarmHandler.UpdateLandParcel)
						r.Delete("/", s.FarmHandler.DeleteLandParcel)
						r.Post("/greenhouses", s.FarmHandler.CreateGreenhouseForLandParcel)
					})
					r.Route("/greenhouses/{greenhouseID}", func(r chi.Router) {
						r.Put("/", s.FarmHandler.UpdateGreenhouse)
						r.Delete("/", s.FarmHandler.DeleteGreenhouse)
					})
				})

				// Роуты на СОЗДАНИЕ, ИЗМЕНЕНИЕ и УДАЛЕНИЕ грядок и камер
				r.Route("/plots", func(r chi.Router) {
					r.Post("/", s.PlotHandler.CreatePlot)
					r.Route("/{plotID}", func(r chi.Router) {
						r.Put("/", s.PlotHandler.UpdatePlot)
						r.Delete("/", s.PlotHandler.DeletePlot)
						r.Post("/cameras", s.CameraHandler.CreateCamera)
					})
				})
				r.Route("/cameras", func(r chi.Router) {
					r.Delete("/{cameraID}", s.CameraHandler.DeleteCamera)
				})
			})
		})

		// Streaming routes (HLS and WebSocket) with query param auth
		r.Group(func(r chi.Router) {
			r.Use(s.mw.QueryParamAuthMiddleware)
			r.Mount("/stream", s.StreamingHandler.Routes())
		})
	})

	return http.ListenAndServe(addr, r)
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}