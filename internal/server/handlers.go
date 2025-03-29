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
		Msg("pinging database")

	if err := a.DB.Ping(ctx.Request.Context()); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to ping database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "service unavailable",
		})
		return
	}

	logger.Info().
		Msg("Health check OK")

	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
