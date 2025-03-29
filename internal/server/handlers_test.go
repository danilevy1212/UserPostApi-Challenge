package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/inmemory"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
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
