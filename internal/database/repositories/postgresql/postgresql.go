package postgresql

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent/migrate"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/rs/zerolog"
	"time"
)

type PostgresqlClient struct {
	*ent.Client
	connection *sql.DB
}

func (pg *PostgresqlClient) Connection() *sql.DB {
	return pg.connection
}

func (pg *PostgresqlClient) Ping(ctx context.Context) error {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.Ping").
		Logger()

	log.Info().
		Msg("Pinging DB")

	if err := pg.Connection().Ping(); err != nil {
		log.Error().AnErr("error", err).Msg("Failed to ping DB")
		return err
	}

	log.Info().Msg("successfully pinged database")

	return nil
}

func (pg *PostgresqlClient) UserCreate(ctx context.Context, user ent.User) (*ent.User, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.UserCreate").
		Logger()

	log.Info().
		Interface("user", user).
		Msg("creating user")

	u, err := pg.User.
		Create().
		SetName(user.Name).
		SetEmail(user.Email).
		Save(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while creating user")
	}

	return u, err
}

func (pg *PostgresqlClient) UserGetAll(ctx context.Context) ([]*ent.User, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.UserGetAll").
		Logger()

	users, err := pg.User.Query().All(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while querying users")
	}

	return users, err
}

func (pg *PostgresqlClient) CreateDB(ctx context.Context, l *zerolog.Logger) error {
	logger := l.With().
		Str("method", "postgresql.CreateDB").
		Logger()

	err := pg.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		logger.Error().AnErr("error", err)
		return err
	}

	return nil
}

func New(dns string, parentLogger *zerolog.Logger) *PostgresqlClient {
	logger := parentLogger.
		With().
		Str("method", "postgresql/New").
		Logger()

	logger.Info().
		Msg("Initiating connection with DB")

	db, err := sql.Open("pgx", dns)

	if err != nil {
		logger.Panic().AnErr("error", err).Msg("Failed to open DB connection")
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	drv := entsql.OpenDB(dialect.Postgres, db)
	entClient := ent.NewClient(ent.Driver(drv))

	logger.Info().
		Msg("Successfully connected to DB")

	return &PostgresqlClient{
		entClient,
		db,
	}
}
