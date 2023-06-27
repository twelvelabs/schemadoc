package core

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yaml")
	app, err := NewApp("", "", "", path)
	defer app.Close()

	assert.NotNil(t, app)
	assert.NoError(t, err)
	assert.Equal(t, path, app.Config.ConfigPath)
}

func TestNewApp_WhenConfigError(t *testing.T) {
	path := filepath.Join("testdata", "config", "malformed.yaml")
	app, err := NewApp("", "", "", path)

	assert.Nil(t, app)
	assert.Error(t, err)
}

func TestNewTestApp(t *testing.T) {
	app := NewTestApp()
	defer app.Close()

	assert.NotNil(t, app)
}

func TestAppForContext(t *testing.T) {
	app := NewTestApp()
	ctx := app.Context()

	assert.Equal(t, app, AppForContext(ctx))
}

func TestApp_Init(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yaml")
	app, err := NewApp("", "", "", path)
	assert.NoError(t, err)

	assert.Nil(t, app.IO)
	assert.Nil(t, app.UI)
	assert.Nil(t, app.Logger)

	err = app.Init()
	assert.NoError(t, err)

	assert.NotNil(t, app.IO)
	assert.NotNil(t, app.UI)
	assert.NotNil(t, app.Logger)
}

func TestApp_Init_WhenNoColor(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yaml")
	app, err := NewApp("", "", "", path)
	assert.NoError(t, err)

	err = app.Init()
	assert.NoError(t, err)
	app.IO.SetColorEnabled(true)

	app.Config.Color = false
	err = app.Init()
	assert.NoError(t, err)

	assert.Equal(t, false, app.IO.IsColorEnabled())
}

func TestApp_Init_WhenNoPrompt(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yaml")
	app, err := NewApp("", "", "", path)
	assert.NoError(t, err)

	err = app.Init()
	assert.NoError(t, err)
	app.IO.SetInteractive(true)

	app.Config.Prompt = false
	err = app.Init()
	assert.NoError(t, err)

	assert.Equal(t, false, app.IO.IsInteractive())
}
