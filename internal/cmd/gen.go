package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/cobra"
	"github.com/twelvelabs/termite/fsutil"
	"github.com/twelvelabs/termite/render"
	"github.com/twelvelabs/termite/validate"

	"github.com/twelvelabs/schemadoc/internal/core"
	"github.com/twelvelabs/schemadoc/internal/jsonschema"
)

func NewGenCmd(app *core.App) *cobra.Command {
	a := &GenAction{
		App: app,
	}

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate documents",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Run(cmd.Context(), args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&a.InPath, "in", "i", a.InPath, "file path or dir to one or more JSON schema files")
	flags.StringVarP(&a.OutDir, "out", "o", a.OutDir, "output dir to generate files to")
	flags.StringVarP(&a.TemplatePath, "template", "t", a.TemplatePath, "optional template path")

	return cmd
}

type GenAction struct {
	*core.App

	InPath       string `validate:"required"`
	OutDir       string `validate:"required" default:"out"`
	SchemaPaths  []string
	TemplatePath string
}

func (a *GenAction) Run(_ context.Context, _ []string) error {
	if err := a.setup(); err != nil {
		return err
	}

	for _, path := range a.SchemaPaths {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		context := jsonschema.NewContext()
		scm, err := context.Get(path)
		if err != nil {
			return err
		}

		var rendered string
		if a.TemplatePath != "" {
			rendered, err = render.File(a.TemplatePath, scm)
		} else {
			var tpl *template.Template
			templatePath := "templates/markdown.tpl.md"
			tpl, err = template.New(filepath.Base(templatePath)).
				Funcs(render.FuncMap).
				ParseFS(jsonschema.Templates, templatePath)
			if err != nil {
				return err
			}
			buf := bytes.Buffer{}
			err = tpl.Execute(&buf, scm)
			rendered = buf.String()
		}
		if err != nil {
			return err
		}

		renderedPath := filepath.Join(a.OutDir, fmt.Sprintf("%s.md", scm.EntityName()))
		if err := os.WriteFile(renderedPath, []byte(rendered), fsutil.DefaultFileMode); err != nil {
			return err
		}
	}

	return nil
}

func (a *GenAction) setup() error {
	start := time.Now()

	if err := defaults.Set(a); err != nil {
		return err
	}
	if err := validate.Struct(a); err != nil {
		msg := err.Error()
		msg = strings.ReplaceAll(msg, "InPath", `'--in'`)
		msg = strings.ReplaceAll(msg, "OutDir", `'--out'`)
		msg = strings.ReplaceAll(msg, "field", "flag")
		return fmt.Errorf(msg)
	}

	info, err := os.Stat(a.InPath)
	if err != nil {
		return fmt.Errorf(`'--in': %w`, err)
	}
	if info.IsDir() {
		pattern := filepath.Join(a.InPath, "*.schema.json")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf(`'--in': %w`, err)
		}
		a.SchemaPaths = matches
	} else {
		a.SchemaPaths = []string{a.InPath}
	}

	if err := fsutil.EnsureDirWritable(a.OutDir); err != nil {
		return fmt.Errorf(`'--out': %w`, err)
	}

	if a.TemplatePath != "" {
		info, err := os.Stat(a.TemplatePath)
		if err != nil {
			return fmt.Errorf(`'--template': %w`, err)
		}
		if info.IsDir() {
			return fmt.Errorf(`'--template': must not be a directory`)
		}
	}

	a.Logger.Debug(
		"Setup",
		"duration", time.Since(start),
		"in", a.SchemaPaths,
		"out", a.OutDir,
		"template", a.TemplatePath,
	)
	return nil
}
