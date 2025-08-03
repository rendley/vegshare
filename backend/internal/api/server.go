// Пакет `server` — это коллекция логики, связанной с запуском HTTP-сервера.
package api

import (
	authhandler "github.com/rendley/backend/internal/auth/handler"
	"github.com/rendley/backend/pkg/config"
	"log"
	"net/http"
)

// Определяем структуру `Server`. В Go структуры используются для объединения данных и методов.
// Здесь она хранит:
// - handler: роутер (пока не используется, добавим позже).
type Server struct {
	cfg         *config.Config
	AuthHandler *authhandler.Handler
	//FarmHandler *farmhandler.Handler
}

// Функция `New` — это конструктор для `Server`.
// Принимает:
// - cfg: указатель на загруженный конфиг.
// - db: указатель на подкючение к базе.
// - log: указатель на логгер
// Возвращает:
// - Указатель на новый экземпляр `Server`.

//func New(cfg *config.Config, hasher security.PasswordHasher, db *sql.DB, logger *logrus.Logger) *Server {
//	return &Server{
//		cfg:            cfg, // Инициализируем поле `cfg` переданным конфигом.
//		db:             db,
//		logger:         logger,
//		passwordHasher: hasher,
//	}
//}

func New(cfg *config.Config, auth *authhandler.Handler) *Server {
	return &Server{
		cfg:         cfg, // Инициализируем поле `cfg` переданным конфигом.
		AuthHandler: auth,
		//FarmHandler: farm,
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
	s.AuthHandler.RegisterRouter(mux)
	// s.FarmHandler.RegisterRoutes(mux)

	// Запускаем сервер:
	// - `ListenAndServe` блокирует выполнение, пока сервер работает.
	// - Если произойдёт ошибка, она вернётся из функции.
	// - Можно передать `nil` вместо роутера (сервер будет возвращать 404 на все запросы).
	return http.ListenAndServe(addr, mux) // Принимает запросы и передаёт их в наш роутер.
}
