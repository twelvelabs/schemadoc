#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

usage() {
    echo "Usage: $(basename "$0") <project>"
}
export project="${1?$(usage)}"

mkdir -p build
go build -o "build/$project" .

rm -rf build/manpages
mkdir -p build/manpages

# Generate Cobra man pages.
# See: https://github.com/spf13/cobra/blob/main/doc/man_docs.md
"build/$project" man | gzip -c -9 >"build/manpages/${project}.1.gz"
