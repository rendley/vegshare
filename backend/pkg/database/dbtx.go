package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBTX - это интерфейс, которому удовлетворяют *sqlx.DB и *sqlx.Tx.
// Он позволяет писать методы, которые могут работать как с подключением к БД, так и с транзакцией.
type DBTX interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
