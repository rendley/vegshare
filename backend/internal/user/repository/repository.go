// package repository объявляет, что весь код в этом файле принадлежит пакету 'repository'.
// Пакеты в Go - это способ организации и разделения кода.
package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	// Импортируем пакет 'sqlx' для удобной работы с базой данных.
	// Он расширяет стандартный пакет 'database/sql'.
	"github.com/jmoiron/sqlx"

	// Импортируем наш пакет 'models', чтобы иметь доступ к структуре UserProfile.
	// Это пример того, как пакеты взаимодействуют: repository использует определения из models.
	"github.com/rendley/vegshare/backend/internal/user/models"
)

// UserRepository - это структура, которая будет содержать всю логику для работы с таблицей users.
// Мы называем ее "репозиторий пользователей".
type UserRepository struct {
	// db - это поле, в котором будет храниться объект для подключения к базе данных.
	// Тип *sqlx.DB означает "указатель на объект sqlx.DB".
	db *sqlx.DB
}

// NewUserRepository - это функция-конструктор, которая создает и возвращает новый экземпляр UserRepository.
// Она принимает на вход подключение к базе данных (db *sqlx.DB).
func NewUserRepository(db *sqlx.DB) *UserRepository {
	// return &UserRepository{db: db} создает новый объект UserRepository,
	// помещает в его поле 'db' переданное подключение и возвращает указатель на этот объект.
	return &UserRepository{db: db}
}

// GetUserByID - это МЕТОД структуры UserRepository. Метод - это функция, привязанная к конкретной структуре.
// Он принимает 'контекст' (для управления таймаутами/отменой) и 'id' пользователя.
// Возвращает он указатель на найденный профиль (*models.UserProfile) и ошибку (error), если что-то пошло не так.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	// Объявляем переменную 'user', в которую мы сложим результат из базы данных.
	// Ее тип - models.UserProfile, который мы импортировали из пакета models.
	var user models.UserProfile

	// query - это строка, содержащая наш SQL-запрос.
	// Мы выбираем все нужные поля из таблицы 'users'.
	// "WHERE id = $1" - это условие поиска. $1 - это специальный placeholder (заполнитель) для первого аргумента,
	// который мы передадим в запрос. Это защищает нас от SQL-инъекций.
	query := "SELECT id, email, full_name, avatar_url, farm_id, created_at, updated_at FROM users WHERE id = $1"

	// r.db.GetContext - это метод sqlx, который выполняет запрос и автоматически "раскладывает"
	// полученные из БД колонки по полям нашей структуры 'user' (сопоставляя имена колонок и теги `db:"..."`).
	// Мы передаем ему контекст, переменную для результата, сам запрос и значение для плейсхолдера $1 (в нашем случае, id).
	err := r.db.GetContext(ctx, &user, query, id)

	// Стандартная проверка на ошибку в Go. Если err не равен 'nil', значит, произошла ошибка.
	if err != nil {
		// Если ошибка есть (например, пользователь не найден), мы возвращаем nil вместо пользователя
		// и саму ошибку, чтобы слой выше (сервис) мог ее обработать.
		return nil, fmt.Errorf("не удалось найти пользователя по id: %w", err)
	}

	// Если все прошло успешно, возвращаем указатель на заполненную структуру 'user' и nil в качестве ошибки.
	return &user, nil
}

// UpdateUser - это метод для обновления данных пользователя в базе данных.
// Он принимает контекст и указатель на объект UserProfile с уже измененными данными.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.UserProfile) error {
	// В запросе UPDATE мы указываем, какие поля нужно обновить.
	// :full_name, :avatar_url, :id - это именованные плейсхолдеры.
	// sqlx автоматически подставит в них значения из одноименных полей структуры 'user', которую мы передадим.
	// Это очень удобно, так как не нужно перечислять аргументы по порядку, как с '$1, $2'.
	// Мы также обновляем поле updated_at на текущее время с помощью функции БД now().
	query := `UPDATE users SET 
				full_name = :full_name, 
				avatar_url = :avatar_url,
				updated_at = now()
			  WHERE id = :id`

	// r.db.NamedExecContext выполняет запрос с именованными плейсхолдерами.
	// Он берет структуру 'user', находит в ней поля 'full_name', 'avatar_url', 'id'
	// и подставляет их значения в запрос.
	_, err := r.db.NamedExecContext(ctx, query, user)

	// Если при выполнении запроса произошла ошибка,
	if err != nil {
		// мы возвращаем ее, обернув в наше сообщение для контекста.
		return fmt.Errorf("не удалось обновить пользователя: %w", err)
	}

	// Если ошибок нет, возвращаем nil.
	return nil
}

// DeleteUser - это метод для удаления пользователя из базы данных по его ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Составляем запрос на удаление.
	// "WHERE id = $1" указывает, какую именно запись нужно удалить.
	query := "DELETE FROM users WHERE id = $1"

	// r.db.ExecContext выполняет запрос, который не возвращает строк (как DELETE, INSERT, UPDATE).
	// Мы передаем ему контекст, сам запрос и ID пользователя для плейсхолдера $1.
	_, err := r.db.ExecContext(ctx, query, id)

	// Если при удалении произошла ошибка,
	if err != nil {
		// возвращаем ее с пояснением.
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}

	// Если все прошло гладко, возвращаем nil.
	return nil
}
