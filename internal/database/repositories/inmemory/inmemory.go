package inmemory

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
)

type PingFunc func(context.Context) error
type UserCreateFunc func(context.Context, ent.User) (*ent.User, error)
type UserGetAllFunc func(context.Context) ([]*ent.User, error)
type UserGetByIDFunc func(context.Context, int) (*ent.User, error)

var InMemoryDBPingFn PingFunc = func(c context.Context) error {
	return nil
}
var InMemoryUserCreateFn UserCreateFunc = func(ctx context.Context, user ent.User) (*ent.User, error) {
	return &ent.User{
		ID:    1,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
var InMemoryUserGetAllFn UserGetAllFunc = func(ctx context.Context) ([]*ent.User, error) {
	return []*ent.User{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "johnnydoe@gmail.com",
		},
		{
			ID:    2,
			Name:  "Daniel Levy Moreno",
			Email: "danielmorenolevy@gmail.com",
		},
	}, nil
}
var InMemoryUserGetByIDFn UserGetByIDFunc = func(ctx context.Context, id int) (*ent.User, error) {
	return &ent.User{
		ID:    1,
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

func (im *InMemoryDB) UserCreate(ctx context.Context, user ent.User) (*ent.User, error) {
	return InMemoryUserCreateFn(ctx, user)
}

func (im *InMemoryDB) UserGetAll(ctx context.Context) ([]*ent.User, error) {
	return InMemoryUserGetAllFn(ctx)
}

func (im *InMemoryDB) UserGetByID(ctx context.Context, id int) (*ent.User, error) {
	return InMemoryUserGetByIDFn(ctx, id)
}
