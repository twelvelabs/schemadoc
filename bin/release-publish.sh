#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

if [[ "${CI:-}" == "true" ]]; then
    goreleaser release --clean
else
    # --snapshot creates everything in dist,
    # but does not publish any artifacts.
    # Good for locally testing release config.
    goreleaser release --clean --snapshot
fi
