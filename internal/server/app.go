package server

// TODO  TEST!!!!

import (
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/config"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
)

type Application struct {
	router *gin.Engine
	logger *zerolog.Logger
	config *config.Config
}

func New() Application {
	c := config.New()

	if !c.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// TODO  Remove hardcoding, pass through config
	r.SetTrustedProxies(strings.Split("127.0.0.1", ","))

	l := logger.New(c.IsDev)

	return Application{
		router: r,
		logger: l,
		config: &c,
	}
}

func (a *Application) Serve(port uint) error {
	return a.router.Run(fmt.Sprintf(":%d", a.config.Port))
}
