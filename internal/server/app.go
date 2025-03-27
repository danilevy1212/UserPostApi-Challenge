package server

// TODO  TEST!!!!

import (
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
)

type Application struct {
	router *gin.Engine
	logger *zerolog.Logger
}

func New() Application {
	// TODO  We will need some way here to set gin.SetMode(gin.ReleaseMode)
	//       Probably in the docker container
	r := gin.New()

	// TODO  Remove hardcoding, pass a config or env variable instead
	r.SetTrustedProxies(strings.Split("127.0.0.1", ","))

	// TODO  Hardcoded for now, until I figure out config
	l := logger.New(true)

	return Application{
		router: r,
		logger: l,
	}
}

func (a *Application) Serve(port uint) error {
	return a.router.Run(fmt.Sprintf(":%d", port))
}
