package server

// TODO  TEST!!!!

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) RegisterRoutes() {
	r := a.router

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
