// Пакет service содержит бизнес-логику, связанную с фермами.
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/farm/models"
	"github.com/rendley/vegshare/backend/internal/farm/repository"
)

// Service - это интерфейс, определяющий контракт для сервиса фермы.
// Он описывает, какие операции бизнес-логики могут быть выполнены.
type Service interface {
	CreateFarm(ctx context.Context, name, location string) (*models.Farm, error)
	GetAllFarms(ctx context.Context) ([]models.Farm, error)
}

// service - это приватная структура, реализующая интерфейс Service.
// Она содержит зависимости, необходимые для выполнения бизнес-логики.
type service struct {
	// repo - это наша зависимость от хранилища. Обратите внимание,
	// мы зависим от ИНТЕРФЕЙСА repository.Repository, а не от конкретной реализации.
	// Это ключ к тестируемости.
	repo repository.Repository
}

// NewFarmService - это конструктор для нашего сервиса.
// Он принимает зависимости (в данном случае, репозиторий) и возвращает новый экземпляр сервиса.
func NewFarmService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// CreateFarm - это реализация метода по созданию фермы.
// Он принимает бизнес-данные (имя и локацию), а не готовую модель.
func (s *service) CreateFarm(ctx context.Context, name, location string) (*models.Farm, error) {
	// --- Начало бизнес-логики ---

	// 1. Подготовка модели для сохранения.
	// Мы создаем объект Farm, который будет сохранен в БД.
	farm := &models.Farm{
		// 2. Генерируем новый уникальный идентификатор (UUID) для фермы.
		// Это и есть бизнес-логика. Сервис отвечает за создание ID, а не хендлер или репозиторий.
		ID:       uuid.New(),
		Name:     name,
		Location: location,
	}

	// --- Конец бизнес-логики ---

	// 3. Вызов слоя данных (репозитория) для сохранения подготовленной модели.
	// Мы передаем управление нашему репозиторию, вызывая его метод CreateFarm.
	err := s.repo.CreateFarm(ctx, farm)
	if err != nil {
		// Если репозиторий вернул ошибку, мы ее "пробрасываем" наверх,
		// добавив свой контекст.
		return nil, fmt.Errorf("не удалось создать ферму в сервисе: %w", err)
	}

	// 4. Возвращаем созданный объект фермы.
	// Теперь он содержит сгенерированный ID.
	return farm, nil
}

// GetAllFarms вызывает репозиторий для получения всех ферм.
func (s *service) GetAllFarms(ctx context.Context) ([]models.Farm, error) {
	// Просто вызываем метод репозитория.
	farms, err := s.repo.GetAllFarms(ctx)
	if err != nil {
		// Если репозиторий вернул ошибку, мы ее "пробрасываем" наверх,
		// добавив свой контекст.
		return nil, fmt.Errorf("не удалось получить фермы в сервисе: %w", err)
	}

	return farms, nil
}
