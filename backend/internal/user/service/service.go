// package service объявляет, что весь код в этом файле принадлежит пакету 'service'.
package service

import (
	"context" // Импортируем пакет context

	"github.com/google/uuid" // Импортируем пакет uuid для работы с ID
	"github.com/rendley/vegshare/backend/internal/user/models" // Импортируем наши модели
	// Импортируем наш пакет 'repository', чтобы сервис мог с ним взаимодействовать.
	"github.com/rendley/vegshare/backend/internal/user/repository"
)

// UserService - это структура, которая будет содержать бизнес-логику для работы с пользователями.
type UserService struct {
	// repo - это поле, в котором будет храниться экземпляр нашего UserRepository.
	// Сервис будет использовать его для доступа к данным.
	repo *repository.UserRepository
}

// NewUserService - это функция-конструктор для нашего сервиса.
// Она принимает на вход репозиторий, от которого зависит. Это называется "внедрение зависимостей" (Dependency Injection).
func NewUserService(repo *repository.UserRepository) *UserService {
	// Создаем экземпляр сервиса, сохраняем в него переданный репозиторий и возвращаем.
	return &UserService{repo: repo}
}

// GetUser - это метод сервиса для получения одного пользователя.
// Он принимает контекст и ID пользователя.
// Обратите внимание, он не знает про SQL или базу данных, он просто вызывает репозиторий.
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	// s.repo - это наш экземпляр репозитория, который мы сохранили при создании сервиса.
	// Мы вызываем его метод GetUserByID, передавая ему те же аргументы.
	// return возвращает результат вызова s.repo.GetUserByID() напрямую.
	// Таким образом, сервис делегирует работу с данными репозиторию.
	return s.repo.GetUserByID(ctx, id)
}

// UpdateUser - это метод сервиса для обновления данных пользователя.
// Он принимает ID пользователя и новые данные (имя и URL аватара).
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, fullName, avatarURL string) (*models.UserProfile, error) {
	// Шаг 1: Получаем текущего пользователя по ID.
	// Это гарантирует, что мы работаем с существующей записью.
	user, err := s.repo.GetUserByID(ctx, id)
	// Если репозиторий вернул ошибку (например, пользователь не найден),
	if err != nil {
		// мы прерываем выполнение и возвращаем эту ошибку наверх.
		return nil, err
	}

	// Шаг 2: Обновляем поля в полученной структуре 'user'.
	user.FullName = fullName
	user.AvatarURL = avatarURL

	// Шаг 3: Вызываем метод репозитория для сохранения обновленной структуры в БД.
	err = s.repo.UpdateUser(ctx, user)
	// Если при сохранении произошла ошибка,
	if err != nil {
		// также возвращаем ее.
		return nil, err
	}

	// Если все прошло успешно, возвращаем обновленный объект пользователя.
	return user, nil
}

// DeleteUser - это метод сервиса для удаления пользователя.
// Он просто вызывает соответствующий метод репозитория.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Вызываем метод DeleteUser из нашего репозитория, передавая ему ID.
	// и возвращаем результат (в данном случае, только ошибку, если она есть).
	return s.repo.DeleteUser(ctx, id)
}
