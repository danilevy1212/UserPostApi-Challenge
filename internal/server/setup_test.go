package server

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/rs/zerolog"

	"github.com/danilevy1212/UserPostApi-Challenge/internal/config"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/inmemory"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
)

var app Application

func addLoggerToContext(req *http.Request) *http.Request {
	return req.WithContext(logger.WithContext(req.Context(), app.Logger))
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Set testing config
	config.ConfigFetcher = func() config.Config {
		return config.Config{
			IsDev: true,
			Port:  8080,
		}
	}

	c :=  config.New()
	l := zerolog.Nop()
	app.Config = &c
	app.Router = gin.New()
	app.DB = &inmemory.InMemoryDB{}
	app.Logger = &l

	// Freeze time
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2025, 3, 27, 12, 0, 0, 0, time.UTC)
	}
	defer func() {
		zerolog.TimestampFunc = time.Now
	}()

	v := m.Run()

	snaps.Clean(m)
	os.Exit(v)
}
