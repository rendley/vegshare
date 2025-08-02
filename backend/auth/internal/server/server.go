// Пакет `server` — это коллекция логики, связанной с запуском HTTP-сервера.
package server

import (
	"database/sql"
	"github.com/rendley/auth/internal/handler"
	"github.com/rendley/auth/pkg/config"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

// Определяем структуру `Server`. В Go структуры используются для объединения данных и методов.
// Здесь она хранит:
// - cfg: конфигурацию (порт, хост и т.д.),
// - handler: роутер (пока не используется, добавим позже).
type Server struct {
	cfg    *config.Config
	db     *sql.DB        // Добавляем поле для БД
	logger *logrus.Logger // Добавляем поле для логгера
}

// Функция `New` — это конструктор для `Server`.
// Принимает:
// - cfg: указатель на загруженный конфиг.
// - db: указатель на подкючение к базе.
// - log: указатель на логгер
// Возвращает:
// - Указатель на новый экземпляр `Server`.

func New(cfg *config.Config, db *sql.DB, logger *logrus.Logger) *Server {
	return &Server{
		cfg:    cfg, // Инициализируем поле `cfg` переданным конфигом.
		db:     db,
		logger: logger,
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
	// Передаём зависимости в handler.New()
	h := handler.New(s.db, s.logger)
	h.SetupRoutes(mux) // Регистрируем роуты

	// Запускаем сервер:
	// - `ListenAndServe` блокирует выполнение, пока сервер работает.
	// - Если произойдёт ошибка, она вернётся из функции.
	// - Можно передать `nil` вместо роутера (сервер будет возвращать 404 на все запросы).
	return http.ListenAndServe(addr, mux) // Принимает запросы и передаёт их в наш роутер.
}
