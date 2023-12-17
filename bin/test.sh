#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

update=""
if [[ "${UPDATE:-}" != "" ]]; then
    update="-update"
fi

go mod tidy
go test --coverprofile=coverage.out ./... "${update}"
