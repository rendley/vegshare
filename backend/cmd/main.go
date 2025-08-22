package main

import (
	"fmt"

	"github.com/rendley/vegshare/backend/internal/api"
	authHandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	authRepository "github.com/rendley/vegshare/backend/internal/auth/repository"
	authService "github.com/rendley/vegshare/backend/internal/auth/service"
	cameraHandler "github.com/rendley/vegshare/backend/internal/camera/handler"
	cameraRepository "github.com/rendley/vegshare/backend/internal/camera/repository"
	cameraService "github.com/rendley/vegshare/backend/internal/camera/service"
	catalogHandler "github.com/rendley/vegshare/backend/internal/catalog/handler"
	catalogRepository "github.com/rendley/vegshare/backend/internal/catalog/repository"
	catalogService "github.com/rendley/vegshare/backend/internal/catalog/service"
	farmHandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	farmService "github.com/rendley/vegshare/backend/internal/farm/service"
	leasingHandler "github.com/rendley/vegshare/backend/internal/leasing/handler"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	leasingService "github.com/rendley/vegshare/backend/internal/leasing/service"
	operationsHandler "github.com/rendley/vegshare/backend/internal/operations/handler"
	operationsRepository "github.com/rendley/vegshare/backend/internal/operations/repository"
	operationsService "github.com/rendley/vegshare/backend/internal/operations/service"
	plotHandler "github.com/rendley/vegshare/backend/internal/plot/handler"
	plotRepository "github.com/rendley/vegshare/backend/internal/plot/repository"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
	streamingHandler "github.com/rendley/vegshare/backend/internal/streaming/handler"
	streamingService "github.com/rendley/vegshare/backend/internal/streaming/service"
	userHandler "github.com/rendley/vegshare/backend/internal/user/handler"
	userRepository "github.com/rendley/vegshare/backend/internal/user/repository"
	userService "github.com/rendley/vegshare/backend/internal/user/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/database"
	"github.com/rendley/vegshare/backend/pkg/jwt"
	"github.com/rendley/vegshare/backend/pkg/logger"
	"github.com/rendley/vegshare/backend/pkg/middleware"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
	"github.com/rendley/vegshare/backend/pkg/security"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)

	log := logger.New()
	log.Info("Starting application...")

	hasher := security.NewBcryptHasher(10)
	jwtGen := jwt.NewGenerator(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	rabbitMQClient, err := rabbitmq.New(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQClient.Close()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected")

	// --- Инициализация модулей ---

	// Repositories
	authRepo := authRepository.NewAuthRepository(db)
	userRepo := userRepository.NewUserRepository(db)
	farmRepo := farmRepository.NewRepository(db)
	leasingRepo := leasingRepository.NewRepository(db)
	operationsRepo := operationsRepository.NewRepository(db)
	catalogRepo := catalogRepository.NewRepository(db)
	cameraRepo := cameraRepository.NewRepository(db)
	plotRepo := plotRepository.NewRepository(db)

	// Services
	authSvc := authService.NewAuthService(authRepo, hasher, jwtGen)
	userSvc := userService.NewUserService(userRepo)
	farmSvc := farmService.NewFarmService(farmRepo)
	plotSvc := plotService.NewService(plotRepo, farmSvc)
	leasingSvc := leasingService.NewLeasingService(leasingRepo, plotSvc)
	catalogSvc := catalogService.NewService(catalogRepo)
	operationsSvc := operationsService.NewOperationsService(operationsRepo, plotSvc, leasingRepo, catalogSvc, rabbitMQClient)
	cameraSvc := cameraService.NewService(cameraRepo, plotSvc)
	streamingSvc := streamingService.NewService(cameraSvc)

	// Middleware
	mw := middleware.NewMiddleware(cfg, log)

	// Handlers
	authHandler := authHandler.NewAuthHandler(authSvc, log)
	userHandler := userHandler.NewUserHandler(userSvc, log)
	cameraHandler := cameraHandler.NewCameraHandler(cameraSvc, log)
	plotHandler := plotHandler.NewPlotHandler(plotSvc, log)
	farmHandler := farmHandler.NewFarmHandler(farmSvc, log)
	leasingHandler := leasingHandler.NewLeasingHandler(leasingSvc, log)
	operationsHandler := operationsHandler.NewOperationsHandler(operationsSvc, log)
	catalogHandler := catalogHandler.NewCatalogHandler(catalogSvc, log)
	streamingHandler := streamingHandler.NewStreamingHandler(streamingSvc, log)

	// Создаем и запускаем сервер
	srv := api.New(cfg, mw, authHandler, userHandler, farmHandler, leasingHandler, operationsHandler, catalogHandler, cameraHandler, plotHandler, streamingHandler)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
