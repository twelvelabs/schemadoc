# schemadoc

[![build](https://github.com/twelvelabs/schemadoc/actions/workflows/build.yml/badge.svg)](https://github.com/twelvelabs/schemadoc/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/schemadoc/branch/main/graph/badge.svg?token=jLMxTJcY08)](https://codecov.io/gh/twelvelabs/schemadoc)

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

TODO

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
