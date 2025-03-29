package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DBConfig struct {
	username string
	password string
	name     string
	host     string
}

func (dbc DBConfig) String() string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?connect_timeout=5", dbc.username, dbc.password, dbc.host, dbc.name)
}

type Config struct {
	IsDev bool
	Port  uint
	DB    DBConfig
}

type ConfigFunc func() Config

var ConfigFetcher ConfigFunc = fetchFromEnvironment

func New() Config {
	return ConfigFetcher()
}

func fetchFromEnvironment() Config {
	isDev := strings.ToLower(os.Getenv("CHALLENGE_SERVER_IS_PRODUCTION")) != "true"
	port, err := strconv.ParseUint(os.Getenv("CHALLENGE_SERVER_PORT"), 10, 32)

	if err != nil {
		panic(fmt.Sprintf("could not parse `CHALLENGE_SERVER_PORT`: %v", err))
	}

	dbHost := os.Getenv("CHALLENGE_DATABASE_HOST")
	if dbHost == "" {
		panic("database host `CHALLENGE_DATABASE_HOST` is not set")
	}

	dbName := os.Getenv("CHALLENGE_DATABASE_NAME")
	if dbName == "" {
		panic("database name `CHALLENGE_DATABASE_NAME` is not set")
	}

	dbUsername := os.Getenv("CHALLENGE_DATABASE_USERNAME")
	if dbUsername == "" {
		panic("database username `CHALLENGE_DATABASE_USERNAME` is not set")
	}

	dbPassword := os.Getenv("CHALLENGE_DATABASE_PASSWORD")
	if dbPassword == "" {
		panic("database password `CHALLENGE_DATABASE_PASSWORD` is not set")
	}

	return Config{
		IsDev: isDev,
		Port:  uint(port),
		DB: DBConfig{
			username: dbUsername,
			password: dbPassword,
			name:     dbName,
			host:     dbHost,
		},
	}
}
