package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	catalogService "github.com/rendley/vegshare/backend/internal/catalog/service"
	operations_repository "github.com/rendley/vegshare/backend/internal/operations/repository"
	plotService "github.com/rendley/vegshare/backend/internal/plot/service"
	"github.com/rendley/vegshare/backend/internal/task/models"
	"github.com/rendley/vegshare/backend/internal/task/repository"
	unitcontent_service "github.com/rendley/vegshare/backend/internal/unitcontent/service"
	"time"
)

// Service определяет интерфейс для бизнес-логики управления задачами.
type Service interface {
	CreateTask(ctx context.Context, operationID uuid.UUID, title, description string) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	AcceptTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error)
	CompleteTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error)
	FailTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error)
}

// service - реализация Service.
type service struct {
	db                 *sqlx.DB
	taskRepo           repository.Repository
	operationRepo      operations_repository.Repository
	unitContentService unitcontent_service.Service
	plotService        plotService.Service
	catalogService     catalogService.Service
}

// NewService - конструктор для сервиса задач.
func NewService(db *sqlx.DB, taskRepo repository.Repository, opRepo operations_repository.Repository, ucs unitcontent_service.Service, ps plotService.Service, cs catalogService.Service) Service {
	return &service{
		db:                 db,
		taskRepo:           taskRepo,
		operationRepo:      opRepo,
		unitContentService: ucs,
		plotService:        ps,
		catalogService:     cs,
	}
}

func (s *service) CreateTask(ctx context.Context, operationID uuid.UUID, title, description string) (*models.Task, error) {
	now := time.Now()
	descPtr := &description
	if description == "" {
		descPtr = nil
	}

	task := &models.Task{
		ID:          uuid.New(),
		OperationID: operationID,
		Status:      models.StatusNew,
		Title:       title,
		Description: descPtr,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("не удалось создать задачу в репозитории: %w", err)
	}

	return task, nil
}

func (s *service) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}

func (s *service) AcceptTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("задача с ID %s не найдена: %w", taskID, err)
	}

	if task.Status != models.StatusNew {
		return nil, fmt.Errorf("задачу можно взять в работу только из статуса 'new', текущий статус: '%s'", task.Status)
	}

	task.Status = models.StatusInProgress
	task.AssigneeID = &userID

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback() // Откат, если что-то пойдет не так

	taskRepoTx := repository.NewRepository(tx)
	opRepoTx := operations_repository.NewRepository(tx)

	if err := taskRepoTx.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("не удалось обновить задачу: %w", err)
	}

	if err := opRepoTx.UpdateOperationLogStatus(ctx, task.OperationID, "in_progress"); err != nil {
		return nil, fmt.Errorf("не удалось обновить статус операции: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось закоммитить транзакцию: %w", err)
	}

	return task, nil
}

// PlantActionParams - структура для парсинга параметров операции 'plant'.
type PlantActionParams struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

func (s *service) CompleteTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("задача с ID %s не найдена: %w", taskID, err)
	}

	if task.Status != models.StatusInProgress {
		return nil, fmt.Errorf("завершить задачу можно только из статуса 'in_progress', текущий статус: '%s'", task.Status)
	}
	if task.AssigneeID == nil || *task.AssigneeID != userID {
		return nil, fmt.Errorf("завершить задачу может только назначенный исполнитель")
	}

	operation, err := s.operationRepo.GetOperationLogByID(ctx, task.OperationID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти связанную операцию %s: %w", task.OperationID, err)
	}

	task.Status = models.StatusCompleted

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	taskRepoTx := repository.NewRepository(tx)
	opRepoTx := operations_repository.NewRepository(tx)

	if err := taskRepoTx.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	if err := opRepoTx.UpdateOperationLogStatus(ctx, task.OperationID, "completed"); err != nil {
		return nil, err
	}

	if operation.ActionType == "plant" {
		var params PlantActionParams
		if err := json.Unmarshal(operation.Parameters, &params); err != nil {
			return nil, fmt.Errorf("ошибка парсинга параметров для операции plant: %w", err)
		}
		if err := s.unitContentService.CreateOrUpdateContent(ctx, operation.UnitID, params.ItemID, operation.UnitType, params.Quantity); err != nil {
			return nil, fmt.Errorf("не удалось обновить содержимое юнита: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *service) FailTask(ctx context.Context, taskID, userID uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("задача с ID %s не найдена: %w", taskID, err)
	}

	if task.AssigneeID != nil && *task.AssigneeID != userID && task.Status == models.StatusInProgress {
		return nil, fmt.Errorf("отменить задачу может только назначенный исполнитель")
	}

	task.Status = models.StatusFailed

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	taskRepoTx := repository.NewRepository(tx)
	opRepoTx := operations_repository.NewRepository(tx)

	if err := taskRepoTx.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	if err := opRepoTx.UpdateOperationLogStatus(ctx, task.OperationID, "failed"); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}