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
		// r.Group() создает группу роутов, к которой можно применить общий middleware.
		// Все роуты внутри этой группы будут сначала проходить через s.mw.AuthMiddleware.
		r.Group(func(r chi.Router) {
			r.Use(s.mw.AuthMiddleware) // ПРОВЕРКА №1: Пользователь залогинен?

			// Обычные пользовательские роуты (просмотр профиля, аренда и т.д.)
			r.Mount("/users", s.UserHandler.Routes())
			r.Mount("/leasing", s.LeasingHandler.Routes())
			r.Mount("/operations", s.OperationsHandler.Routes())

			// --- Иерархия фермы: РЕГИОНЫ ---
			// r.Route() группирует роуты по общему префиксу, делая код чище.
			r.Route("/farm/regions", func(r chi.Router) {
				// GET /farm/regions - получить все регионы (доступно всем залогиненным)
				r.Get("/", s.FarmHandler.GetAllRegions)
				// r.With() применяет middleware только к ОДНОМУ следующему обработчику.
				// POST /farm/regions - создать регион (ПРОВЕРКА №2: только админ)
				r.With(s.mw.AdminMiddleware).Post("/", s.FarmHandler.CreateRegion)

				// Группа для роутов с параметром {regionID}
				r.Route("/{regionID}", func(r chi.Router) {
					r.Get("/", s.FarmHandler.GetRegionByID)
					r.Get("/land-parcels", s.FarmHandler.GetLandParcelsByRegion)
					// Действия над конкретным регионом доступны только админу
					r.With(s.mw.AdminMiddleware).Put("/", s.FarmHandler.UpdateRegion)
					r.With(s.mw.AdminMiddleware).Delete("/", s.FarmHandler.DeleteRegion)
					// Создание дочерней сущности - тоже только админу
					r.With(s.mw.AdminMiddleware).Post("/land-parcels", s.FarmHandler.CreateLandParcelForRegion)
				})
			})

			// --- Иерархия фермы: УЧАСТКИ ---
			r.Route("/farm/land-parcels/{parcelID}", func(r chi.Router) {
				r.Get("/", s.FarmHandler.GetLandParcelByID)
				r.Get("/structures", s.FarmHandler.GetStructuresByLandParcel)
				r.With(s.mw.AdminMiddleware).Put("/", s.FarmHandler.UpdateLandParcel)
				r.With(s.mw.AdminMiddleware).Delete("/", s.FarmHandler.DeleteLandParcel)
				r.With(s.mw.AdminMiddleware).Post("/structures", s.FarmHandler.CreateStructureForLandParcel)
			})

			// --- Иерархия фермы: СТРОЕНИЯ ---
			r.Route("/farm/structures", func(r chi.Router) {
				r.Get("/types", s.FarmHandler.GetStructureTypes)
				r.Route("/{structureID}", func(r chi.Router) {
					r.Get("/", s.FarmHandler.GetStructureByID)
					r.With(s.mw.AdminMiddleware).Put("/", s.FarmHandler.UpdateStructure)
					r.With(s.mw.AdminMiddleware).Delete("/", s.FarmHandler.DeleteStructure)
				})
			})

			// --- ГРЯДКИ (Plots) ---
			r.Route("/plots", func(r chi.Router) {
				// GET /plots?structure_id=... - получить грядки в строении (для всех)
				r.Get("/", s.PlotHandler.GetPlots)
				// POST /plots - создать грядку (только админ)
				r.With(s.mw.AdminMiddleware).Post("/", s.PlotHandler.CreatePlot)

				// Группа для конкретной грядки: /plots/{plotID}
				r.Route("/{plotID}", func(r chi.Router) {
					r.Get("/", s.PlotHandler.GetPlotByID)
					r.Get("/cameras", s.CameraHandler.GetCamerasByPlotID)
					r.With(s.mw.AdminMiddleware).Put("/", s.PlotHandler.UpdatePlot)
					r.With(s.mw.AdminMiddleware).Delete("/", s.PlotHandler.DeletePlot)
					r.With(s.mw.AdminMiddleware).Post("/cameras", s.CameraHandler.CreateCamera)
				})
			})

			// --- КАМЕРЫ (только удаление по ID) ---
			r.Route("/cameras", func(r chi.Router) {
				r.With(s.mw.AdminMiddleware).Delete("/{cameraID}", s.CameraHandler.DeleteCamera)
			})

			// --- Админ-панель: Управление пользователями ---
			r.Route("/admin/users", func(r chi.Router) {
				r.Use(s.mw.AdminMiddleware) // Защищаем все роуты в этой группе
				r.Get("/", s.UserHandler.GetAllUsers)
				r.Put("/{userID}/role", s.UserHandler.UpdateUserRole)
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