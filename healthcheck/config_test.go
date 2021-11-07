package healthcheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Default Values", func(t *testing.T) {
		cfg := config{}

		assert.False(t, cfg.TimeoutEnabled, "timeout should be false by default")
		assert.Empty(t, cfg.TimeoutPeriod, "timeoutPeriod should be 0 by default")
		assert.Empty(t, cfg.ServiceName, "serviceName should be \"\" by default")
		assert.Nil(t, cfg.HealthChecks, "map should be nil by default")
		assert.False(t, cfg.ShowErrors, "showErrors should be false by default")
	})

	t.Run("EnableTimeout", func(t *testing.T) {
		t.Run("should enable timeout and set default period if not set", func(t *testing.T) {
			cfg := config{}
			EnableTimeout()(&cfg)

			assert.True(t, cfg.TimeoutEnabled)
			assert.Equal(t, time.Duration(30), cfg.TimeoutPeriod)
		})

		t.Run("should not override timeout period", func(t *testing.T) {
			cfg := config{}
			SetTimeoutPeriod(15)(&cfg)
			EnableTimeout()(&cfg)

			assert.True(t, cfg.TimeoutEnabled)
			assert.Equal(t, time.Duration(15), cfg.TimeoutPeriod)
		})
	})

	t.Run("ShowErrors", func(t *testing.T) {
		cfg := config{}
		ShowErrors()(&cfg)

		assert.True(t, cfg.ShowErrors, "should set show errors to true")
	})

	t.Run("SetTimeoutPeriod", func(t *testing.T) {
		t.Run("should set timeout period and auto enable timeout if false", func(t *testing.T) {
			cfg := config{}
			SetTimeoutPeriod(20)(&cfg)

			assert.True(t, cfg.TimeoutEnabled)
			assert.Equal(t, time.Duration(20), cfg.TimeoutPeriod)
		})
	})

	t.Run("RegisterHealthChecks", func(t *testing.T) {
		cfg := config{}
		checks := HealthchecksMap{"test": func() error { return nil }}

		RegisterHealthChecks(checks)(&cfg)

		assert.Equal(t, checks, cfg.HealthChecks)
	})
}
