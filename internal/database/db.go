package database

import (
	"context"
	"database/sql"
)

type DBRepository interface {
	Connection() *sql.DB
	Ping(ctx context.Context) error
}
