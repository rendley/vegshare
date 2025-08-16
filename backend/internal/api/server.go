// Пакет `server` — это коллекция логики, связанной с запуском HTTP-сервера.
package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	authhandler "github.com/rendley/vegshare/backend/internal/auth/handler"
	farmhandler "github.com/rendley/vegshare/backend/internal/farm/handler"
	userhandler "github.com/rendley/vegshare/backend/internal/user/handler"
	"github.com/rendley/vegshare/backend/pkg/config"
)

// Определяем структуру `Server`. В Go структуры используются для объединения данных и методов.
// Здесь она хранит:
// - handler: роутер (пока не используется, добавим позже).
type Server struct {
	cfg         *config.Config
	AuthHandler *authhandler.AuthHandler
	UserHandler *userhandler.UserHandler
	FarmHandler *farmhandler.FarmHandler
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

func New(cfg *config.Config, auth *authhandler.AuthHandler, user *userhandler.UserHandler, farm *farmhandler.FarmHandler) *Server {
	return &Server{
		cfg:         cfg, // Инициализируем поле `cfg` переданным конфигом.
		AuthHandler: auth,
		UserHandler: user,
		FarmHandler: farm,
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

	// Создаем новый роутер chi
	r := chi.NewRouter()

	// Добавляем стандартные middleware для логирования и восстановления после паник
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Регистрируем маршруты наших хендлеров
	s.AuthHandler.RegisterRouter(r)
	s.UserHandler.RegisterRouter(r)
	s.FarmHandler.RegisterRouter(r)

	// Запускаем сервер:
	// - `ListenAndServe` блокирует выполнение, пока сервер работает.
	// - Если произойдёт ошибка, она вернётся из функции.
	return http.ListenAndServe(addr, r) // Принимает запросы и передаёт их в наш chi роутер.
}
