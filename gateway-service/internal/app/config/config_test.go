//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Load config successfully", func(t *testing.T) {
		config := MustInitConfig("../../../.env.sample")

		assert.Equal(t, LogLeveler("info"), config.LogLevel)
		assert.Equal(t, false, config.TracingEnabled)
		assert.Equal(t, 3001, config.HTTP.Port)
		assert.Equal(t, false, config.HTTP.PprofEnabled)
		assert.Equal(t, 3002, config.HTTP.PprofPort)
	})
}
