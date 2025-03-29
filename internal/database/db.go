package database

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
)

type DBRepository interface {
	Connection() *sql.DB
	Ping(ctx context.Context) error
	UserCreate(ctx context.Context, user ent.User) (*ent.User, error)
	UserGetAll(ctx context.Context) ([]*ent.User, error)
	UserGetByID(ctx context.Context, id int) (*ent.User, error)
}
