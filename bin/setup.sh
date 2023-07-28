#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

export HOMEBREW_NO_INSTALL_CLEANUP="1"

ensure-dependency() {
    local dependency="${1}"
    local install_command="${2}"
    if ! command -v "${dependency}" >/dev/null 2>&1; then
        $install_command
    fi
}
ensure-dependency "brew" "echo 'Please follow the instructions at https://brew.sh' && exit 1"
ensure-dependency "depctl" "brew install --quiet twelvelabs/tap/depctl"

depctl up default

if [[ "${CI:-}" == "true" ]]; then
    depctl up ci
else
    depctl up local

    # Ensure git hooks.
    if [[ -d .git ]]; then
        echo "Updating .git/hooks."
        mkdir -p .git/hooks
        rm -f .git/hooks/*.sample
        cp -f bin/githooks/* .git/hooks/
        chmod +x .git/hooks/*
    fi
fi
