package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/operations/models"
	"github.com/rendley/vegshare/backend/pkg/database"
)

// Repository определяет контракт для хранилища операций.
type Repository interface {
	CreateOperationLog(ctx context.Context, log *models.OperationLog) error
	GetOperationLogByID(ctx context.Context, logID uuid.UUID) (*models.OperationLog, error)
	GetOperationLogsForUnit(ctx context.Context, unitID uuid.UUID) ([]models.OperationLog, error)
	DeleteOperationLog(ctx context.Context, logID uuid.UUID) error
	UpdateOperationLogStatus(ctx context.Context, logID uuid.UUID, status string) error
}

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOperationLog(ctx context.Context, log *models.OperationLog) error {
	query := `INSERT INTO operation_log (id, unit_id, unit_type, user_id, action_type, parameters, status, executed_at, created_at, updated_at) 
	          VALUES (:id, :unit_id, :unit_type, :user_id, :action_type, :parameters, :status, :executed_at, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("не удалось создать запись в журнале операций: %w", err)
	}
	return nil
}

func (r *repository) GetOperationLogByID(ctx context.Context, logID uuid.UUID) (*models.OperationLog, error) {
	var log models.OperationLog
	query := `SELECT * FROM operation_log WHERE id = $1`
	err := r.db.GetContext(ctx, &log, query, logID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить операцию по ID: %w", err)
	}
	return &log, nil
}

func (r *repository) GetOperationLogsForUnit(ctx context.Context, unitID uuid.UUID) ([]models.OperationLog, error) {
	var logs []models.OperationLog
	query := `SELECT * FROM operation_log WHERE unit_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &logs, query, unitID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить журнал операций для юнита: %w", err)
	}
	return logs, nil
}

func (r *repository) DeleteOperationLog(ctx context.Context, logID uuid.UUID) error {
	query := `DELETE FROM operation_log WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, logID)
	if err != nil {
		return fmt.Errorf("не удалось удалить запись из журнала операций: %w", err)
	}
	return nil
}

func (r *repository) UpdateOperationLogStatus(ctx context.Context, logID uuid.UUID, status string) error {
	query := `UPDATE operation_log SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, logID)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус операции: %w", err)
	}
	return nil
}
