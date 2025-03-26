package server

import "github.com/gin-gonic/gin"

func (a *Application) RegisterMiddleware() {
	r := a.router

	// Recover from panics
	r.Use(gin.Recovery())
	// TODO  Zerolog logger
}
