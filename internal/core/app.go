package core

import (
	"context"
	"os"
	"time"

	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

const (
	ctxKeyApp ctxKey = "github.com/twelvelabs/schemadoc/internal/core.App"
)

// App contains global and/or singleton application data.
type App struct {
	Config *Config
	Logger *Logger
	Meta   *Meta
	IO     *ui.IOStreams
	UI     *ui.UserInterface

	ctx context.Context //nolint: containedctx
}

// AppForContext returns the app singleton stored in the given context.
func AppForContext(ctx context.Context) *App {
	return ctx.Value(ctxKeyApp).(*App)
}

// NewApp returns the default App singleton.
// It will be minimally initialized with metadata and config.
// Call `Init()` after flag parsing to complete initialization.
func NewApp(version, commit, date, path string) (*App, error) {
	config, err := NewConfigFromPath(path)
	if err != nil {
		return nil, err
	}
	app := &App{
		Config: config,
		Meta:   NewMeta(version, commit, date),
	}
	return app, nil
}

// NewTestApp returns the test App singleton.
// All properties will be configured for testing (mocks, stubs, etc).
func NewTestApp() *App {
	config, _ := NewTestConfig()
	app := &App{
		Config: config,
		Meta:   NewMeta("test", "", time.Now().Format(time.RFC3339)),
	}
	app.InitForTest()
	return app
}

// Init initializes and configures the app.
// It must be called once flags have been parsed.
func (a *App) Init() error {
	start := time.Now()

	if !a.Config.Color {
		_ = os.Setenv("NO_COLOR", "1")
	}

	a.IO = ui.NewIOStreams()
	if !a.Config.Prompt {
		a.IO.SetInteractive(false)
	}
	a.UI = ui.NewUserInterface(a.IO)
	a.Logger = NewLogger(a.IO, a.Config)

	a.Logger.Debug(
		"App initialized",
		"config", a.Config.ConfigPath,
		"duration", time.Since(start),
	)
	return nil
}

// Init initializes and configures the app for unit testing.
func (a *App) InitForTest() {
	a.IO = ui.NewTestIOStreams()
	a.UI = ui.NewUserInterface(a.IO)
	a.Logger = NewLogger(a.IO, a.Config)
}

// Close ensures all app resources have been closed.
func (a *App) Close() error {
	// Add any db/file closing here if needed.
	return nil
}

// Context returns the root [context.Context] for the app.
func (a *App) Context() context.Context {
	if a.ctx == nil {
		a.ctx = context.WithValue(context.Background(), ctxKeyApp, a)
	}
	return a.ctx
}
