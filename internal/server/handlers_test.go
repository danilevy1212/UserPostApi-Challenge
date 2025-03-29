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
}
