#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Ensure gitlint.
# Not included in dependencies JSON because it takes a long time to install
# in CI (since it depends on python@3.x) and we only use it locally.
ensure-dependency "gitlint" "brew install --quiet gitlint"

# Ensure $USER owns /usr/local/{bin,share}.
# Allows for running `make install` w/out sudo interruptions.
if [[ ! -w /usr/local/bin ]]; then
    echo "Changing owner of /usr/local/bin to $USER."
    sudo chown -R "${USER}" /usr/local/bin
fi
if [[ ! -w /usr/local/share ]]; then
    echo "Changing owner of /usr/local/share to $USER."
    sudo chown -R "${USER}" /usr/local/share
fi

# Ensure local repo.
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    if gum confirm "Create local git repo?"; then
        echo "Creating local repo."
        git init
        git add .
        git commit -m "feat: initial commit"
    fi
fi

# Ensure local repo hooks.
if [[ -d .git ]]; then
    echo "Updating .git/hooks."
    mkdir -p .git/hooks
    rm -f .git/hooks/*.sample
    cp -f bin/githooks/* .git/hooks/
    chmod +x .git/hooks/*
fi

# Ensure remote repo.
if ! gh repo view --json url &>/dev/null; then
    if gum confirm "Create remote git repo?"; then
        echo "Creating remote repo."
        gh repo create
        sleep 1

        echo "Setting 'remotes/origin/HEAD'."
        git remote set-head origin --auto

        echo "Remote repo created: $(gh repo view --json url --jq .url)"
    fi
fi

# Ensure repo secrets.
if gh repo view --json url &>/dev/null; then
    # Set GH_PAT secret
    if ! gh secret list | grep -q GH_PAT; then
        if gum confirm "Add a personal access token to remote git repo?"; then
            if [[ "${GH_PAT:-}" == "" ]]; then
                echo "Enter a GitHub PAT:"
                echo "  - To create a new one, go to https://github.com/settings/tokens/new"
                GH_PAT=$(gum input --password)
            fi
            echo "Setting GH_PAT repo secret."
            gh secret set GH_PAT --body "$GH_PAT"
        fi
    fi

    # Set GH_COMMIT_SIGNING_{KEY,PASS} secrets
    if ! gh secret list | grep -q GH_COMMIT_SIGNING_KEY; then
        if gum confirm "Add commit signing key to remote git repo?"; then
            if [[ "${GH_COMMIT_SIGNING_KEY:-}" == "" ]]; then
                echo "Select a GPG secret key:"
                key_id=$(gpg --list-secret-keys --with-colons |
                    grep '^uid:' |
                    cut -d':' -f10 |
                    gum choose)
                GH_COMMIT_SIGNING_KEY=$(gpg --armor --export-secret-key "$key_id")
            fi
            echo "Setting GH_COMMIT_SIGNING_KEY repo secret."
            gh secret set GH_COMMIT_SIGNING_KEY --body "$GH_COMMIT_SIGNING_KEY"

            if [[ "${GH_COMMIT_SIGNING_PASS:-}" == "" ]]; then
                echo "Enter the password for '$key_id':"
                echo "  - Press enter if no password."
                GH_COMMIT_SIGNING_PASS=$(gum input --password)
            fi
            echo "Setting GH_COMMIT_SIGNING_PASS repo secret."
            gh secret set GH_COMMIT_SIGNING_PASS --body "$GH_COMMIT_SIGNING_PASS"
        fi
    fi
fi
