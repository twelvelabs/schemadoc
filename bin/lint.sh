#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

stylist check
echo "[lint] ✅"
