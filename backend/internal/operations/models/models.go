package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OperationLog - это запись о любом действии, выполненном над юнитом.
type OperationLog struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	UnitID      uuid.UUID       `db:"unit_id" json:"unit_id"`
	UnitType    string          `db:"unit_type" json:"unit_type"`
	UserID      uuid.UUID       `db:"user_id" json:"user_id"`
	ActionType  string          `db:"action_type" json:"action_type"`
	Parameters  json.RawMessage `db:"parameters" json:"parameters"`
	Status      string          `db:"status" json:"status"`
	ExecutedAt  time.Time       `db:"executed_at" json:"executed_at"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
