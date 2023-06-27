package core

import (
	"time"

	// cspell: words charmbracelet
	"github.com/charmbracelet/log"
	"github.com/twelvelabs/termite/ui"
)

// Logger is a simple wrapper around the [charmbracelet/log] package.
type Logger struct {
	*log.Logger
}

// NewLogger returns a new Logger.
func NewLogger(ios *ui.IOStreams, config *Config) *Logger {
	level := config.LogLevel
	if config.Debug {
		level = "debug"
	}
	return &Logger{
		Logger: log.NewWithOptions(ios.Err, log.Options{
			Level:           log.ParseLevel(level),
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.Kitchen,
		}),
	}
}
