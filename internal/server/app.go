package server

import (
	"fmt"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/config"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"strings"
)

type Application struct {
	Router *gin.Engine
	Logger *zerolog.Logger
	Config *config.Config
	DB     database.DBRepository
}

func New() Application {
	c := config.New()

	if !c.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.SetTrustedProxies(strings.Split("127.0.0.1", ","))
	r.RemoveExtraSlash = true

	l := logger.New(c.IsDev)

	db := postgresql.New(c.DB.String(), l)

	return Application{
		Router: r,
		Logger: l,
		Config: &c,
		DB:     db,
	}
}

func (a *Application) Serve(port uint) error {
	return a.Router.Run(fmt.Sprintf(":%d", a.Config.Port))
}
