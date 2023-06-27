#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/build.sh"

if [[ -f "${BIN_BUILD_PATH}" ]]; then
    install -v -d "$(dirname "BIN_INSTALL_PATH")"
    install -v -m755 "${BIN_BUILD_PATH}" "${BIN_INSTALL_PATH}"
fi

if [[ -f "${CMP_BUILD_PATH}" ]]; then
    install -v -d "$(dirname "CMP_INSTALL_PATH")"
    install -v -m755 "${CMP_BUILD_PATH}" "${CMP_INSTALL_PATH}"
fi

if [[ -f "${MAN_BUILD_PATH}" ]]; then
    install -v -d "$(dirname "MAN_INSTALL_PATH")"
    install -v -m644 "${MAN_BUILD_PATH}" "${MAN_INSTALL_PATH}"
fi
