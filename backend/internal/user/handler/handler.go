package handler

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Корневая структура для всех обработчиков модуля user
type UserHandler struct {
	db     *sql.DB
	logger *logrus.Logger
	//repository repository.UserRepository
	validate *validator.Validate
}

func NewUserHandler(db *sql.DB, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
		//repository: repository.NewUserRepository(db),
		validate: validator.New(),
	}
}

// GET /users/me
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста (должен быть установлен в auth middleware)
	// Возвращаем профиль в формате JSON
}

// PATCH /users/me
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON из тела запроса (только обновляемые поля)
	// Валидация данных
	// Частичное обновление (PATCH, а не PUT)
}

// DELETE /users/me
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Мягкое удаление (помечаем is_deleted=true)
	// Или полное удаление из БД
	// Очистка сессий/токенов
}

// POST /users/me/password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Проверка старого пароля
	// Валидация нового пароля
	// Хеширование и сохранение
}

// GET /users/me/avatar
func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	// Отдача файла или redirect на S3/CDN
}
