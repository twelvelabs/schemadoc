package main

import (
	"os"

	"github.com/twelvelabs/schemadoc/internal/cmd"
	"github.com/twelvelabs/schemadoc/internal/core"
)

// Vars passed in by goreleaser (see `builds[0].ldflags` in .goreleaser.yaml).
var (
	version = "dev" // release version
	commit  = ""    // release commit SHA
	date    = ""    // release commit date
)

// The actual `main` logic.
// Broken out so we can safely use defer (see [os.Exit] docs).
func run() error {
	path, err := core.ConfigPath(os.Args)
	if err != nil {
		return err
	}

	app, err := core.NewApp(version, commit, date, path)
	if err != nil {
		return err
	}
	defer app.Close()

	return cmd.NewRootCmd(app).ExecuteContext(app.Context())
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
