#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

usage() {
    echo "Usage: $(basename "$0") <project>"
}
export project="${1?$(usage)}"

mkdir -p build
go build -o "build/$project" .

rm -rf build/completions
mkdir -p build/completions

# Generate Cobra shell completion scripts.
# See: https://github.com/spf13/cobra/blob/main/shell_completions.md
for sh in bash zsh fish; do
    "build/$project" completion "$sh" >"build/completions/$project.$sh"
done
