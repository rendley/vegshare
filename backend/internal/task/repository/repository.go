package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/task/models"
	"github.com/rendley/vegshare/backend/pkg/database"
	"time"
)

// Repository определяет интерфейс для взаимодействия с хранилищем задач.
type Repository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
}

// postgresRepository - реализация Repository для PostgreSQL.
type repository struct {
	db database.DBTX
}

// NewRepository - конструктор для postgresRepository.
func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (id, operation_id, status, title, description, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, task.ID, task.OperationID, task.Status, task.Title, task.Description, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *repository) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE id = $1`
	err := r.db.GetContext(ctx, &task, query, taskID)
	return &task, err
}

func (r *repository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM tasks ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &tasks, query)
	return tasks, err
}

func (r *repository) UpdateTask(ctx context.Context, task *models.Task) error {
	task.UpdatedAt = time.Now()
	query := `UPDATE tasks SET status = :status, assignee_id = :assignee_id, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}