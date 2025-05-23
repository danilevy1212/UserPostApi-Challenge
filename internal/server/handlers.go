package server

import (
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/logger"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// USERS
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

	dbUser, err := a.DB.UserCreate(reqContext, user)
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

	user.ID = dbUser.ID

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
			ID:    dbU.ID,
			Name:  dbU.Name,
			Email: dbU.Email,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *Application) UserGetByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "UserGetByID").
		Logger()

	idRaw := ctx.Param("id")

	id, err := strconv.ParseUint(idRaw, 10, 64)

	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dbUser, err := a.DB.UserGetByID(reqContext, id)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("user not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		log.Error().
			Err(err).
			Msg("error querying database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	user := models.User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}

	ctx.JSON(http.StatusOK, user)
}

func (a *Application) UserDeleteByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "UserDeleteByID").
		Logger()

	idRaw := ctx.Param("id")

	id, err := strconv.ParseUint(idRaw, 10, 64)

	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = a.DB.UserDeleteByID(reqContext, id)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("user not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		log.Error().
			Uint64("user.id", id).
			Err(err).
			Msg("error deleting user in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (a *Application) UserUpdateByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "UserUpdateByID").
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

	idRaw := ctx.Param("id")
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	updatedUser, err := a.DB.UserUpdate(reqContext, models.UserUpdate{
		ID:    &id,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("user not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if ent.IsConstraintError(err) {
			log.Info().
				Interface("user", user).
				Msg("email already exists")

			ctx.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}

		log.Error().
			Err(err).
			Msg("error updating user in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

// POSTS
func (a *Application) PostCreate(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "PostCreate").
		Logger()

	var post models.Post
	err := ctx.ShouldBindBodyWithJSON(&post)

	if err != nil {
		log.Info().
			Err(err).
			Msg("error validating new post")

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "bad entity",
		})
		return
	}

	dbPost, err := a.DB.PostCreate(reqContext, post)
	if err != nil {
		if ent.IsConstraintError(err) {
			log.Info().
				Interface("post", post).
				Msg("associated userID not in DB")

			ctx.JSON(http.StatusConflict, gin.H{"error": "userID doesn't exist"})
			return
		}

		log.Error().
			Err(err).
			Msg("error inserting post in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	post.ID = dbPost.ID

	ctx.JSON(http.StatusCreated, post)
}

func (a *Application) PostGetAll(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "PostGetAll").
		Logger()

	dbPosts, err := a.DB.PostGetAll(reqContext)

	if err != nil {
		log.Error().
			Err(err).
			Msg("error querying database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	result := make([]models.Post, 0, len(dbPosts))
	for _, dbP := range dbPosts {
		post := models.Post{
			ID:      dbP.ID,
			Title:   dbP.Title,
			Content: dbP.Content,
			UserID:  dbP.UserID,
		}

		result = append(result, post)
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *Application) PostGetByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "PostGetByID").
		Logger()

	idRaw := ctx.Param("id")
	id, err := strconv.ParseUint(idRaw, 10, 64)

	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dbPost, err := a.DB.PostGetByID(reqContext, id)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("post not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		log.Error().
			Err(err).
			Msg("error querying database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	post := models.Post{
		ID:      dbPost.ID,
		Title:   dbPost.Title,
		Content: dbPost.Content,
		UserID:  dbPost.UserID,
	}

	ctx.JSON(http.StatusOK, post)
}

func (a *Application) PostDeleteByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "PostDeleteByID").
		Logger()

	idRaw := ctx.Param("id")
	id, err := strconv.ParseUint(idRaw, 10, 64)

	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = a.DB.PostDeleteByID(reqContext, id)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("post not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		log.Error().
			Err(err).
			Msg("error deleting post in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (a *Application) PostUpdateByID(ctx *gin.Context) {
	reqContext := ctx.Request.Context()
	log := logger.FromContext(reqContext).
		With().
		Str("handler", "PostUpdateByID").
		Logger()

	var post models.PostUpdate
	err := ctx.ShouldBindBodyWithJSON(&post)

	if err != nil {
		log.Info().
			Err(err).
			Msg("error validating new post")

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "bad entity",
		})
		return
	}

	idRaw := ctx.Param("id")
	id, err := strconv.ParseUint(idRaw,10,64)

	if err != nil {
		log.Info().
			Str("id", idRaw).
			Msg("invalid id")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post.ID = &id
	updatedPost, err := a.DB.PostUpdate(reqContext, post)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Info().
				Uint64("id", id).
				Msg("post not found")

			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		if ent.IsConstraintError(err) {
			log.Info().
				Interface("post", post).
				Msg("associated userID not in DB")

			ctx.JSON(http.StatusConflict, gin.H{"error": "userID doesn't exist"})
			return
		}

		log.Error().
			Err(err).
			Msg("error updating post in database")

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "service unavailable",
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedPost)
}
