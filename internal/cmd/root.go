package cmd

import (
	"github.com/spf13/cobra"

	"github.com/twelvelabs/schemadoc/internal/core"
)

func NewRootCmd(app *core.App) *cobra.Command {
	noColor := false
	noPrompt := false
	verbosity := 0

	cmd := &cobra.Command{
		Use:     "schemadoc",
		Short:   "Generate markdown documents from JSON schema files",
		Version: app.Meta.Version,
		Args:    cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if noColor {
				app.Config.Color = false
			}
			if noPrompt {
				app.Config.Prompt = false
			}
			if verbosity == 1 {
				app.Config.LogLevel = "info"
			} else if verbosity >= 2 {
				app.Config.LogLevel = "debug"
			}
			return app.Init()
		},
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&app.Config.ConfigPath, "config", "c", app.Config.ConfigPath, "config path")
	flags.CountVarP(&verbosity, "verbose", "v", "enable verbose logging (increase via -vv)")
	flags.BoolVar(&noColor, "no-color", noColor, "do not use color output")
	flags.BoolVar(&noPrompt, "no-prompt", noPrompt, "do not prompt for input")

	cmd.SetVersionTemplate("{{.Version}}\n")

	cmd.AddCommand(NewGenCmd(app))
	cmd.AddCommand(NewManCmd(app))
	cmd.AddCommand(NewVersionCmd(app))

	return cmd
}
