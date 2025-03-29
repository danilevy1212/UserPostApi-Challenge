package main

import (
	"context"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/config"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	log := logger.New(true)

	c := config.New()

	client := postgresql.New(c.DB.String(), log)
	if err := client.CreateDB(context.Background(), log); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create database")
	}

	log.Info().
		Msg("migration succesfully executed, bye!")
}
