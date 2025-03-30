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
	UserUpdate(ctx context.Context, user models.UserUpdate) (*models.User, error)

	PostCreate(ctx context.Context, post models.Post) (*models.Post, error)
	PostGetAll(ctx context.Context) ([]*models.Post, error)
	PostGetByID(ctx context.Context, id int) (*models.Post, error)
	PostDeleteByID(ctx context.Context, id int) error
	PostUpdate(ctx context.Context, post models.PostUpdate) (*models.Post, error)
}
