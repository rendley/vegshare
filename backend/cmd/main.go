package main

import (
	"fmt"
	"github.com/rendley/backend/internal/api"
	authhandler "github.com/rendley/backend/internal/auth/handler"
	userhandler "github.com/rendley/backend/internal/user/handler"
	"github.com/rendley/backend/pkg/config"
	"github.com/rendley/backend/pkg/database"
	"github.com/rendley/backend/pkg/jwt"
	"github.com/rendley/backend/pkg/logger"
	"github.com/rendley/backend/pkg/security"
)

func main() {
	// 1. Загружаем конфиги (порт, секреты) из YAML.
	// Функция `Load()` читает файл и парсит его в структуру `Config`.
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

	//  Создаём экземпляр сервера, передавая ему конфиг.
	// `New()` — это конструктор, который инициализирует `Server`.
	//srv := api.New(cfg, hasher, db, log)

	authHandler := authhandler.NewAuthHandler(db, hasher, log, jwtGen)
	userHandler := userhandler.NewUserHandler(db, log)
	srv := api.New(cfg, authHandler, userHandler)

	// Запускаем сервер.
	// Если `Start()` вернёт ошибку, программа завершится с логом.
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err) // `Fatalf` выводит сообщение и вызывает `os.Exit(1)`.
	}

}
