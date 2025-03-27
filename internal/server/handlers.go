package server

import (
	serverLogger "github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Application) HealthCheck(ctx *gin.Context) {
	logger := serverLogger.FromContext(ctx.Request.Context()).
		With().
		Str("handler", "HealthCheck").
		Logger()

	logger.Info().
		Msg("Health check OK")

	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
