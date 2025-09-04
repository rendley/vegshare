// Пакет repository отвечает за доступ к данным, связанным с арендой.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	cameraModels "github.com/rendley/vegshare/backend/internal/camera/models"
	leaseModels "github.com/rendley/vegshare/backend/internal/leasing/models"
	plotModels "github.com/rendley/vegshare/backend/internal/plot/models"
	"github.com/rendley/vegshare/backend/pkg/database"
)

// Repository определяет контракт для хранилища данных аренды.
type Repository interface {
	CreateLease(ctx context.Context, lease *leaseModels.Lease) error
	GetEnrichedLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leaseModels.EnrichedLease, error)
	GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leaseModels.Lease, error)
}

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLease(ctx context.Context, lease *leaseModels.Lease) error {
	query := `INSERT INTO leases (id, unit_id, unit_type, user_id, start_date, end_date, status, created_at, updated_at) 
	          VALUES (:id, :unit_id, :unit_type, :user_id, :start_date, :end_date, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, lease)
	if err != nil {
		return fmt.Errorf("не удалось создать запись аренды: %w", err)
	}
	return nil
}

// GetEnrichedLeasesByUserID получает список аренд пользователя с полной информацией о грядках и камерах.
func (r *repository) GetEnrichedLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leaseModels.EnrichedLease, error) {
	query := `
        SELECT
            l.id AS "id",
            l.user_id AS "user_id",
            l.start_date AS "start_date",
            l.end_date AS "end_date",
            l.status AS "status",
            l.created_at AS "created_at",
            l.updated_at AS "updated_at",
            p.id AS "plot.id",
            p.name AS "plot.name",
            p.size AS "plot.size",
            p.status AS "plot.status",
            c.id AS "plot.cameras.id",
            c.name AS "plot.cameras.name",
            c.rtsp_path_name AS "plot.cameras.rtsp_path_name"
        FROM
            leases l
        JOIN
            plots p ON l.unit_id = p.id AND l.unit_type = 'plot'
        LEFT JOIN
            cameras c ON p.id = c.plot_id
        WHERE
            l.user_id = $1 AND l.status = 'active'
        ORDER BY
            l.created_at DESC, c.created_at ASC;
    `

	type flatLeaseData struct {
		ID               uuid.UUID      `db:"id"`
		UserID           uuid.UUID      `db:"user_id"`
		StartDate        time.Time      `db:"start_date"`
		EndDate          time.Time      `db:"end_date"`
		Status           string         `db:"status"`
		CreatedAt        time.Time      `db:"created_at"`
		UpdatedAt        time.Time      `db:"updated_at"`
		PlotID           uuid.UUID      `db:"plot.id"`
		PlotName         string         `db:"plot.name"`
		PlotSize         string         `db:"plot.size"`
		PlotStatus       string         `db:"plot.status"`
		CameraID         uuid.NullUUID  `db:"plot.cameras.id"`
		CameraName       sql.NullString `db:"plot.cameras.name"`
		CameraRTSPPath   sql.NullString `db:"plot.cameras.rtsp_path_name"`
	}

	var flatData []flatLeaseData
	if err := r.db.SelectContext(ctx, &flatData, query, userID); err != nil {
		return nil, fmt.Errorf("не удалось получить обогащенный список аренд: %w", err)
	}

	// Группируем плоские данные в иерархическую структуру
	leasesMap := make(map[uuid.UUID]*leaseModels.EnrichedLease)
	for _, row := range flatData {
		if _, ok := leasesMap[row.ID]; !ok {
			leasesMap[row.ID] = &leaseModels.EnrichedLease{
				Lease: leaseModels.Lease{
					ID:        row.ID,
					UserID:    row.UserID,
					UnitID:    row.PlotID, // Важно: UnitID теперь PlotID
					UnitType:  "plot",
					StartDate: row.StartDate,
					EndDate:   row.EndDate,
					Status:    row.Status,
					CreatedAt: row.CreatedAt,
					UpdatedAt: row.UpdatedAt,
				},
				Plot: &leaseModels.EnrichedPlot{
					Plot: plotModels.Plot{
						ID:     row.PlotID,
						Name:   row.PlotName,
						Size:   row.PlotSize,
						Status: row.PlotStatus,
					},
					Cameras: []cameraModels.Camera{},
				},
			}
		}

		// Добавляем камеру, если она есть
		if row.CameraID.Valid {
			lease := leasesMap[row.ID]
			lease.Plot.Cameras = append(lease.Plot.Cameras, cameraModels.Camera{
				ID:           row.CameraID.UUID,
				PlotID:       row.PlotID,
				Name:         row.CameraName.String,
				RTSPPathName: row.CameraRTSPPath.String,
			})
		}
	}

	// Преобразуем карту в срез
	enrichedLeases := make([]leaseModels.EnrichedLease, 0, len(leasesMap))
	for _, lease := range leasesMap {
		enrichedLeases = append(enrichedLeases, *lease)
	}

	return enrichedLeases, nil
}
func (r *repository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leaseModels.Lease, error) {
	var leases []leaseModels.Lease
	query := `SELECT * FROM leases WHERE user_id = $1 AND status = 'active'`
	if err := r.db.SelectContext(ctx, &leases, query, userID); err != nil {
		return nil, fmt.Errorf("не удалось получить список аренд для пользователя: %w", err)
	}
	return leases, nil
}