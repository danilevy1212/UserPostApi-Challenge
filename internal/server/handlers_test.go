package server

import (
	"context"
	"errors"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/inmemory"
	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

		inmemory.InMemoryUserCreateFn = func(ctx context.Context, u ent.User) (*ent.User, error) {
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

		inmemory.InMemoryUserCreateFn = func(ctx context.Context, u ent.User) (*ent.User, error) {
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
		inmemory.InMemoryUserGetAllFn = func(ctx context.Context) ([]*ent.User, error) {
			return []*ent.User{}, nil
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
		inmemory.InMemoryUserGetAllFn = func(ctx context.Context) ([]*ent.User, error) {
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

	t.Run("should return 404 when user is not found", func(t *testing.T) {
		oldUserGetByIDFunc := inmemory.InMemoryUserGetByIDFn
		defer func() {
			inmemory.InMemoryUserGetByIDFn = oldUserGetByIDFunc
		}()
		inmemory.InMemoryUserGetByIDFn = func(ctx context.Context, id int) (*ent.User, error) {
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
		inmemory.InMemoryUserGetByIDFn = func(ctx context.Context, id int) (*ent.User, error) {
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
