package server

// TODO  TEST!!!!

import (
	"github.com/gin-gonic/gin"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
)

func (a *Application) RegisterMiddleware() {
	r := a.router

	// Recover from panics
	r.Use(gin.Recovery())
	// TODO  Zerolog logger
	r.Use(logger.NewMiddleware(a.logger))
}
