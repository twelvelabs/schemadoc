#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/build.sh"

rm -fv "${BIN_INSTALL_PATH}"
rm -fv "${CMP_INSTALL_PATH}"
rm -fv "${MAN_INSTALL_PATH}"
