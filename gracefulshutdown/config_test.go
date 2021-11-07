package gracefulshutdown

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("should return default values", func(t *testing.T) {
		cfg := Default()

		assert.Equal(t, time.Duration(15), cfg.Period)
		assert.True(t, cfg.Enabled)
		assert.Nil(t, cfg.ListenErrorHandler)
		assert.Nil(t, cfg.ShutdownFns)
		assert.Nil(t, cfg.Signals)
	})

	t.Run("should return default and set shutdownFns", func(t *testing.T) {
		fns := []func() error{
			func() error {
				return errors.New("mock-error")
			},
			func() error {
				return nil
			},
		}
		cfg := WithShutdownFns(fns)

		assert.Equal(t, time.Duration(15), cfg.Period)
		assert.True(t, cfg.Enabled)
		assert.Nil(t, cfg.ListenErrorHandler)
		assert.Equal(t, fns, cfg.ShutdownFns)
		assert.Nil(t, cfg.Signals)
	})
}
