package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	operationsModels "github.com/rendley/vegshare/backend/internal/operations/models"
	"github.com/rendley/vegshare/backend/internal/operations/repository"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
)

// ActionRequest - это структура для запроса на создание нового действия.
type ActionRequest struct {
	UnitID     uuid.UUID       `json:"unit_id"`
	UnitType   string          `json:"unit_type"`
	ActionType string          `json:"action_type"`
	Parameters json.RawMessage `json:"parameters"`
}

// Service определяет контракт для бизнес-логики операций.
type Service interface {
	CreateAction(ctx context.Context, userID uuid.UUID, req ActionRequest) (*operationsModels.OperationLog, error)
	GetActionsForUnit(ctx context.Context, unitID uuid.UUID) ([]operationsModels.OperationLog, error)
	CancelAction(ctx context.Context, logID uuid.UUID) error
}

type service struct {
	repo        repository.Repository
	leasingRepo leasingRepository.Repository
	rabbitmq    rabbitmq.ClientInterface
	cfg         *config.Config
}

// NewOperationsService - конструктор для сервиса.
func NewOperationsService(repo repository.Repository, leasingRepo leasingRepository.Repository, rabbitmq rabbitmq.ClientInterface, cfg *config.Config) Service {
	return &service{
		repo:        repo,
		leasingRepo: leasingRepo,
		rabbitmq:    rabbitmq,
		cfg:         cfg,
	}
}

func (s *service) CreateAction(ctx context.Context, userID uuid.UUID, req ActionRequest) (*operationsModels.OperationLog, error) {
	// 1. Проверяем, что у пользователя есть активная аренда для данного юнита
	leases, err := s.leasingRepo.GetLeasesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить аренду: %w", err)
	}

	var hasActiveLease bool
	for _, lease := range leases {
		if lease.UnitID == req.UnitID && string(lease.UnitType) == req.UnitType && lease.Status == "active" {
			hasActiveLease = true
			break
		}
	}
	if !hasActiveLease {
		return nil, fmt.Errorf("у пользователя нет активной аренды для юнита %s", req.UnitID)
	}

	// 2. TODO: Валидация action_type и parameters (например, через отдельный сервис-валидатор)
	//    Сейчас мы доверяем данным, приходящим от клиента.

	// 3. Создаем запись в журнале операций
	now := time.Now()
	logEntry := &operationsModels.OperationLog{
		ID:         uuid.New(),
		UnitID:     req.UnitID,
		UnitType:   req.UnitType,
		UserID:     userID,
		ActionType: req.ActionType,
		Parameters: req.Parameters,
		Status:     "pending", // Начальный статус
		ExecutedAt: now,     // Можно установить в null и обновлять в воркере
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateOperationLog(ctx, logEntry); err != nil {
		return nil, err
	}

	// 4. Отправляем сообщение в RabbitMQ
	body, err := json.Marshal(logEntry)
	if err != nil {
		// В реальном приложении здесь нужна более сложная логика обработки ошибок,
		// возможно, откат транзакции создания записи в логе.
		return nil, fmt.Errorf("failed to marshal action message: %w", err)
	}

	queueName := s.cfg.RabbitMQ.Queues["actions"] // Универсальная очередь
	if err := s.rabbitmq.Publish(queueName, string(body)); err != nil {
		// Та же логика обработки ошибок
		return nil, fmt.Errorf("failed to publish action message: %w", err)
	}

	return logEntry, nil
}

func (s *service) GetActionsForUnit(ctx context.Context, unitID uuid.UUID) ([]operationsModels.OperationLog, error) {
	return s.repo.GetOperationLogsForUnit(ctx, unitID)
}

func (s *service) CancelAction(ctx context.Context, logID uuid.UUID) error {
	// TODO: Добавить логику отмены. 
	// Например, можно обновить статус в БД на 'cancelled' 
	// и, возможно, отправить компенсирующее событие в RabbitMQ.
	// Сейчас просто удаляем для простоты.
	return s.repo.DeleteOperationLog(ctx, logID)
}