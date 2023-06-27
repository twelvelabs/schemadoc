#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

go mod tidy
go test --coverprofile=coverage.out ./...
