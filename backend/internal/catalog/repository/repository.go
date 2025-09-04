package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/catalog/models"
)

// Repository defines the contract for the catalog data storage.
type Repository interface {
	GetItemsByType(ctx context.Context, itemType string) ([]models.CatalogItem, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.CatalogItem, error)
	CreateItem(ctx context.Context, item *models.CatalogItem) error
}

// repository implements the Repository interface.
type repository struct {
	db *sqlx.DB
}

// NewRepository is a constructor for the repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetItemsByType(ctx context.Context, itemType string) ([]models.CatalogItem, error) {
	var items []models.CatalogItem
	query := "SELECT * FROM catalog_items WHERE item_type = $1"
	err := r.db.SelectContext(ctx, &items, query, itemType)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить элементы каталога по типу '%s': %w", itemType, err)
	}
	return items, nil
}

func (r *repository) GetItemByID(ctx context.Context, id uuid.UUID) (*models.CatalogItem, error) {
	var item models.CatalogItem
	query := `SELECT * FROM catalog_items WHERE id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить элемент каталога по ID: %w", err)
	}
	return &item, nil
}

func (r *repository) CreateItem(ctx context.Context, item *models.CatalogItem) error {
	query := `INSERT INTO catalog_items (id, item_type, name, description, attributes, created_at, updated_at) 
			  VALUES (:id, :item_type, :name, :description, :attributes, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("не удалось создать элемент каталога: %w", err)
	}
	return nil
}