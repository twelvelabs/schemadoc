#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/release-status.sh"

if [[ "$CURRENT_VERSION" == "$NEXT_VERSION" ]]; then
    echo "Nothing to tag."
    exit 0
fi

if [[ "${CI:-}" != "true" ]]; then
    # When running locally, prompt before creating the tag (Safety First™).
    if ! gum confirm --default=false "Create and push tag $NEXT_VERSION"; then
        echo "Aborting."
        exit 0
    fi
fi

git tag \
    --sign "$NEXT_VERSION" \
    --message "$NEXT_VERSION"
git push origin --tags
