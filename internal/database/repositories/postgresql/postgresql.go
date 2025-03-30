package postgresql

import (
	"context"
	"database/sql"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/rs/zerolog"

	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent/migrate"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
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

// USER
func (pg *PostgresqlClient) UserCreate(ctx context.Context, user models.User) (*models.User, error) {
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

	return &models.User{
		ID:    &u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, err
}

func (pg *PostgresqlClient) UserGetAll(ctx context.Context) ([]*models.User, error) {
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

	result := make([]*models.User, 0, len(users))
	for _, u := range users {
		result = append(result, &models.User{
			ID:    &u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return result, err
}

func (pg *PostgresqlClient) UserGetByID(ctx context.Context, id int) (*models.User, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.UserGet").
		Logger()

	user, err := pg.User.Get(ctx, id)

	if err != nil {
		if !ent.IsNotFound(err) {
			log.Err(err).
				Msg("error while querying user")
		}

		return nil, err
	}

	return &models.User{
		ID:    &id,
		Email: user.Email,
		Name:  user.Name,
	}, err
}

func (pg *PostgresqlClient) UserDeleteByID(ctx context.Context, id int) error {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.UserDeleteByID").
		Logger()

	err := pg.User.DeleteOneID(id).Exec(ctx)

	if err != nil && !ent.IsNotFound(err) {
		log.Err(err).
			Msg("error while deleting user")
	}

	return err
}

func (pg *PostgresqlClient) UserUpdate(ctx context.Context, user models.User) (*models.User, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.UserUpdate").
		Logger()

	u, err := pg.User.UpdateOneID(*user.ID).
		SetName(user.Name).
		SetEmail(user.Email).
		Save(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while updating user")
		return nil, err
	}

	return &models.User{
		ID:    &u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, err
}

// POST

func (pg *PostgresqlClient) PostCreate(ctx context.Context, post models.Post) (*models.Post, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.PostCreate").
		Logger()

	log.Info().
		Interface("post", post).
		Msg("creating post")

	p, err := pg.Post.
		Create().
		SetTitle(post.Title).
		SetContent(post.Content).
		SetUserID(post.UserID).
		Save(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while creating post")
		return nil, err
	}

	return &models.Post{
		ID:      &p.ID,
		Title:   p.Title,
		Content: p.Content,
		UserID:  p.ID,
	}, err
}

func (pg *PostgresqlClient) PostGetAll(ctx context.Context) ([]*models.Post, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.PostGetAll").
		Logger()

	posts, err := pg.Post.Query().All(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while querying posts")
		return nil, err
	}

	result := make([]*models.Post, 0, len(posts))
	for _, p := range posts {
		result = append(result, &models.Post{
			ID:      &p.ID,
			Title:   p.Title,
			Content: p.Content,
			UserID:  p.UserID,
		})
	}

	return result, err
}

func (pg *PostgresqlClient) PostGetByID(ctx context.Context, id int) (*models.Post, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.PostGetByID").
		Logger()

	post, err := pg.Post.Get(ctx, id)

	if err != nil {
		if !ent.IsNotFound(err) {
			log.Err(err).
				Msg("error while querying post")
		}

		return nil, err
	}

	return &models.Post{
		ID:      &id,
		Title:   post.Title,
		Content: post.Content,
		UserID:  post.UserID,
	}, err
}

func (pg *PostgresqlClient) PostDeleteByID(ctx context.Context, id int) error {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.PostDeleteByID").
		Logger()

	err := pg.Post.DeleteOneID(id).Exec(ctx)

	if err != nil && !ent.IsNotFound(err) {
		log.Err(err).
			Msg("error while deleting post")
	}

	return err
}

func (pg *PostgresqlClient) PostUpdate(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
	log := logger.
		FromContext(ctx).
		With().
		Str("method", "postgresql.PostUpdate").
		Logger()

	p, err := pg.Post.UpdateOneID(*post.ID).
		SetTitle(post.Title).
		SetContent(post.Content).
		Save(ctx)

	if err != nil {
		log.Err(err).
			Msg("error while updating post")
		return nil, err
	}

	return &models.Post{
		ID:      &p.ID,
		Title:   p.Title,
		Content: p.Content,
		UserID:  p.UserID,
	}, err
}

// OTHER
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
		logger.Error().Err(err)
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
