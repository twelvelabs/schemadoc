package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/twelvelabs/termite/conf"
)

const (
	ConfigPathDefault = ".schemadoc.yaml"
	ConfigPathEnv     = "SCHEMADOC_CONFIG"
)

type Config struct {
	ConfigPath string
	Color      bool   `yaml:"color" env:"SCHEMADOC_COLOR" default:"true"`
	Debug      bool   `yaml:"debug" env:"SCHEMADOC_DEBUG"`
	Prompt     bool   `yaml:"prompt" env:"SCHEMADOC_PROMPT" default:"true"`
	LogLevel   string `yaml:"log_level" env:"SCHEMADOC_LOG_LEVEL" default:"warn" validate:"oneof=debug info warn error fatal"` //nolint: lll
}

// NewTestConfig returns a new Config for unit tests
// populated with default values.
func NewTestConfig() (*Config, error) {
	return NewConfigFromPath("")
}

// NewConfigFromPath returns a new config for the file at path.
func NewConfigFromPath(path string) (*Config, error) {
	config, err := conf.NewLoader(&Config{}, path).Load()
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}
	config.ConfigPath = path
	return config, nil
}

// ConfigPath resolves and returns the config path.
// Lookup order:
//   - Flag
//   - Environment variable
//   - Default path name
func ConfigPath(args []string) (string, error) {
	path := ConfigPathDefault
	if p := os.Getenv(ConfigPathEnv); p != "" {
		path = p
	}

	// Create a minimal, duplicate flag set to determine just the config path
	// (the remaining flags are defined on the cobra.Command flag set).
	// Using two different sets because Cobra doesn't parse flags until _after_
	// we have instantiated the app (and thus the Config).
	fs := pflag.NewFlagSet("config-args", pflag.ContinueOnError)
	fs.StringVarP(&path, "config", "c", path, "")
	// Ignore all the flags used by the main Cobra flagset.
	fs.ParseErrorsAllowlist.UnknownFlags = true
	// Suppress the default usage shown when the `--help` flag is present
	// (otherwise we end up w/ a duplicate of what Cobra shows).
	fs.Usage = func() {}

	err := fs.Parse(args)
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		return "", fmt.Errorf("unable to parse config flag: %w", err)
	}

	return path, nil
}
