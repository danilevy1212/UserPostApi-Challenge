package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"

	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/inmemory"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/models"
)

func Test_Application_Health(t *testing.T) {
	app.Router.GET("/health", app.HealthCheck)

	t.Run("should return 200 if PING to DB is ok", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/health", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 if PING to DB fails", func(t *testing.T) {
		oldPingFunc := inmemory.InMemoryDBPingFn
		inmemory.InMemoryDBPingFn = func(ctx context.Context) error {
			return errors.New("some error")
		}
		defer func() {
			inmemory.InMemoryDBPingFn = oldPingFunc
		}()

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/health", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

// USERS
func Test_Application_UserCreate(t *testing.T) {
	app.Router.POST("/users", app.UserCreate)

	tests := []struct {
		Name        string
		StatusCode  int
		RequestBody string
	}{
		{
			"should return 201 if user is created on DB",
			201,
			`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`,
		},
		{
			"should return 422 if user is malformed",
			422,
			`{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := strings.NewReader(tt.RequestBody)
			req := addLoggerToContext(httptest.NewRequest(http.MethodPost, "/users", reader))
			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.StatusCode, w.Code)
			snaps.MatchJSON(t, w.Body.String())
		})
	}

	t.Run("should return 409 if email is already is in used", func(t *testing.T) {
		oldUserCreateFn := inmemory.InMemoryUserCreateFn
		defer func() {
			inmemory.InMemoryUserCreateFn = oldUserCreateFn
		}()

		inmemory.InMemoryUserCreateFn = func(ctx context.Context, u models.User) (*models.User, error) {
			return nil, &ent.ConstraintError{}
		}

		reader := strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)
		req := addLoggerToContext(httptest.NewRequest(http.MethodPost, "/users", reader))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, 409, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 if unknown error occurs", func(t *testing.T) {
		oldUserCreateFn := inmemory.InMemoryUserCreateFn
		defer func() {
			inmemory.InMemoryUserCreateFn = oldUserCreateFn
		}()

		inmemory.InMemoryUserCreateFn = func(ctx context.Context, u models.User) (*models.User, error) {
			return nil, errors.New("something terrible happened")
		}

		reader := strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)
		req := addLoggerToContext(httptest.NewRequest(http.MethodPost, "/users", reader))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, 503, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_UserGetAll(t *testing.T) {
	app.Router.GET("/users", app.UserGetAll)

	t.Run("should return 200 with all data", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 200 with all data when empty", func(t *testing.T) {
		oldUserGetAllFunc := inmemory.InMemoryUserGetAllFn
		defer func() {
			inmemory.InMemoryUserGetAllFn = oldUserGetAllFunc
		}()
		inmemory.InMemoryUserGetAllFn = func(ctx context.Context) ([]*models.User, error) {
			return []*models.User{}, nil
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldUserGetAllFunc := inmemory.InMemoryUserGetAllFn
		defer func() {
			inmemory.InMemoryUserGetAllFn = oldUserGetAllFunc
		}()
		inmemory.InMemoryUserGetAllFn = func(ctx context.Context) ([]*models.User, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_UserGetByID(t *testing.T) {
	app.Router.GET("/users/:id", app.UserGetByID)

	t.Run("should return 200 with user data", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users/hahaha", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 404 when user is not found", func(t *testing.T) {
		oldUserGetByIDFunc := inmemory.InMemoryUserGetByIDFn
		defer func() {
			inmemory.InMemoryUserGetByIDFn = oldUserGetByIDFunc
		}()
		inmemory.InMemoryUserGetByIDFn = func(ctx context.Context, id int) (*models.User, error) {
			return nil, &ent.NotFoundError{}
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldUserGetByIDFunc := inmemory.InMemoryUserGetByIDFn
		defer func() {
			inmemory.InMemoryUserGetByIDFn = oldUserGetByIDFunc
		}()
		inmemory.InMemoryUserGetByIDFn = func(ctx context.Context, id int) (*models.User, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_UserDeleteByID(t *testing.T) {
	app.Router.DELETE("/users/:id", app.UserDeleteByID)

	t.Run("should return 204 when user is deleted", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, w.Body.String(), "")
	})

	t.Run("should return 404 when user is not found", func(t *testing.T) {
		oldUserDeleteByIDFunc := inmemory.InMemoryUserDeleteByIDFn
		defer func() {
			inmemory.InMemoryUserDeleteByIDFn = oldUserDeleteByIDFunc
		}()
		inmemory.InMemoryUserDeleteByIDFn = func(ctx context.Context, id int) error {
			return &ent.NotFoundError{}
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/users/hahaha", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldUserDeleteByIDFunc := inmemory.InMemoryUserDeleteByIDFn
		defer func() {
			inmemory.InMemoryUserDeleteByIDFn = oldUserDeleteByIDFunc
		}()
		inmemory.InMemoryUserDeleteByIDFn = func(ctx context.Context, id int) error {
			return errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/users/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_UserUpdateByID(t *testing.T) {
	app.Router.PUT("/users/:id", app.UserUpdateByID)

	t.Run("should return 200 when user is updated", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 422 when user is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/hahaha", strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 404 when user is not found", func(t *testing.T) {
		oldUserUpdateFunc := inmemory.InMemoryUserUpdateFn
		defer func() {
			inmemory.InMemoryUserUpdateFn = oldUserUpdateFunc
		}()
		inmemory.InMemoryUserUpdateFn = func(ctx context.Context, user models.UserUpdate) (*models.User, error) {
			return nil, &ent.NotFoundError{}
		}
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 409 when email is already in use", func(t *testing.T) {
		oldUserUpdateFunc := inmemory.InMemoryUserUpdateFn
		defer func() {
			inmemory.InMemoryUserUpdateFn = oldUserUpdateFunc
		}()
		inmemory.InMemoryUserUpdateFn = func(ctx context.Context, user models.UserUpdate) (*models.User, error) {
			return nil, &ent.ConstraintError{}
		}
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldUserUpdateFunc := inmemory.InMemoryUserUpdateFn
		defer func() {
			inmemory.InMemoryUserUpdateFn = oldUserUpdateFunc
		}()
		inmemory.InMemoryUserUpdateFn = func(ctx context.Context, user models.UserUpdate) (*models.User, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Daniel Levy Moreno","email":"danielmorenolevy@gmail.com"}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

// POSTS
func Test_Application_PostCreate(t *testing.T) {
	app.Router.POST("/posts", app.PostCreate)

	tests := []struct {
		Name        string
		StatusCode  int
		RequestBody string
	}{
		{
			"should return 201 if post is created on DB",
			201,
			`{"title":"Post Title","content":"Post Content","user_id":1}`,
		},
		{
			"should return 422 if post is malformed",
			422,
			`{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := strings.NewReader(tt.RequestBody)
			req := addLoggerToContext(httptest.NewRequest(http.MethodPost, "/posts", reader))
			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.StatusCode, w.Code)
			snaps.MatchJSON(t, w.Body.String())
		})
	}

	t.Run("should return 503 if unknown error occurs", func(t *testing.T) {
		oldPostCreateFn := inmemory.InMemoryPostCreateFn
		defer func() {
			inmemory.InMemoryPostCreateFn = oldPostCreateFn
		}()

		inmemory.InMemoryPostCreateFn = func(ctx context.Context, p models.Post) (*models.Post, error) {
			return nil, errors.New("something terrible happened")
		}

		reader := strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)
		req := addLoggerToContext(httptest.NewRequest(http.MethodPost, "/posts", reader))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, 503, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_PostGetAll(t *testing.T) {
	app.Router.GET("/posts", app.PostGetAll)

	t.Run("should return 200 with all data", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 200 with all data when empty", func(t *testing.T) {
		oldPostGetAllFunc := inmemory.InMemoryPostGetAllFn
		defer func() {
			inmemory.InMemoryPostGetAllFn = oldPostGetAllFunc
		}()
		inmemory.InMemoryPostGetAllFn = func(ctx context.Context) ([]*models.Post, error) {
			return []*models.Post{}, nil
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldPostGetAllFunc := inmemory.InMemoryPostGetAllFn
		defer func() {
			inmemory.InMemoryPostGetAllFn = oldPostGetAllFunc
		}()
		inmemory.InMemoryPostGetAllFn = func(ctx context.Context) ([]*models.Post, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_PostGetByID(t *testing.T) {
	app.Router.GET("/posts/:id", app.PostGetByID)

	t.Run("should return 200 with post data", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts/hahaha", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 404 when post is not found", func(t *testing.T) {
		oldPostGetByIDFunc := inmemory.InMemoryPostGetByIDFn
		defer func() {
			inmemory.InMemoryPostGetByIDFn = oldPostGetByIDFunc
		}()
		inmemory.InMemoryPostGetByIDFn = func(ctx context.Context, id int) (*models.Post, error) {
			return nil, &ent.NotFoundError{}
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldPostGetByIDFunc := inmemory.InMemoryPostGetByIDFn
		defer func() {
			inmemory.InMemoryPostGetByIDFn = oldPostGetByIDFunc
		}()
		inmemory.InMemoryPostGetByIDFn = func(ctx context.Context, id int) (*models.Post, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodGet, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_PostDeleteByID(t *testing.T) {
	app.Router.DELETE("/posts/:id", app.PostDeleteByID)

	t.Run("should return 204 when post is deleted", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, w.Body.String(), "")
	})

	t.Run("should return 404 when post is not found", func(t *testing.T) {
		oldPostDeleteByIDFunc := inmemory.InMemoryPostDeleteByIDFn
		defer func() {
			inmemory.InMemoryPostDeleteByIDFn = oldPostDeleteByIDFunc
		}()
		inmemory.InMemoryPostDeleteByIDFn = func(ctx context.Context, id int) error {
			return &ent.NotFoundError{}
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/posts/hahaha", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldPostDeleteByIDFunc := inmemory.InMemoryPostDeleteByIDFn
		defer func() {
			inmemory.InMemoryPostDeleteByIDFn = oldPostDeleteByIDFunc
		}()
		inmemory.InMemoryPostDeleteByIDFn = func(ctx context.Context, id int) error {
			return errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodDelete, "/posts/1", nil))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}

func Test_Application_PostUpdateByID(t *testing.T) {
	app.Router.PUT("/posts/:id", app.PostUpdateByID)

	t.Run("should return 200 when post is updated", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 400 when id is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/hahaha", strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 404 when post is not found", func(t *testing.T) {
		oldPostUpdateFn := inmemory.InMemoryPostUpdateFn
		defer func() {
			inmemory.InMemoryPostUpdateFn = oldPostUpdateFn
		}()
		inmemory.InMemoryPostUpdateFn = func(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
			return nil, &ent.NotFoundError{}
		}
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 409 when user id doesn't exists", func(t *testing.T) {
		oldPostUpdateFn := inmemory.InMemoryPostUpdateFn
		defer func() {
			inmemory.InMemoryPostUpdateFn = oldPostUpdateFn
		}()
		inmemory.InMemoryPostUpdateFn = func(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
			return nil, &ent.ConstraintError{}
		}
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 422 when post is malformed", func(t *testing.T) {
		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(`{}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})

	t.Run("should return 503 when unexpected error happens", func(t *testing.T) {
		oldPostUpdateFn := inmemory.InMemoryPostUpdateFn
		defer func() {
			inmemory.InMemoryPostUpdateFn = oldPostUpdateFn
		}()
		inmemory.InMemoryPostUpdateFn = func(ctx context.Context, post models.PostUpdate) (*models.Post, error) {
			return nil, errors.New("You've met a terrible fate, haven't you?")
		}

		req := addLoggerToContext(httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(`{"title":"Post Title","content":"Post Content","user_id":1}`)))
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		snaps.MatchJSON(t, w.Body.String())
	})
}
