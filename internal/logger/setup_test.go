package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Gin test mode
	gin.SetMode(gin.TestMode)

	// Freeze time
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2025, 3, 27, 12, 0, 0, 0, time.UTC)
	}
	defer func() {
		zerolog.TimestampFunc = time.Now
	}()

	m.Run()
}
