# schemadoc

[![build](https://github.com/twelvelabs/schemadoc/actions/workflows/build.yml/badge.svg)](https://github.com/twelvelabs/schemadoc/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/schemadoc/branch/main/graph/badge.svg)](https://codecov.io/gh/twelvelabs/schemadoc)

Generate markdown documents from JSON schema files ✨.

## Installation

Choose one of the following:

- Download and manually install the latest [release](https://github.com/twelvelabs/schemadoc/releases/latest)
- Install with [Homebrew](https://brew.sh/) 🍺

  ```bash
  brew install twelvelabs/tap/schemadoc
  ```

- Install from source

  ```bash
  go install github.com/twelvelabs/schemadoc@latest
  ```

## Usage

```shell
# Renders `./my.schema.json` to `./out/SchemaTitle.md`.
schemadoc gen --in ./my.schema.json

# Renders all json schema files in `./schemas` to `./docs`.
schemadoc gen --in ./schemas --out ./docs
```

To see schemadoc in action, check out
[Generator.md](https://github.com/twelvelabs/stamp/blob/main/docs/Generator.md)
which is rendered from
[stamp.schema.json](https://github.com/twelvelabs/stamp/blob/main/docs/stamp.schema.json)
at build time.

## Customizing

Schemadoc ships with a built in [template](./internal/jsonschema/templates/markdown.tpl.md) for rendering markdown.
To customize (or render something other than markdown)
you can supply your own Go [text/template](https://pkg.go.dev/text/template) file:

```shell
schemadoc gen --in ./schemas --out ./dest --template path/to/my-xml-template.tpl --outfile "{{ .EntityName }}.xml"
```

Each top-level JSON schema in `./schemas` will be parsed into a
[Schema](./internal/jsonschema/schema.go) struct and passed into
`my-xml-template.tpl`.
The rendered files will be written to `./dest/$SchemaName.xml`.

## Development

```shell
git clone git@github.com:twelvelabs/schemadoc.git
cd schemadoc

# Ensures all required dependencies are installed
# and bootstraps the project for local development.
make setup

make build
make test
make install

# Show help.
make
```
