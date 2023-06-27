#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

ensure-dependency() {
    local dependency="${1}"
    local install_command="${2}"

    if command -v "${dependency}" >/dev/null 2>&1; then
        echo "Found dependency: ${dependency}."
    else
        echo "Installing dependency: ${dependency}..."
        $install_command
    fi
}

export HOMEBREW_NO_INSTALL_CLEANUP="1"

ensure-dependency "brew" "echo 'Please follow the instructions at https://brew.sh' && exit 1"

ensure-dependency "actionlint" "brew install --quiet actionlint"
ensure-dependency "gh" "brew install --quiet gh"
ensure-dependency "git" "brew install --quiet git"
ensure-dependency "gitleaks" "brew install --quiet gitleaks"
ensure-dependency "go" "brew install --quiet go"
ensure-dependency "gocovsh" "brew install --quiet orlangure/tap/gocovsh"
ensure-dependency "golangci-lint" "brew install --quiet golangci-lint"
ensure-dependency "goreleaser" "brew install --quiet goreleaser"
ensure-dependency "gum" "brew install --quiet gum"
ensure-dependency "jq" "brew install --quiet jq"
ensure-dependency "npm" "brew install --quiet node"
ensure-dependency "shellcheck" "brew install --quiet shellcheck"
ensure-dependency "shfmt" "brew install --quiet shfmt"
ensure-dependency "stylist" "brew install --quiet twelvelabs/tap/stylist"
ensure-dependency "cspell" "npm install --global --no-audit --no-fund --quiet cspell"
ensure-dependency "markdownlint" "npm install --global --no-audit --no-fund --quiet markdownlint-cli"
ensure-dependency "pin-github-action" "npm install --global --no-audit --no-fund --quiet pin-github-action"

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
if [[ "${CI:-}" == "true" ]]; then
    source "${SCRIPT_DIR}/setup-ci.sh"
else
    source "${SCRIPT_DIR}/setup-local.sh"
fi
