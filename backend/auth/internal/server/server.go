// Пакет `server` — это коллекция логики, связанной с запуском HTTP-сервера.
package server

import (
	"github.com/rendley/auth/internal/handler"
	"github.com/rendley/auth/pkg/config"
	"log"
	"net/http"
)

// Определяем структуру `Server`. В Go структуры используются для объединения данных и методов.
// Здесь она хранит:
// - cfg: конфигурацию (порт, хост и т.д.),
// - handler: роутер (пока не используется, добавим позже).
type Server struct {
	cfg *config.Config
}

// Функция `New` — это конструктор для `Server`.
// Принимает:
// - cfg: указатель на загруженный конфиг.
// Возвращает:
// - Указатель на новый экземпляр `Server`.

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg, // Инициализируем поле `cfg` переданным конфигом.
	}
}

// Метод `Start` запускает HTTP-сервер.
// Возвращает:
// - error: если сервер не смог стартовать.
func (s *Server) Start() error {
	// Формируем адрес сервера из конфига
	addr := s.cfg.HTTP.Host + ":" + s.cfg.HTTP.Port

	// Логгируем адрес для отладки.
	log.Printf("Starting server on %s", addr)

	mux := http.NewServeMux()
	h := handler.New()
	h.SetupRoutes(mux) // Регистрируем роуты

	// Запускаем сервер:
	// - `ListenAndServe` блокирует выполнение, пока сервер работает.
	// - Если произойдёт ошибка, она вернётся из функции.
	// - Можно передать `nil` вместо роутера (сервер будет возвращать 404 на все запросы).
	return http.ListenAndServe(addr, mux) // Принимает запросы и передаёт их в наш роутер.
}
