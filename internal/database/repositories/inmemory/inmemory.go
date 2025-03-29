package inmemory

import (
	"context"
	"database/sql"
)

type PingFunc func(context.Context) error

var InMemoryDBPingFn PingFunc = func(c context.Context) error {
	return nil
}

type InMemoryDB struct{}

func (im *InMemoryDB) Connection() *sql.DB {
	return nil
}

func (im *InMemoryDB) Ping(ctx context.Context) error {
	return InMemoryDBPingFn(ctx)
}
