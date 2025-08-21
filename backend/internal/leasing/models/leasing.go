package models

import (
	"time"
	"github.com/google/uuid"
)

// PlotLease представляет собой факт аренды грядки пользователем на определенный срок.
// Эта таблица связывает пользователей и грядки.
type PlotLease struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PlotID    uuid.UUID `db:"plot_id" json:"plot_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	// Status показывает состояние аренды: 'active', 'expired'.
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
