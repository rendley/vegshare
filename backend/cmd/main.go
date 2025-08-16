package main

import (
	"fmt"

	"github.com/rendley/vegshare/backend/internal/api"
	authHandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	farmHandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	farmRepository "github.com/rendley/vegshare/backend/internal/farm/repository"
	farmService "github.com/rendley/vegshare/backend/internal/farm/service"
	leasingHandler "github.com/rendley/vegshare/backend/internal/leasing/handler"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	leasingService "github.com/rendley/vegshare/backend/internal/leasing/service"
	userHandler "github.com/rendley/vegshare/backend/internal/user/handler"
	userRepository "github.com/rendley/vegshare/backend/internal/user/repository"
	userService "github.com/rendley/vegshare/backend/internal/user/service"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/database"
	"github.com/rendley/vegshare/backend/pkg/jwt"
	"github.com/rendley/vegshare/backend/pkg/logger"
	"github.com/rendley/vegshare/backend/pkg/security"
)

func main() {
	// 1. Загружаем конфиги (порт, секреты) из YAML.
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)

	// 2. Инициализируем логгер
	log := logger.New()
	log.Info("Starting application...")

	// Создаём password hasher
	hasher := security.NewBcryptHasher(10)

	// Инициализация JWT
	jwtGen := jwt.NewGenerator(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Подключаемся к PostgreSQL
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected")

	// --- Инициализация модулей ---

	// Модуль Auth
	authHandler := authHandler.NewAuthHandler(db, hasher, log, jwtGen)

	// Модуль User
	userRepository := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepository)
	userHandler := userHandler.NewUserHandler(userService, log)

	// Модуль Farm
	farmRepository := farmRepository.NewRepository(db)
	farmService := farmService.NewFarmService(farmRepository)
	farmHandler := farmHandler.NewFarmHandler(farmService, log)

	// Модуль Leasing
	leasingRepository := leasingRepository.NewRepository(db)
	leasingService := leasingService.NewLeasingService(leasingRepository, farmRepository)
	leasingHandler := leasingHandler.NewLeasingHandler(leasingService, log)

	// Создаем и запускаем сервер
	srv := api.New(cfg, authHandler, userHandler, farmHandler, leasingHandler)

	// Запускаем сервер.
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}