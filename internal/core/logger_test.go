package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/ui"
)

func TestNewLogger(t *testing.T) {
	ios := ui.NewTestIOStreams()
	config, _ := NewTestConfig()
	assert.Equal(t, false, config.Debug)

	logger := NewLogger(ios, config)
	assert.Equal(t, "warn", logger.GetLevel().String())
}

func TestNewLogger_WhenDebug(t *testing.T) {
	ios := ui.NewTestIOStreams()
	config, _ := NewTestConfig()
	config.Debug = true

	logger := NewLogger(ios, config)
	assert.Equal(t, "debug", logger.GetLevel().String())
}
