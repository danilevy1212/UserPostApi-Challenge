package inmemory

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
)

type PingFunc func(context.Context) error
type UserCreateFunc func(context.Context, models.User) (*models.User, error)
type UserGetAllFunc func(context.Context) ([]*models.User, error)
type UserGetByIDFunc func(context.Context, int) (*models.User, error)
type UserDeleteByIDFunc func(context.Context, int) error
type UserUpdateFunc func(context.Context, models.User) (*models.User, error)

var InMemoryDBPingFn PingFunc = func(c context.Context) error {
	return nil
}
var InMemoryUserCreateFn UserCreateFunc = func(ctx context.Context, user models.User) (*models.User, error) {
	id := 1
	return &models.User{
		ID:    &id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
var InMemoryUserGetAllFn UserGetAllFunc = func(ctx context.Context) ([]*models.User, error) {
	id1 := 1
	id2 := 2
	return []*models.User{
		{
			ID:    &id1,
			Name:  "John Doe",
			Email: "johnnydoe@gmail.com",
		},
		{
			ID:    &id2,
			Name:  "Daniel Levy Moreno",
			Email: "danielmorenolevy@gmail.com",
		},
	}, nil
}
var InMemoryUserGetByIDFn UserGetByIDFunc = func(ctx context.Context, id int) (*models.User, error) {
	return &models.User{
		ID:    &id,
		Name:  "Daniel Levy Moreno",
		Email: "danielmorenolevy@gmail.com",
	}, nil
}
var InMemoryUserDeleteByIDFn UserDeleteByIDFunc = func(ctx context.Context, i int) error {
	return nil
}
var InMemoryUserUpdateFn UserUpdateFunc = func(ctx context.Context, u models.User) (*models.User, error) {
	id := 1
	return &models.User{
		ID:    &id,
		Name:  "Daniel Levy Moreno",
		Email: "danielmorenolevy@gmail.com",
	}, nil
}

type InMemoryDB struct{}

func (im *InMemoryDB) Connection() *sql.DB {
	return nil
}

func (im *InMemoryDB) Ping(ctx context.Context) error {
	return InMemoryDBPingFn(ctx)
}

func (im *InMemoryDB) UserCreate(ctx context.Context, user models.User) (*models.User, error) {
	return InMemoryUserCreateFn(ctx, user)
}

func (im *InMemoryDB) UserGetAll(ctx context.Context) ([]*models.User, error) {
	return InMemoryUserGetAllFn(ctx)
}

func (im *InMemoryDB) UserGetByID(ctx context.Context, id int) (*models.User, error) {
	return InMemoryUserGetByIDFn(ctx, id)
}

func (im *InMemoryDB) UserDeleteByID(ctx context.Context, id int) error {
	return InMemoryUserDeleteByIDFn(ctx, id)
}

func (im *InMemoryDB) UserUpdate(ctx context.Context, user models.User) (*models.User, error) {
	return InMemoryUserUpdateFn(ctx, user)
}
