package main

import (
	"fmt"
	"github.com/rendley/auth/internal/server"
	"github.com/rendley/auth/pkg/config"
	"github.com/rendley/auth/pkg/database"
	"github.com/rendley/auth/pkg/logger"
)

func main() {
	// 1. Загружаем конфиги (порт, секреты) из YAML.
	// Функция `Load()` читает файл и парсит его в структуру `Config`.
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)

	// 2. Инициализируем логгер
	log := logger.New()
	log.Info("Starting application...")

	// 3. Подключаемся к PostgreSQL
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected")

	//  Создаём экземпляр сервера, передавая ему конфиг.
	// `New()` — это конструктор, который инициализирует `Server`.
	srv := server.New(cfg, db, log)

	// Запускаем сервер.
	// Если `Start()` вернёт ошибку, программа завершится с логом.
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err) // `Fatalf` выводит сообщение и вызывает `os.Exit(1)`.
	}

}
