package main

import (
	"fmt"
	"github.com/rendley/auth/internal/server"
	"github.com/rendley/auth/pkg/config"
	"log"
)

func main() {
	// Загружаем конфиги (порт, секреты) из YAML.
	// Функция `Load()` читает файл и парсит его в структуру `Config`.
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)

	//  Создаём экземпляр сервера, передавая ему конфиг.
	// `New()` — это конструктор, который инициализирует `Server`.
	srv := server.New(cfg)

	// Запускаем сервер.
	// Если `Start()` вернёт ошибку, программа завершится с логом.
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err) // `Fatalf` выводит сообщение и вызывает `os.Exit(1)`.
	}

}
