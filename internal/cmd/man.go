package cmd

import (
	"fmt"
	"time"

	// cspell: words mangoc roff
	mangoc "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"

	"github.com/twelvelabs/schemadoc/internal/core"
)

func NewManCmd(app *core.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates manpages for the app",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mangoc.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			copyright := fmt.Sprintf(
				"(C) %v Skip Baney.\nReleased under MIT license.",
				time.Now().Year(),
			)
			manPage = manPage.WithSection("Copyright", copyright)

			_, err = fmt.Fprint(app.IO.Out, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}
