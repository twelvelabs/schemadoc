#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# cspell: words koozz taggerdate
if ! gh extension list | grep -q "koozz/gh-semver"; then
    gh extension install koozz/gh-semver &>/dev/null
fi

CURRENT_BRANCH=$(git symbolic-ref --short HEAD)
CURRENT_VERSION=$(gh release view --json tagName --jq .tagName || true)
CURRENT_VERSION_LOCAL=$(git tag --sort=taggerdate | tail -1)

if [[ "${CURRENT_VERSION}" != "${CURRENT_VERSION_LOCAL}" ]]; then
    # The semver plugin uses local tags to determine the current version
    # and will report an incorrect version if they're not up to date.
    echo "Fetching tags..."
    git fetch --all --tags 1>&2
fi

NEXT_VERSION=$(gh semver)
if [[ "${CURRENT_VERSION}" != "" ]]; then
    COMMITS_REF="$CURRENT_VERSION...HEAD"
else
    COMMITS_REF="HEAD"
fi
COMMITS=$(git log --color=always --format=" - %C(yellow)%h%Creset %s" "$COMMITS_REF") # cspell: disable-line

echo ""
echo "Current branch:  $CURRENT_BRANCH"
echo "Current version: $CURRENT_VERSION"
echo -n "Unreleased commits:"
if [[ "$COMMITS" != "" ]]; then
    echo ""
    echo "$COMMITS"
else
    echo " <none>"
fi
echo ""
if [[ "$CURRENT_VERSION" != "$NEXT_VERSION" ]]; then
    echo "Next version: $NEXT_VERSION"
else
    echo "Next version: <none>"
fi
echo ""
