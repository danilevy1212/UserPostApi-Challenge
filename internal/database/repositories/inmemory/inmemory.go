package inmemory

import (
	"context"
	"database/sql"
)

type InMemoryDB struct{}

func (im *InMemoryDB) Connection() *sql.DB {
	return nil
}

func (im *InMemoryDB) Ping(ctx context.Context) error {
	return nil
}
