// Пакет repository отвечает за прямой доступ к базе данных для сущностей фермы.
package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/farm/models"
)

// Repository - это ИНТЕРФЕЙС, который определяет "контракт" для нашего хранилища.
// Любой сервис, который захочет работать с данными фермы, будет использовать этот интерфейс,
// а не конкретную структуру. Это позволяет нам подменять реализацию в тестах.
type Repository interface {
	// CreateFarm определяет метод для создания новой фермы.
	CreateFarm(ctx context.Context, farm *models.Farm) error
	// GetAllFarms определяет метод для получения всех ферм.
	GetAllFarms(ctx context.Context) ([]models.Farm, error)
}

// repository - это СТРУКТУРА, которая реализует интерфейс Repository.
// Она приватна для пакета (начинается с маленькой буквы), так как внешнему коду
// не нужно знать о ее существовании. Все взаимодействие идет через интерфейс.
type repository struct {
	db *sqlx.DB // Подключение к базе данных.
}

// NewRepository - это функция-конструктор.
// Обратите внимание, она возвращает интерфейс Repository, а не конкретную структуру.
func NewRepository(db *sqlx.DB) Repository {
	// Мы возвращаем указатель на нашу приватную структуру, но "под маской" интерфейса.
	return &repository{db: db}
}

// CreateFarm - это реализация метода интерфейса Repository.
// Приемник (r *repository) указывает, что этот метод принадлежит нашей приватной структуре.
func (r *repository) CreateFarm(ctx context.Context, farm *models.Farm) error {
	// Текст SQL-запроса для вставки новой записи в таблицу 'farms'.
	query := "INSERT INTO farms (id, name, location) VALUES (:id, :name, :location)"

	// Выполняем запрос, используя именованные параметры из структуры farm.
	_, err := r.db.NamedExecContext(ctx, query, farm)

	// Если возникает ошибка, "оборачиваем" ее для добавления контекста.
	if err != nil {
		return fmt.Errorf("не удалось создать ферму в репозитории: %w", err)
	}

	return nil
}

// GetAllFarms - это реализация метода для получения всех ферм из базы данных.
func (r *repository) GetAllFarms(ctx context.Context) ([]models.Farm, error) {
	// Создаем срез для хранения результатов.
	var farms []models.Farm

	// Текст SQL-запроса для выбора всех записей из таблицы 'farms'.
	query := "SELECT * FROM farms"

	// Выполняем запрос и сканируем результаты в срез структур.
	err := r.db.SelectContext(ctx, &farms, query)
	if err != nil {
		// В случае ошибки, возвращаем ее с дополнительным контекстом.
		return nil, fmt.Errorf("не удалось получить список ферм: %w", err)
	}

	return farms, nil
}
