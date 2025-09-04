package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JSONB is a custom type for our JSONB database fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface, which allows this type to be written to the database.
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface, which allows this type to be read from the database.
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// CatalogItem представляет собой универсальную запись в каталоге.
type CatalogItem struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ItemType    string    `db:"item_type" json:"item_type"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Attributes  JSONB     `db:"attributes" json:"attributes"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}