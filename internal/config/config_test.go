package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	sampleConfig := Config{
		IsDev: true,
		Port:  9999,
	}

	t.Run("should be able to override fetcher using `ConfigFetcher`", func(t *testing.T) {
		oldConfigFetcher := ConfigFetcher
		ConfigFetcher = func() Config {
			return sampleConfig
		}
		newConfig := New()
		ConfigFetcher = oldConfigFetcher

		assert.Equal(t, sampleConfig, newConfig, "are not the same value")
	})
}

func Test_fetchFromEnvironment(t *testing.T) {
	simpleTests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{"should fetch value from environment when set (true)", "CHALLENGE_SERVER_IS_PRODUCTION", "TRUE", false},
		{"should fetch value from environment when set (false)", "CHALLENGE_SERVER_IS_PRODUCTION", "FALSE", true},
	}
	t.Setenv("CHALLENGE_SERVER_PORT", "3000")

	for _, tt := range simpleTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)
			config := fetchFromEnvironment()
			assert.Equal(t, tt.expected, config.IsDev, "unexpected value")
			assert.Equal(t, config.Port, uint(3000), "should be 3000")
		})
	}

	t.Run("should validate `port` is a valid number", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Setenv("CHALLENGE_SERVER_PORT", "WRONG")
			fetchFromEnvironment()
		}, "should have panicked")
	})
}
