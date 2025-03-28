package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	IsDev bool
	Port  uint
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

	return Config{
		IsDev: isDev,
		Port:  uint(port),
	}
}
