package server

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Application struct {
	router *gin.Engine
}

func New() Application {
	// TODO  We will need some way here to set gin.SetMode(gin.ReleaseMode)
	//       Probably in the docker container
	r := gin.New()

	// TODO  Remove hardcoding, pass a config or env variable instead
	r.SetTrustedProxies(strings.Split("127.0.0.1", ","))

	return Application{
		router: r,
	}
}

func (a *Application) Serve(port uint) error {
	return a.router.Run(fmt.Sprintf(":%d", port))
}
