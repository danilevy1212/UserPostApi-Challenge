package server

import (
	"github.com/gin-gonic/gin"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
)

func (a *Application) RegisterMiddleware() {
	r := a.Router

	// GLOBAL
	// Recover from panics
	r.Use(gin.Recovery())
	// Zerolog logger
	r.Use(logger.NewMiddleware(a.Logger))
}
