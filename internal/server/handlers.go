package server

import (
	"net/http"

	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
	"github.com/gin-gonic/gin"
)

func (a *Application) HealthCheck(ctx *gin.Context) {
	log := logger.FromContext(ctx.Request.Context()).
		With().
		Str("handler", "HealthCheck").
		Logger()

	log.Info().
		Msg("pinging database")

	if err := a.DB.Ping(ctx.Request.Context()); err != nil {
		log.Error().
			Err(err).
			Msg("failed to ping database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "service unavailable",
		})
		return
	}

	log.Info().
		Msg("Health check OK")

	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (a *Application) UserCreate(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "UserCreate").
		Logger()

	var user models.User

	err := ctx.ShouldBindBodyWithJSON(&user)

	if err != nil {
		log.Info().
			Err(err).
			Msg("error validating new user")

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "bad entity",
		})
		return
	}

	dbUser, err := a.DB.UserCreate(reqContext, *user.ToEnt())
	if err != nil {
		if ent.IsConstraintError(err) {
			log.Info().
				Interface("user", user).
				Msg("user already exists")

			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		log.Error().
			Err(err).
			Msg("error inserting user in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	user.ID = &dbUser.ID

	ctx.JSON(http.StatusCreated, user)
}

func (a *Application) UserGetAll(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "UserGetAll").
		Logger()

	dbUsers, err := a.DB.UserGetAll(reqContext)

	if err != nil {
		log.Error().
			Err(err).
			Msg("error querying database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	result := make([]models.User, 0, len(dbUsers))
	for _, dbU := range dbUsers {
		user := models.User{
			ID:    &dbU.ID,
			Name:  dbU.Name,
			Email: dbU.Email,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}
