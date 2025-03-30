package database

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
)

type DBRepository interface {
	Connection() *sql.DB
	Ping(ctx context.Context) error
	UserCreate(ctx context.Context, user models.User) (*models.User, error)
	UserGetAll(ctx context.Context) ([]*models.User, error)
	UserGetByID(ctx context.Context, id int) (*models.User, error)
	UserDeleteByID(ctx context.Context, id int) error
	UserUpdate(ctx context.Context, user models.User) (*models.User, error)
}
