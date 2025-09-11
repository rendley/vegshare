package models

import (
	"github.com/google/uuid"
	"time"
)

type TaskStatus string

const (
	StatusNew        TaskStatus = "new"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	OperationID uuid.UUID  `json:"operation_id" db:"operation_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id" db:"assignee_id"`
	Status      TaskStatus `json:"status" db:"status"`
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
