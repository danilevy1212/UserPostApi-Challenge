package inmemory

import (
	"context"
	"database/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
)

type PingFunc func(context.Context) error
type UserCreateFunc func(context.Context, ent.User) (*ent.User, error)

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
