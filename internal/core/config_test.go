package core

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		want      *Config
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "should return valid config",
			args: args{
				path: filepath.Join("testdata", "config", "valid.yaml"),
			},
			want: &Config{
				ConfigPath: filepath.Join("testdata", "config", "valid.yaml"),
				Color:      false,
				Debug:      false,
				Prompt:     true,
				LogLevel:   "warn",
			},
			assertion: assert.NoError,
		},
		{
			name: "should return error if malformed",
			args: args{
				path: filepath.Join("testdata", "config", "malformed.yaml"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFromPath(tt.args.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigPath(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		setup     func(t *testing.T)
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "should return default config path",
			want:      ConfigPathDefault,
			assertion: assert.NoError,
		},
		{
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(ConfigPathEnv, "from_env.yaml")
			},
			name:      "should return config path from env",
			want:      "from_env.yaml",
			assertion: assert.NoError,
		},
		{
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(ConfigPathEnv, "from_env.yaml")
			},
			args: args{
				args: []string{
					"--config",
					"from_flag.yaml",
				},
			},
			name:      "should return config path from long flag",
			want:      "from_flag.yaml",
			assertion: assert.NoError,
		},
		{
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(ConfigPathEnv, "from_env.yaml")
			},
			args: args{
				args: []string{
					"-c",
					"from_flag.yaml",
				},
			},
			name:      "should return config path from short flag",
			want:      "from_flag.yaml",
			assertion: assert.NoError,
		},
		{
			args: args{
				args: []string{
					"---",
				},
			},
			name:      "should return an error when unable to parse args",
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			got, err := ConfigPath(tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
